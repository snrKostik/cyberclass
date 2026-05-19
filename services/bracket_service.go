package services

import (
	"context"
	"fmt"

	"github.com/snrkostik/cyberclass/models"
	"github.com/snrkostik/cyberclass/store"
)

type BracketService struct {
	store *store.SQLiteStore
}

func NewBracketService(
	store *store.SQLiteStore,
) *BracketService {
	return &BracketService{
		store: store,
	}
}

func (s *BracketService) CreateSingleEliminationBracket(
	ctx context.Context,
	tournamentID int64,
	teamIDs []int64,
) error {

	if len(teamIDs) != 16 {
		return fmt.Errorf(
			"single elimination requires exactly 16 teams",
		)
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	matchMap := map[string]int64{}

	rounds := map[int]int{
		1: 8,
		2: 4,
		3: 2,
		4: 1,
	}

	// create empty matches
	for round, count := range rounds {

		for pos := 1; pos <= count; pos++ {

			res, err := tx.ExecContext(
				ctx,
				`
				INSERT INTO matches (
					tournament_id,
					round,
					position,
					status
				)
				VALUES (?, ?, ?, ?)
				`,
				tournamentID,
				round,
				pos,
				models.MatchStatusPending,
			)

			if err != nil {
				return err
			}

			id, err := res.LastInsertId()
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%d:%d", round, pos)

			matchMap[key] = id
		}
	}

	// connect matches
	for round := 1; round <= 3; round++ {

		matchCount := rounds[round]

		for pos := 1; pos <= matchCount; pos++ {

			currentKey := fmt.Sprintf(
				"%d:%d",
				round,
				pos,
			)

			currentID := matchMap[currentKey]

			nextRound := round + 1

			nextPos := (pos + 1) / 2

			nextKey := fmt.Sprintf(
				"%d:%d",
				nextRound,
				nextPos,
			)

			nextID := matchMap[nextKey]

			slot := 1

			if pos%2 == 0 {
				slot = 2
			}

			_, err := tx.ExecContext(
				ctx,
				`
				UPDATE matches
				SET
					next_match_id = ?,
					next_match_slot = ?
				WHERE id = ?
				`,
				nextID,
				slot,
				currentID,
			)

			if err != nil {
				return err
			}
		}
	}

	// seed first round
	for i := 0; i < 8; i++ {

		key := fmt.Sprintf("1:%d", i+1)

		matchID := matchMap[key]

		team1 := teamIDs[i*2]
		team2 := teamIDs[i*2+1]

		_, err := tx.ExecContext(
			ctx,
			`
			UPDATE matches
			SET
				team1_id = ?,
				team2_id = ?
			WHERE id = ?
			`,
			team1,
			team2,
			matchID,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
