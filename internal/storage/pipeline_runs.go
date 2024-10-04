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
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Duration   float64   `db:"duration"`
}

type PipelineRunStore struct {
	db *sql.DB
}

func (s *PipelineRunStore) Create(pipelineRun PipelineRun) (string, error) {
	var id string
	query := `
	INSERT INTO pipeline_runs (commit_sha, config_file, branch, status, duration)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	err := s.db.QueryRow(query, pipelineRun.CommitSHA, pipelineRun.ConfigFile, pipelineRun.Branch, pipelineRun.Status, pipelineRun.Duration).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
