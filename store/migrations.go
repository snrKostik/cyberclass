package store

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS accounts (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,

    username        TEXT NOT NULL UNIQUE,

    password_hash   TEXT NOT NULL,

    created_at      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tournaments (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,

    name            TEXT NOT NULL,
    game            TEXT NOT NULL,

    format          INTEGER NOT NULL,
    status          INTEGER NOT NULL DEFAULT 0,

    created_at      INTEGER NOT NULL,
    started_at      INTEGER,
    ended_at        INTEGER
);

CREATE TABLE IF NOT EXISTS teams (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL UNIQUE,
	slogan TEXT NOT NULL DEFAULT '',
    created_at      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS players (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,

    nickname        TEXT NOT NULL UNIQUE,

    real_name       TEXT,

    created_at      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS team_members (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,

    team_id         INTEGER NOT NULL,

    player_id       INTEGER NOT NULL,

    joined_at       INTEGER NOT NULL,

    FOREIGN KEY (team_id)
    REFERENCES teams(id)
    ON DELETE CASCADE,

    FOREIGN KEY (player_id)
    REFERENCES players(id)
    ON DELETE CASCADE,

    UNIQUE(team_id, player_id)
);

CREATE TABLE IF NOT EXISTS tournament_teams (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,

    tournament_id   INTEGER NOT NULL,

    team_id         INTEGER NOT NULL,

    seed            INTEGER,

    joined_at       INTEGER NOT NULL,

    FOREIGN KEY (tournament_id)
    REFERENCES tournaments(id)
    ON DELETE CASCADE,

    FOREIGN KEY (team_id)
    REFERENCES teams(id)
    ON DELETE CASCADE,

    UNIQUE(tournament_id, team_id),

    UNIQUE(tournament_id, seed)
);

CREATE TABLE IF NOT EXISTS matches (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,

    tournament_id       INTEGER NOT NULL,

    round               INTEGER NOT NULL,

    position            INTEGER NOT NULL,

    status              INTEGER NOT NULL DEFAULT 0,

    team1_id            INTEGER,
    team2_id            INTEGER,

    score1              INTEGER NOT NULL DEFAULT 0,
    score2              INTEGER NOT NULL DEFAULT 0,

    winner_team_id      INTEGER,

    next_match_id       INTEGER,
    next_match_slot     INTEGER,

    started_at          INTEGER,
    ended_at            INTEGER,

    FOREIGN KEY (tournament_id)
    REFERENCES tournaments(id)
    ON DELETE CASCADE,

    FOREIGN KEY (team1_id)
    REFERENCES teams(id)
    ON DELETE SET NULL,

    FOREIGN KEY (team2_id)
    REFERENCES teams(id)
    ON DELETE SET NULL,

    FOREIGN KEY (winner_team_id)
    REFERENCES teams(id)
    ON DELETE SET NULL,

    FOREIGN KEY (next_match_id)
    REFERENCES matches(id)
    ON DELETE SET NULL,

    CHECK (
        next_match_slot IN (1, 2)
        OR next_match_slot IS NULL
    ),

    CHECK (
        status IN (0, 1, 2)
    ),

    UNIQUE (tournament_id, round, position)
);

CREATE TABLE IF NOT EXISTS awards (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,

    tournament_id       INTEGER NOT NULL,

    title               TEXT NOT NULL,

    description         TEXT,

    team_id             INTEGER,

    created_at          INTEGER NOT NULL,

    FOREIGN KEY (tournament_id)
    REFERENCES tournaments(id)
    ON DELETE CASCADE,

    FOREIGN KEY (team_id)
    REFERENCES teams(id)
    ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS timer_states (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,

    match_id            INTEGER NOT NULL UNIQUE,

    duration_seconds    INTEGER NOT NULL,

    started_at          INTEGER,

    is_running          INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (match_id)
    REFERENCES matches(id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_matches_tournament
ON matches(tournament_id);

CREATE INDEX IF NOT EXISTS idx_matches_next_match
ON matches(next_match_id);

CREATE INDEX IF NOT EXISTS idx_tournament_teams_tournament
ON tournament_teams(tournament_id);

CREATE INDEX IF NOT EXISTS idx_team_members_team
ON team_members(team_id);
`
