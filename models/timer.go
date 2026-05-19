package models

type TimerState struct {
	ID int64

	MatchID int64

	DurationSeconds int

	StartedAt *int64

	IsRunning bool
}
