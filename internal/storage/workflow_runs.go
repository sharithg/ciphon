package storage

import (
	"database/sql"
	"time"
)

type WorkflowRun struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	PipelineRunID string    `db:"pipeline_run_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type WorkflowRunStore struct {
	db *sql.DB
}

func (s *WorkflowRunStore) Create(workflowRun WorkflowRun) (string, error) {
	var id string
	query := `
	INSERT INTO workflow_runs (name, pipeline_run_id)
	VALUES ($1, $2)
	RETURNING id
	`
	err := s.db.QueryRow(query, workflowRun.Name, workflowRun.PipelineRunID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
