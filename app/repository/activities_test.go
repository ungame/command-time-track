package repository

import (
	"context"
	"github.com/ungame/command-time-track/app/models"
	"github.com/ungame/command-time-track/app/pointer"
	"github.com/ungame/command-time-track/db"
	"sort"
	"testing"
	"time"
)

func TestActivitiesRepository(t *testing.T) {

	var (
		ctx      = context.Background()
		conn     = db.New()
		repo     = NewActivitiesRepository(ctx, conn)
		id       int64
		err      error
		activity *models.Activity
	)

	// testing delete and cleanup test data...
	defer func() {
		t.Run("Delete should delete an existing activity by id", func(_ *testing.T) {
			rows, err := repo.Delete(ctx, id)
			if err != nil {
				t.Errorf("unexpected error on delete activity: %s", err.Error())
			}
			if rows != 1 {
				t.Errorf("unexpected affected rows on delete activity: expected=1, got=%d", rows)
			}
		})
	}()

	t.Run("Create should create an activity", func(_ *testing.T) {

		activity = &models.Activity{
			Category:    "test",
			Description: "creating activity",
			StartedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		id, err = repo.Create(ctx, activity)
		if err != nil {
			t.Errorf("unexpected error on create activity: %s", err.Error())
		}
	})

	var existing *models.Activity

	t.Run("Get should retrieve existing activity by id", func(_ *testing.T) {
		existing, err = repo.Get(ctx, id)
		if err != nil {
			t.Errorf("unexpected error on get activity by id: %s", err.Error())
		}

		if existing.Category != activity.Category {
			t.Errorf("unexpected category on get activity by id: expected=%s, got=%s", activity.Category, existing.Category)
		}

		if existing.Description != activity.Description {
			t.Errorf("unexpected description on get activity by id: expected=%s, got=%s", activity.Description, existing.Description)
		}

		if existing.Status != models.StatusStarted {
			t.Errorf("unexpected status on get activity by id: expected=%v, got=%s", models.StatusStarted, existing.Status)
		}

		if existing.StartedAt.IsZero() {
			t.Errorf("unexpected time field on get activity by id: StartedAt=%s", existing.StartedAt.String())
		}

		if existing.UpdatedAt.IsZero() {
			t.Errorf("unexpected time field on get activity by id: UpdatedAt=%s", existing.StartedAt.String())
		}
	})

	t.Run("Update should update an activity by id", func(_ *testing.T) {
		existing.Description = "updating activity"
		existing.Status = models.StatusFinished
		existing.FinishedAt = pointer.New(time.Now())
		existing.UpdatedAt = time.Now()

		rows, err := repo.Update(ctx, existing)
		if err != nil {
			t.Errorf("unexpected error on update activity: %s", err.Error())
		}
		if rows != 1 {
			t.Errorf("unexpected affected rows on update activity: expected=1, got=%d", rows)
		}

	})

	t.Run("GetAll should return a slice of activities", func(_ *testing.T) {
		items, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("unexpected error on get all activities: %s", err.Error())
		}

		index := sort.Search(len(items), func(i int) bool {
			return items[i].ID >= id
		})

		if items[index].ID != existing.ID {
			t.Errorf("unexpected item on get all activities: expected=%v, got=%v", existing, items[index])
		}

		if items[index].Status != models.StatusFinished {
			t.Errorf("unexpected status on get all activities: expected=%v, got=%v", models.StatusFinished, items[index].Status)
		}
	})

	t.Run("Search should return a slice of activities filtering by some term", func(_ *testing.T) {

		term := "test"

		items, err := repo.Search(ctx, term)
		if err != nil {
			t.Errorf("unexpected error on search activity by term: %s", err.Error())
		}
		if len(items) == 0 {
			t.Errorf("unexpected length items on search activity by term: %d", len(items))
		}
	})

	t.Run("GetByStatus should return a slice of activities with status 0", func(_ *testing.T) {

		items, err := repo.GetByStatus(ctx, models.StatusFinished)
		if err != nil {
			t.Errorf("unexpected error on get activity by status: %s", err.Error())
		}
		if len(items) == 0 {
			t.Errorf("unexpected length items on search activity by term: %d", len(items))
		}

		index := sort.Search(len(items), func(i int) bool {
			return items[i].ID >= id
		})

		if items[index].ID != existing.ID {
			t.Errorf("unexpected item on get all activities: expected=%v, got=%v", existing, items[index])
		}

		if items[index].Status != models.StatusFinished {
			t.Errorf("unexpected status on get all activities: expected=%v, got=%v", models.StatusFinished, items[index].Status)
		}
	})
}
