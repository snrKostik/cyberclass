package services

import (
	"context"
	"errors"
	"time"

	"github.com/snrkostik/cyberclass/models"
	"github.com/snrkostik/cyberclass/store"
)

type MatchService struct {
	store *store.SQLiteStore
}

func NewMatchService(
	store *store.SQLiteStore,
) *MatchService {
	return &MatchService{
		store: store,
	}
}

func (s *MatchService) CompleteMatch(
	ctx context.Context,
	matchID int64,
	score1 int,
	score2 int,
) error {

	if score1 == score2 {
		return errors.New(
			"draws are not allowed",
		)
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var match models.Match

	row := tx.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			tournament_id,

			round,
			position,

			status,

			team1_id,
			team2_id,

			score1,
			score2,

			winner_team_id,

			next_match_id,
			next_match_slot,

			started_at,
			ended_at
		FROM matches
		WHERE id = ?
		`,
		matchID,
	)

	err = row.Scan(
		&match.ID,
		&match.TournamentID,

		&match.Round,
		&match.Position,

		&match.Status,

		&match.Team1ID,
		&match.Team2ID,

		&match.Score1,
		&match.Score2,

		&match.WinnerTeamID,

		&match.NextMatchID,
		&match.NextMatchSlot,

		&match.StartedAt,
		&match.EndedAt,
	)

	if err != nil {
		return err
	}

	if match.Team1ID == nil ||
		match.Team2ID == nil {

		return errors.New(
			"match is not ready",
		)
	}

	var winnerID int64

	if score1 > score2 {
		winnerID = *match.Team1ID
	} else {
		winnerID = *match.Team2ID
	}

	now := time.Now().Unix()

	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE matches
		SET
			score1 = ?,
			score2 = ?,

			winner_team_id = ?,

			status = ?,

			ended_at = ?
		WHERE id = ?
		`,
		score1,
		score2,

		winnerID,

		models.MatchStatusCompleted,

		now,

		matchID,
	)

	if err != nil {
		return err
	}

	// final match
	if match.NextMatchID == nil {
		return tx.Commit()
	}

	if match.NextMatchSlot == nil {
		return errors.New(
			"next_match_slot is null",
		)
	}

	if *match.NextMatchSlot == 1 {

		_, err = tx.ExecContext(
			ctx,
			`
			UPDATE matches
			SET team1_id = ?
			WHERE id = ?
			`,
			winnerID,
			*match.NextMatchID,
		)

	} else {

		_, err = tx.ExecContext(
			ctx,
			`
			UPDATE matches
			SET team2_id = ?
			WHERE id = ?
			`,
			winnerID,
			*match.NextMatchID,
		)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}
