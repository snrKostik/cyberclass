package services

import (
	"context"

	"github.com/snrkostik/cyberclass/models"
	"github.com/snrkostik/cyberclass/store"
)

type StandingsService struct {
	store store.Store
}

func NewStandingsService(
	store store.Store,
) *StandingsService {

	return &StandingsService{
		store: store,
	}
}

func (s *StandingsService) ComputeStandings(
	ctx context.Context,
	tournamentID int64,
) (*models.Standings, error) {

	rounds, err := s.store.GetBracketView(
		ctx,
		tournamentID,
	)

	if err != nil {
		return nil, err
	}

	if len(rounds) < 4 {
		return nil, nil
	}

	finalRound := rounds[3]

	if len(finalRound.Matches) == 0 {
		return nil, nil
	}

	final := finalRound.Matches[0]

	if final.WinnerTeamID == nil {
		return nil, nil
	}

	standings := &models.Standings{}

	if final.Team1ID != nil &&
		*final.Team1ID == *final.WinnerTeamID {

		standings.First = final.Team1Name
		standings.Second = final.Team2Name

	} else {

		standings.First = final.Team2Name
		standings.Second = final.Team1Name
	}

	semiRound := rounds[2]

	for _, match := range semiRound.Matches {

		if match.WinnerTeamID == nil {
			continue
		}

		if match.Team1ID != nil &&
			*match.Team1ID == *match.WinnerTeamID {

			standings.Third = append(
				standings.Third,
				match.Team2Name,
			)

		} else {

			standings.Third = append(
				standings.Third,
				match.Team1Name,
			)
		}
	}

	return standings, nil
}
