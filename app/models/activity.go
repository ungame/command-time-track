package models

import "time"

type Activity struct {
	ID          int64
	Category    string
	Description string
	StartedAt   time.Time
	UpdatedAt   time.Time
	StoppedAt   *time.Time
}
