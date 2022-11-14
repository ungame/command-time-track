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
		ctx  = context.Background()
		conn = db.New()
		repo = NewActivitiesRepository(ctx, conn)
		id   int64
		err  error
	)

	// testing delete and cleanup test data...
	defer func() {
		rows, err := repo.Delete(ctx, id)
		if err != nil {
			t.Errorf("unexpected error on delete activity: %s", err.Error())
		}
		if rows != 1 {
			t.Errorf("unexpected affected rows on delete activity: expected=1, got=%d", rows)
		}
	}()

	activity := &models.Activity{
		Category:    "test",
		Description: "creating activity",
		StartedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err = repo.Create(ctx, activity)
	if err != nil {
		t.Errorf("unexpected error on create activity: %s", err.Error())
	}

	existing, err := repo.Get(ctx, id)
	if err != nil {
		t.Errorf("unexpected error on get activity by id: %s", err.Error())
	}

	if existing.Category != activity.Category {
		t.Errorf("unexpected category on get activity by id: expected=%s, got=%s", activity.Category, existing.Category)
	}

	if existing.Description != activity.Description {
		t.Errorf("unexpected description on get activity by id: expected=%s, got=%s", activity.Description, existing.Description)
	}

	if existing.StartedAt.IsZero() {
		t.Errorf("unexpected time field on get activity by id: StartedAt=%s", existing.StartedAt.String())
	}

	if existing.UpdatedAt.IsZero() {
		t.Errorf("unexpected time field on get activity by id: UpdatedAt=%s", existing.StartedAt.String())
	}

	existing.Description = "updating activity"
	existing.StoppedAt = pointer.New(time.Now())
	existing.UpdatedAt = time.Now()

	rows, err := repo.Update(ctx, existing)
	if err != nil {
		t.Errorf("unexpected error on update activity: %s", err.Error())
	}
	if rows != 1 {
		t.Errorf("unexpected affected rows on update activity: expected=1, got=%d", rows)
	}

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
}
