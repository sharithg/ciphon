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

type WorkflowRunInfo struct {
	CommitSHA    string    `db:"commit_sha" json:"commitSha"`
	RepoName     string    `db:"repo_name" json:"repoName"`
	WorkflowName string    `db:"repo_name" json:"workflowName"`
	PipelineID   string    `db:"pipeline_id" json:"pipelineId"`
	WorkflowID   string    `db:"workflow_id" json:"workflowId"`
	Status       *string   `db:"status" json:"status,omitempty"`
	Branch       string    `db:"branch" json:"branch"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	Duration     *int      `db:"duration" json:"duration,omitempty"`
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

func (s *WorkflowRunStore) GetWorkflowRuns() ([]WorkflowRunInfo, error) {
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
	`

	rows, err := s.db.Query(query)
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
