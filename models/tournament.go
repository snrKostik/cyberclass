package models

const (
	TournamentFormatSingleElimination = 1
)

const (
	TournamentStatusPending = iota
	TournamentStatusActive
	TournamentStatusCompleted
)

type Tournament struct {
	ID int64

	Name string
	Game string

	Format int
	Status int

	CreatedAt int64
	StartedAt *int64
	EndedAt   *int64
}
