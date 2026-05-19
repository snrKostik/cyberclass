package models

type Award struct {
	ID int64

	TournamentID int64

	Title       string
	Description *string

	TeamID *int64

	CreatedAt int64
}
