package types

type StartActivityInput struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

type ActivityOutput struct {
	ID          int64   `json:"id"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	StartedAt   string  `json:"started_at"`
	UpdatedAt   string  `json:"updated_at"`
	FinishedAt  *string `json:"stopped_at"`
}

type UpdateActivityInput struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type GetActivityInput struct {
	ID int64 `json:"id"`
}

type DeleteActivityInput struct {
	ID int64 `json:"id"`
}
