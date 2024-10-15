package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	Nodes interface {
		All(ctx context.Context) ([]Node, error)
		GetById(ctx context.Context, id string) (*Node, error)
		Create(ctx context.Context, node Node) (string, error)
		UpdateStatus(ctx context.Context, nodeID string, status string) error
	}
	Repos interface {
		All(ctx context.Context) ([]TsListRepo, error)
		Create(ctx context.Context, repo CreateRepo) (string, error)
	}
	PipelineRunsStore interface {
		Create(ctx context.Context, pipelineRun PipelineRun) (string, error)
	}
	WorkflowRunsStore interface {
		Create(ctx context.Context, workflowRun WorkflowRun) (string, error)
		GetWorkflowRuns(ctx context.Context) ([]TsWorkflowRunInfo, error)
		GetById(ctx context.Context, id string) ([]WorkflowRunSteps, error)
		UpdateStatus(ctx context.Context, id string, status string) error
		UpdateDuration(ctx context.Context, id string, duration float64) error
		UpdateAllStatuses(ctx context.Context, workflowID string) error
	}
	JobRunsStore interface {
		Create(ctx context.Context, jobRun JobRun) (string, error)
		GetByWorkflowId(ctx context.Context, workflowId string) ([]TsJobs, error)
		UpdateStatus(ctx context.Context, id string, status string) error
	}
	StepRunsStore interface {
		Create(ctx context.Context, stepRun StepRun) (string, error)
		GetByJobId(ctx context.Context, jobId string) ([]TsSteps, error)
		UpdateStatus(ctx context.Context, id string, status string) error
		CreateCommandOutput(ctx context.Context, cmd TsCommandOutput) (string, error)
		GetByStepID(ctx context.Context, stepID string) ([]TsCommandOutput, error)
	}
	UsersStore interface {
		Create(ctx context.Context, user User, githubInfo GitHubUserInfo) (*User, error)
		GetByExternalId(ctx context.Context, id string) (*User, error)
		GetById(ctx context.Context, id string) (*TsUserDisplay, error)
	}
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		Nodes:             &NodeStore{pool},
		Repos:             &RepoStore{pool},
		PipelineRunsStore: &PipelineRunStore{pool},
		WorkflowRunsStore: &WorkflowRunStore{pool},
		JobRunsStore:      &JobRunStore{pool},
		StepRunsStore:     &StepRunStore{pool},
		UsersStore:        &UserStore{pool},
	}
}
