package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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

type TsSteps struct {
	Type    string  `json:"type" db:"type"`
	ID      string  `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Command string  `json:"command" db:"command"`
	Status  *string `json:"status" db:"status"`
}

type StepRunStore struct {
	pool *pgxpool.Pool
}

func NewStepRunStore(pool *pgxpool.Pool) *StepRunStore {
	return &StepRunStore{pool: pool}
}

func (s *StepRunStore) Create(ctx context.Context, st StepRun) (string, error) {
	var id string
	query := `
	INSERT INTO step_runs (job_id, step_order, type, name, command, keys, paths)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	err := s.pool.QueryRow(ctx, query, st.JobID, st.StepOrder, st.Type, st.Name, st.Command, st.Keys, st.Paths).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *StepRunStore) GetByJobId(ctx context.Context, jobId string) ([]TsSteps, error) {
	var steps []TsSteps

	query := `
	SELECT type, s.id, s.name, s.command, s.status
	FROM step_runs s
	JOIN job_runs j ON s.job_id = j.id
	WHERE j.id = $1
	ORDER BY step_order
	`

	rows, err := s.pool.Query(ctx, query, jobId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var step TsSteps
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

func (s *StepRunStore) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
	UPDATE step_runs
	SET status = $1
	WHERE id = $2
	`

	_, err := s.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}
	return nil
}

type TsCommandOutput struct {
	ID        string    `json:"id" db:"id"`
	StepID    string    `json:"step_id" db:"step_id"`
	Stdout    string    `json:"stdout" db:"stdout"`
	Type      *string   `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s *StepRunStore) CreateCommandOutput(ctx context.Context, cmd TsCommandOutput) (string, error) {
	var id string
	query := `
	INSERT INTO command_output (step_id, stdout, type)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	err := s.pool.QueryRow(ctx, query, cmd.StepID, cmd.Stdout, cmd.Type).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *StepRunStore) GetByStepID(ctx context.Context, stepID string) ([]TsCommandOutput, error) {
	var outputs []TsCommandOutput

	query := `
	SELECT id, step_id, stdout, type, created_at
	FROM command_output
	WHERE step_id = $1
	`

	rows, err := s.pool.Query(ctx, query, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var output TsCommandOutput
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
