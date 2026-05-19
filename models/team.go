package models

type Team struct {
	ID int64

	Name string

	CreatedAt int64
}

type TournamentTeam struct {
	ID int64

	TournamentID int64
	TeamID       int64

	Seed *int

	JoinedAt int64
}

type TeamMember struct {
	ID int64

	TeamID   int64
	PlayerID int64

	JoinedAt int64
}
