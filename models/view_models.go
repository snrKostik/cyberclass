package models

type MatchView struct {
	ID int64

	Round    int
	Position int

	Status int

	Team1Name string
	Team2Name string

	Score1 int
	Score2 int

	WinnerTeamID *int64

	Team1ID *int64
	Team2ID *int64
}

type BracketRoundView struct {
	Round int

	Matches []MatchView
}
