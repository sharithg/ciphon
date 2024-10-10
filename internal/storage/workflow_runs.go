package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type WorkflowRun struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	PipelineRunID string    `db:"pipeline_run_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type WorkflowRunStore struct {
	pool *pgxpool.Pool
}

type WorkflowRunInfo struct {
	CommitSHA    string    `db:"commit_sha" json:"commitSha"`
	RepoName     string    `db:"repo_name" json:"repoName"`
	WorkflowName string    `db:"repo_name" json:"workflowName"`
	PipelineID   string    `db:"pipeline_id" json:"pipelineId"`
	WorkflowID   string    `db:"workflow_id" json:"workflowId"`
	Status       *string   `db:"status" json:"status,omitempty"`
	Branch       string    `db:"branch" json:"branch"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	Duration     *float64  `db:"duration" json:"duration,omitempty"`
}

type WorkflowRunSteps struct {
	JobID     string    `db:"job_id"`
	StepID    string    `db:"step_id"`
	Command   string    `db:"command"`
	Type      string    `db:"type"`
	Keys      *[]string `db:"keys"`
	Paths     *[]string `db:"paths"`
	StepOrder int       `db:"step_order"`
	Url       string    `db:"url"`
	RepoName  string    `db:"repo_name"`
	CommitSHA string    `db:"commit_sha"`
	Branch    string    `db:"branch"`
	Docker    string    `db:"docker"`
}

func (s *WorkflowRunStore) Create(ctx context.Context, workflowRun WorkflowRun) (string, error) {
	var id string
	query := `
	INSERT INTO workflow_runs (name, pipeline_run_id)
	VALUES ($1, $2)
	RETURNING id
	`
	err := s.pool.QueryRow(ctx, query, workflowRun.Name, workflowRun.PipelineRunID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *WorkflowRunStore) GetWorkflowRuns(ctx context.Context) ([]WorkflowRunInfo, error) {
	var repos []WorkflowRunInfo

	query := `
		select pr.commit_sha,
			r.name,
			pr.id as pipeline_id,
			w.id  as workflow_id,
			w.status,
			w.name as workflow_name,
			pr.branch,
			pr.created_at,
			w.duration
		from workflow_runs w
				join pipeline_runs pr on pr.id = w.pipeline_run_id
				join github_repos r on r.repo_id = pr.repo_id
		order by w.created_at desc
		limit 20
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo WorkflowRunInfo
		err := rows.Scan(&repo.CommitSHA, &repo.RepoName, &repo.PipelineID, &repo.WorkflowID, &repo.Status, &repo.WorkflowName, &repo.Branch, &repo.CreatedAt, &repo.Duration)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *WorkflowRunStore) GetById(ctx context.Context, id string) ([]WorkflowRunSteps, error) {
	query := `
		select j.id as job_id,
			s.id as step_id,
			s.command,
			s.type,
			s.keys,
			s.paths,
			s.step_order,
			r.url,
			r.name as repo_name,
			pr.commit_sha,
			pr.branch,
			j.docker
		from workflow_runs w
			join pipeline_runs pr on pr.id = w.pipeline_run_id
			join github_repos r on r.repo_id = pr.repo_id
			join job_runs j on j.workflow_id = w.id
			join step_runs s on s.job_id = j.id
		where w.id = $1
		order by s.step_order
	`

	rows, err := s.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var stepRuns []WorkflowRunSteps
	for rows.Next() {
		var step WorkflowRunSteps
		err := rows.Scan(
			&step.JobID,
			&step.StepID,
			&step.Command,
			&step.Type,
			&step.Keys,
			&step.Paths,
			&step.StepOrder,
			&step.Url,
			&step.RepoName,
			&step.CommitSHA,
			&step.Branch,
			&step.Docker,
		)
		if err != nil {
			fmt.Println("err", err)
			return nil, err
		}
		stepRuns = append(stepRuns, step)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stepRuns, nil
}

func (s *WorkflowRunStore) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		update workflow_runs
		set status = $1
		where id = $2
	`
	_, err := s.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkflowRunStore) UpdateDuration(ctx context.Context, id string, duration float64) error {
	query := `
		update workflow_runs
		set duration = $1
		where id = $2
	`
	_, err := s.pool.Exec(ctx, query, duration, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkflowRunStore) UpdateAllStatuses(ctx context.Context, workflowID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "UPDATE workflow_runs SET status = NULL WHERE id = $1", workflowID)
	if err != nil {
		return fmt.Errorf("failed to update workflow_runs status: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE job_runs SET status = NULL WHERE workflow_id = $1", workflowID)
	if err != nil {
		return fmt.Errorf("failed to update jobs status: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE step_runs SET status = NULL WHERE job_id IN (SELECT id FROM job_runs WHERE workflow_id = $1)", workflowID)
	if err != nil {
		return fmt.Errorf("failed to update steps status: %w", err)
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM command_output
		WHERE step_id IN (SELECT id
						FROM step_runs
						WHERE job_id in (SELECT id
										FROM job_runs
										WHERE workflow_id = $1))
	`, workflowID)

	if err != nil {
		return fmt.Errorf("failed to delete step commands: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
