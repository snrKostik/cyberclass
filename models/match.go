package models

const (
	MatchStatusPending = iota
	MatchStatusActive
	MatchStatusCompleted
)

type Match struct {
	ID int64

	TournamentID int64

	Round    int
	Position int

	Status int

	Team1ID *int64
	Team2ID *int64

	Score1 int
	Score2 int

	WinnerTeamID *int64

	NextMatchID   *int64
	NextMatchSlot *int

	StartedAt *int64
	EndedAt   *int64
}
