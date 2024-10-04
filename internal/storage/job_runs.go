package storage

import (
	"database/sql"
	"time"
)

type JobRun struct {
	ID         string    `db:"id"`
	WorkflowID string    `db:"workflow_id"`
	Name       string    `db:"name"`
	Docker     string    `db:"docker"`
	Node       string    `db:"node"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type JobRunStore struct {
	db *sql.DB
}

func (s *JobRunStore) Create(jobRun JobRun) (string, error) {
	var id string
	query := `
	INSERT INTO job_runs (workflow_id, name, docker, node)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	err := s.db.QueryRow(query, jobRun.WorkflowID, jobRun.Name, jobRun.Docker, jobRun.Node).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
