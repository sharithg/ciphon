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

type Jobs struct {
	ID     string `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Status string `db:"status" json:"status"`
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

func (s *JobRunStore) GetByWorkflowId(workflowId string) ([]Jobs, error) {
	var jobs []Jobs

	query := `
	select id, name, status from job_runs
	where workflow_id = $1
	`

	rows, err := s.db.Query(query, workflowId)

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

func (s *JobRunStore) UpdateStatus(id, status string) error {
	query := `
		update job_runs
		set status = $1
		where id = $2
	`
	err := s.db.QueryRow(query, status, id).Err()
	if err != nil {
		return err
	}
	return nil
}
