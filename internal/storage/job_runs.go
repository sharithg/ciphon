package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type JobRun struct {
	ID         string    `db:"id"`
	WorkflowID string    `db:"workflow_id"`
	Name       string    `db:"name"`
	Docker     string    `db:"docker"`
	Node       string    `db:"node"`
	Requires   []string  `db:"requires"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Jobs struct {
	ID     string  `db:"id" json:"id"`
	Name   string  `db:"name" json:"name"`
	Status *string `db:"status" json:"status"`
}

type JobRunStore struct {
	pool *pgxpool.Pool
}

func (s *JobRunStore) Create(ctx context.Context, jobRun JobRun) (string, error) {
	var id string
	query := `
	INSERT INTO job_runs (id, workflow_id, name, docker, node, requires)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	err := s.pool.QueryRow(ctx, query, jobRun.ID, jobRun.WorkflowID, jobRun.Name, jobRun.Docker, jobRun.Node, jobRun.Requires).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *JobRunStore) GetByWorkflowId(ctx context.Context, workflowId string) ([]Jobs, error) {
	var jobs []Jobs

	query := `
	SELECT id, name, status 
	FROM job_runs 
	WHERE workflow_id = $1
	`

	rows, err := s.pool.Query(ctx, query, workflowId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var job Jobs
		err := rows.Scan(&job.ID, &job.Name, &job.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobRunStore) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
	UPDATE job_runs
	SET status = $1
	WHERE id = $2
	`
	_, err := s.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}
