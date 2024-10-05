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
	StepOrder int       `db:"step_order"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type StepRunStore struct {
	db *sql.DB
}

func (s *StepRunStore) Create(st StepRun) (string, error) {
	var id string
	query := `
	INSERT INTO step_runs (job_id, step_order, type, name, command, keys, paths)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	err := s.db.QueryRow(query, st.JobID, st.StepOrder, st.Type, st.Name, st.Command, pq.Array(st.Keys), pq.Array(st.Paths)).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
