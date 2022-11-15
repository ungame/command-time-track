package models

import (
	"github.com/ungame/command-time-track/app/pointer"
	"github.com/ungame/command-time-track/app/types"
	"time"
)

type Status string

const (
	StatusFinished Status = "0"
	StatusStarted  Status = "1"
)

func (s Status) String() string {
	if s == "0" {
		return "FINISHED"
	}
	return "STARTED"
}

type Activity struct {
	ID          int64
	Category    string
	Description string
	Status      Status
	StartedAt   time.Time
	UpdatedAt   time.Time
	FinishedAt  *time.Time
}

func (a *Activity) GetFinishedAt() string {
	if a.FinishedAt == nil {
		return ""
	}
	return a.FinishedAt.String()
}

func (a *Activity) Out() *types.ActivityOutput {
	return &types.ActivityOutput{
		ID:          a.ID,
		Category:    a.Category,
		Description: a.Description,
		Status:      a.Status.String(),
		StartedAt:   a.StartedAt.String(),
		UpdatedAt:   a.UpdatedAt.String(),
		FinishedAt:  pointer.New(a.GetFinishedAt()),
	}
}
