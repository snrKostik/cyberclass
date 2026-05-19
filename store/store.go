package store

import (
	"context"

	"github.com/snrkostik/cyberclass/models"
)

type Store interface {
	CreateTournament(
		ctx context.Context,
		tournament *models.Tournament,
	) error

	GetTournament(
		ctx context.Context,
		id int64,
	) (*models.Tournament, error)

	CreateTeam(
		ctx context.Context,
		team *models.Team,
	) error

	GetTeamsByTournament(
		ctx context.Context,
		tournamentID int64,
	) ([]models.Team, error)

	CreateMatch(
		ctx context.Context,
		match *models.Match,
	) error

	UpdateMatch(
		ctx context.Context,
		match *models.Match,
	) error

	GetMatchesByTournament(
		ctx context.Context,
		tournamentID int64,
	) ([]models.Match, error)

	AddTeamToTournament(
		ctx context.Context,
		tournamentTeam *models.TournamentTeam,
	) error

	GetBracketView(
		ctx context.Context,
		tournamentID int64,
	) ([]models.BracketRoundView, error)

	GetMatchByID(
		ctx context.Context,
		id int64,
	) (*models.Match, error)

	GetAllTournaments(
		ctx context.Context,
	) ([]models.Tournament, error)

	TournamentHasBracket(
		ctx context.Context,
		tournamentID int64,
	) (bool, error)
}
