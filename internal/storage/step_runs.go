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

type Steps struct {
	Type    string  `json:"type" db:"type"`
	ID      string  `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Command string  `json:"command" db:"command"`
	Status  *string `json:"status" db:"status"`
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

func (s *StepRunStore) GetByJobId(jobId string) ([]Steps, error) {
	var steps []Steps

	query := `
	select type, s.id, s.name, s.command, s.status
	from step_runs s
	join job_runs j on s.job_id = j.id
	where j.id = $1
	order by step_order
	`

	rows, err := s.db.Query(query, jobId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var step Steps
		err := rows.Scan(&step.Type, &step.ID, &step.Name, &step.Command, &step.Status)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return steps, nil
}

func (s *StepRunStore) UpdateStatus(id, status string) error {
	query := `
		update step_runs
		set status = $1
		where id = $2
	`
	err := s.db.QueryRow(query, status, id).Err()
	if err != nil {
		return err
	}
	return nil
}

type CommandOutput struct {
	ID        string    `json:"id" db:"id"`
	StepID    string    `json:"step_id" db:"step_id"`
	Stdout    string    `json:"stdout" db:"stdout"`
	Type      *string   `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s *StepRunStore) CreateCommandOutput(cmd CommandOutput) (string, error) {
	var id string
	query := `
	INSERT INTO command_output (step_id, stdout, type)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	err := s.db.QueryRow(query, cmd.StepID, cmd.Stdout, cmd.Type).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *StepRunStore) GetByStepID(stepID string) ([]CommandOutput, error) {
	var outputs []CommandOutput

	query := `
	select id, step_id, stdout, type, created_at
	from command_output
	where step_id = $1
	`

	rows, err := s.db.Query(query, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var output CommandOutput
		err := rows.Scan(&output.ID, &output.StepID, &output.Stdout, &output.Type, &output.CreatedAt)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return outputs, nil
}
