package storage

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type StepRun struct {
	ID        string    `db:"id"`
	JobID     string    `db:"job_id"`
	Type      string    `db:"type"`
	Name      string    `db:"name"`
	Command   string    `db:"command"`
	Keys      []string  `db:"keys"`
	Paths     []string  `db:"paths"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type StepRunStore struct {
	db *sql.DB
}

func (s *StepRunStore) Create(stepRun StepRun) (string, error) {
	var id string
	query := `
	INSERT INTO step_runs (job_id, type, name, command, keys, paths)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	err := s.db.QueryRow(query, stepRun.JobID, stepRun.Type, stepRun.Name, stepRun.Command, pq.Array(stepRun.Keys), pq.Array(stepRun.Paths)).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
