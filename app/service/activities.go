package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ungame/command-time-track/app/models"
	"github.com/ungame/command-time-track/app/observer"
	"github.com/ungame/command-time-track/app/pointer"
	"github.com/ungame/command-time-track/app/repository"
	"github.com/ungame/command-time-track/app/types"
	"log"
	"sync"
	"time"
)

type ActivitiesService interface {
	StartActivity(ctx context.Context, input *types.StartActivityInput) (*types.ActivityOutput, error)
	StopActivity(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error)
	UpdateActivityCategory(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error)
	UpdateActivityDescription(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error)
	GetActivityByID(ctx context.Context, input *types.GetActivityInput) (*types.ActivityOutput, error)
	ListActivities(ctx context.Context) ([]*types.ActivityOutput, error)
	SearchActivities(ctx context.Context, term string) ([]*types.ActivityOutput, error)
	DeleteActivityByID(ctx context.Context, input *types.DeleteActivityInput) (int64, error)
	Close()
}

type activitiesService struct {
	activitiesRepository repository.ActivitiesRepository
	activitiesObserver   observer.ActivitiesObserver
	waitGroup            *sync.WaitGroup
}

func NewActivitiesService(activitiesRepository repository.ActivitiesRepository, activitiesObserver observer.ActivitiesObserver) ActivitiesService {
	return &activitiesService{
		activitiesRepository: activitiesRepository,
		activitiesObserver:   activitiesObserver,
		waitGroup:            &sync.WaitGroup{},
	}
}

func (s *activitiesService) StartActivity(ctx context.Context, input *types.StartActivityInput) (*types.ActivityOutput, error) {

	started, err := s.activitiesRepository.GetByStatus(ctx, models.StatusStarted)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if len(started) != 0 {
		s.asyncStopActivities(started)
	}

	activity := &models.Activity{
		Category:    input.Category,
		Description: input.Description,
		StartedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	activity.ID, err = s.activitiesRepository.Create(ctx, activity)
	if err != nil {
		return nil, err
	}

	s.activitiesObserver.Count(activity.Category)

	activity.Status = models.StatusStarted

	log.Printf("Activity created: ID=%v\n", activity.ID)

	return activity.Out(), nil

}

func (s *activitiesService) StopActivity(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error) {

	existing, err := s.activitiesRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if existing.Status != models.StatusFinished {

		existing.Status = models.StatusFinished
		existing.FinishedAt = pointer.New(time.Now().UTC())
		existing.UpdatedAt = time.Now().UTC()

		_, err := s.activitiesRepository.Update(ctx, existing)
		if err != nil {
			return nil, err
		}

		s.activitiesObserver.DurationOf(existing.Category, existing.StartedAt)

		log.Printf("Activity stopped: ID=%v\n", existing.ID)
	}

	return existing.Out(), nil
}

func (s *activitiesService) asyncStopActivities(activities []*models.Activity) {
	for _, activity := range activities {
		s.waitGroup.Add(1)
		go func(id int64) {
			defer s.waitGroup.Done()
			_, err := s.StopActivity(context.Background(), &types.UpdateActivityInput{ID: id})
			if err != nil {
				log.Println("Error on stop activity in background:", err.Error())
			}
		}(activity.ID)
	}
}

func (s *activitiesService) UpdateActivityCategory(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error) {

	existing, err := s.activitiesRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if existing.Category != input.Category {

		existing.Category = input.Category
		existing.UpdatedAt = time.Now().UTC()

		_, err := s.activitiesRepository.Update(ctx, existing)
		if err != nil {
			return nil, err
		}

		log.Printf("Activity category updated: ID=%v\n", existing.ID)
	}

	return existing.Out(), nil
}

func (s *activitiesService) UpdateActivityDescription(ctx context.Context, input *types.UpdateActivityInput) (*types.ActivityOutput, error) {

	existing, err := s.activitiesRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if existing.Description != input.Description {
		existing.Description = input.Description
		existing.UpdatedAt = time.Now().UTC()

		_, err := s.activitiesRepository.Update(ctx, existing)
		if err != nil {
			return nil, err
		}

		log.Printf("Activity description updated: ID=%v\n", existing.ID)
	}

	return existing.Out(), nil
}

func (s *activitiesService) ListActivities(ctx context.Context) ([]*types.ActivityOutput, error) {
	activities, err := s.activitiesRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*types.ActivityOutput, 0, len(activities))
	for _, activity := range activities {
		output = append(output, activity.Out())
	}
	return output, nil
}

func (s *activitiesService) GetActivityByID(ctx context.Context, input *types.GetActivityInput) (*types.ActivityOutput, error) {
	activity, err := s.activitiesRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return activity.Out(), nil
}

func (s *activitiesService) SearchActivities(ctx context.Context, term string) ([]*types.ActivityOutput, error) {
	activities, err := s.activitiesRepository.Search(ctx, term)
	if err != nil {
		return nil, err
	}
	output := make([]*types.ActivityOutput, 0, len(activities))
	for _, activity := range activities {
		output = append(output, activity.Out())
	}
	return output, nil
}

func (s *activitiesService) DeleteActivityByID(ctx context.Context, input *types.DeleteActivityInput) (int64, error) {
	rows, err := s.activitiesRepository.Delete(ctx, input.ID)
	if err != nil {
		return 0, err
	}

	if rows == 0 {
		return 0, fmt.Errorf("unable to delete activity: ID=%v", input.ID)
	}

	log.Printf("Activity deleted: ID=%v\n", input.ID)

	return input.ID, nil
}

func (s *activitiesService) Close() {
	s.waitGroup.Wait()
}
