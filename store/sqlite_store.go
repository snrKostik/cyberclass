package store

import (
	"context"
	"database/sql"

	"github.com/snrkostik/cyberclass/models"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{
		db: db,
	}
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

func (s *SQLiteStore) BeginTx(
	ctx context.Context,
) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *SQLiteStore) CreateTournament(
	ctx context.Context,
	t *models.Tournament,
) error {

	res, err := s.db.ExecContext(
		ctx,
		`
		INSERT INTO tournaments (
			name,
			game,
			format,
			status,
			created_at
		)
		VALUES (?, ?, ?, ?, ?)
		`,
		t.Name,
		t.Game,
		t.Format,
		t.Status,
		t.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = id

	return nil
}

func (s *SQLiteStore) GetTournament(
	ctx context.Context,
	id int64,
) (*models.Tournament, error) {

	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			name,
			game,
			format,
			status,
			created_at,
			started_at,
			ended_at
		FROM tournaments
		WHERE id = ?
		`,
		id,
	)

	var t models.Tournament

	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Game,
		&t.Format,
		&t.Status,
		&t.CreatedAt,
		&t.StartedAt,
		&t.EndedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *SQLiteStore) CreateTeam(
	ctx context.Context,
	team *models.Team,
) error {

	res, err := s.db.ExecContext(
		ctx,
		`
		INSERT INTO teams (
			name,
			slogan,
			created_at
		)
		VALUES (?, ?, ?)
		`,
		team.Name,
		team.Slogan,
		team.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	team.ID = id

	return nil
}

func (s *SQLiteStore) AddTeamToTournament(
	ctx context.Context,
	tt *models.TournamentTeam,
) error {

	res, err := s.db.ExecContext(
		ctx,
		`
		INSERT INTO tournament_teams (
			tournament_id,
			team_id,
			seed,
			joined_at
		)
		VALUES (?, ?, ?, ?)
		`,
		tt.TournamentID,
		tt.TeamID,
		tt.Seed,
		tt.JoinedAt,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	tt.ID = id

	return nil
}

func (s *SQLiteStore) GetTeamsByTournament(
	ctx context.Context,
	tournamentID int64,
) ([]models.Team, error) {

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT
			t.id,
			t.name,
			t.slogan,
			t.created_at
		FROM tournament_teams tt
		INNER JOIN teams t
			ON t.id = tt.team_id
		WHERE tt.tournament_id = ?
		ORDER BY tt.seed ASC
		`,
		tournamentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teams []models.Team

	for rows.Next() {

		var t models.Team

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Slogan,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	return teams, rows.Err()
}

func (s *SQLiteStore) CreateMatch(
	ctx context.Context,
	match *models.Match,
) error {

	res, err := s.db.ExecContext(
		ctx,
		`
		INSERT INTO matches (
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
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		match.TournamentID,
		match.Round,
		match.Position,
		match.Status,

		match.Team1ID,
		match.Team2ID,

		match.Score1,
		match.Score2,

		match.WinnerTeamID,

		match.NextMatchID,
		match.NextMatchSlot,

		match.StartedAt,
		match.EndedAt,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	match.ID = id

	return nil
}

func (s *SQLiteStore) UpdateMatch(
	ctx context.Context,
	match *models.Match,
) error {

	_, err := s.db.ExecContext(
		ctx,
		`
		UPDATE matches
		SET
			status = ?,

			team1_id = ?,
			team2_id = ?,

			score1 = ?,
			score2 = ?,

			winner_team_id = ?,

			next_match_id = ?,
			next_match_slot = ?,

			started_at = ?,
			ended_at = ?
		WHERE id = ?
		`,
		match.Status,

		match.Team1ID,
		match.Team2ID,

		match.Score1,
		match.Score2,

		match.WinnerTeamID,

		match.NextMatchID,
		match.NextMatchSlot,

		match.StartedAt,
		match.EndedAt,

		match.ID,
	)

	return err
}

func (s *SQLiteStore) GetMatchesByTournament(
	ctx context.Context,
	tournamentID int64,
) ([]models.Match, error) {

	rows, err := s.db.QueryContext(
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
		WHERE tournament_id = ?
		ORDER BY round, position
		`,
		tournamentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var matches []models.Match

	for rows.Next() {

		var m models.Match

		err := rows.Scan(
			&m.ID,
			&m.TournamentID,

			&m.Round,
			&m.Position,

			&m.Status,

			&m.Team1ID,
			&m.Team2ID,

			&m.Score1,
			&m.Score2,

			&m.WinnerTeamID,

			&m.NextMatchID,
			&m.NextMatchSlot,

			&m.StartedAt,
			&m.EndedAt,
		)

		if err != nil {
			return nil, err
		}

		matches = append(matches, m)
	}

	return matches, rows.Err()
}

func (s *SQLiteStore) GetBracketView(
	ctx context.Context,
	tournamentID int64,
) ([]models.BracketRoundView, error) {

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT
			m.id,

			m.round,
			m.position,

			m.status,

			COALESCE(t1.name, ''),
			COALESCE(t2.name, ''),

			m.score1,
			m.score2,

			m.team1_id,
			m.team2_id,

			m.winner_team_id

		FROM matches m

		LEFT JOIN teams t1
			ON t1.id = m.team1_id

		LEFT JOIN teams t2
			ON t2.id = m.team2_id

		WHERE m.tournament_id = ?

		ORDER BY
			m.round ASC,
			m.position ASC
		`,
		tournamentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roundMap := map[int][]models.MatchView{}

	for rows.Next() {

		var mv models.MatchView

		err := rows.Scan(
			&mv.ID,

			&mv.Round,
			&mv.Position,

			&mv.Status,

			&mv.Team1Name,
			&mv.Team2Name,

			&mv.Score1,
			&mv.Score2,

			&mv.Team1ID,
			&mv.Team2ID,

			&mv.WinnerTeamID,
		)

		if err != nil {
			return nil, err
		}

		roundMap[mv.Round] = append(
			roundMap[mv.Round],
			mv,
		)
	}

	var rounds []models.BracketRoundView

	for round := 1; round <= 4; round++ {

		rounds = append(
			rounds,
			models.BracketRoundView{
				Round:   round,
				Matches: roundMap[round],
			},
		)
	}

	return rounds, nil
}

func (s *SQLiteStore) GetMatchByID(
	ctx context.Context,
	id int64,
) (*models.Match, error) {

	row := s.db.QueryRowContext(
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
		id,
	)

	var m models.Match

	err := row.Scan(
		&m.ID,
		&m.TournamentID,

		&m.Round,
		&m.Position,

		&m.Status,

		&m.Team1ID,
		&m.Team2ID,

		&m.Score1,
		&m.Score2,

		&m.WinnerTeamID,

		&m.NextMatchID,
		&m.NextMatchSlot,

		&m.StartedAt,
		&m.EndedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *SQLiteStore) GetAllTournaments(
	ctx context.Context,
) ([]models.Tournament, error) {

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT
			id,
			name,
			game,
			format,
			status,
			created_at,
			started_at,
			ended_at
		FROM tournaments
		ORDER BY id DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tournaments []models.Tournament

	for rows.Next() {

		var t models.Tournament

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Game,
			&t.Format,
			&t.Status,
			&t.CreatedAt,
			&t.StartedAt,
			&t.EndedAt,
		)

		if err != nil {
			return nil, err
		}

		tournaments = append(
			tournaments,
			t,
		)
	}

	return tournaments, rows.Err()
}

func (s *SQLiteStore) TournamentHasBracket(
	ctx context.Context,
	tournamentID int64,
) (bool, error) {

	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM matches
		WHERE tournament_id = ?
		`,
		tournamentID,
	)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SQLiteStore) DeleteTeam(
	ctx context.Context,
	id int64,
) error {

	_, err := s.db.ExecContext(
		ctx,
		`
		DELETE FROM teams
		WHERE id = ?
		`,
		id,
	)

	return err
}

func (s *SQLiteStore) DeleteTournament(
	ctx context.Context,
	id int64,
) error {

	_, err := s.db.ExecContext(
		ctx,
		`
		DELETE FROM tournaments
		WHERE id = ?
		`,
		id,
	)

	return err
}

func (s *SQLiteStore) UpdateTournamentStatus(
	ctx context.Context,
	id int64,
	status int,
) error {

	_, err := s.db.ExecContext(
		ctx,
		`
		UPDATE tournaments
		SET status = ?
		WHERE id = ?
		`,
		status,
		id,
	)

	return err
}
