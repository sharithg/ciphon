package storage

import (
	"database/sql"
	"time"
)

type PipelineRun struct {
	ID         string    `db:"id"`
	CommitSHA  string    `db:"commit_sha"`
	ConfigFile string    `db:"config_file"`
	Branch     string    `db:"branch"`
	Status     string    `db:"status"`
	RepoId     int64     `db:"repo_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Duration   float64   `db:"duration"`
}

type PipelineRunStore struct {
	db *sql.DB
}

func (s *PipelineRunStore) Create(p PipelineRun) (string, error) {
	var id string
	query := `
	INSERT INTO pipeline_runs (commit_sha, repo_id, config_file, branch, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	err := s.db.QueryRow(query, p.CommitSHA, p.RepoId, p.ConfigFile, p.Branch, p.Status).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
