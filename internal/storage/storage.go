package storage

import "database/sql"

type Storage struct {
	Nodes interface {
		All() ([]Node, error)
		GetById(id string) (*Node, error)
		Create(node Node) (string, error)
		UpdateStatus(nodeID string, status string) error
	}
	Repos interface {
		All() ([]ListRepo, error)
		Create(repo CreateRepo) (string, error)
	}
	PipelineRunsStore interface {
		Create(pipelineRun PipelineRun) (string, error)
	}
	WorkflowRunsStore interface {
		Create(workflowRun WorkflowRun) (string, error)
		GetWorkflowRuns() ([]WorkflowRunInfo, error)
		GetById(id string) ([]WorkflowRunSteps, error)
		UpdateStatus(id string, status string) error
		UpdateDuration(id string, duration float64) error
		UpdateAllStatuses(workflowID string) error
	}
	JobRunsStore interface {
		Create(jobRun JobRun) (string, error)
		GetByWorkflowId(workflowId string) ([]Jobs, error)
		UpdateStatus(id string, status string) error
	}
	StepRunsStore interface {
		Create(stepRun StepRun) (string, error)
		GetByJobId(jobId string) ([]Steps, error)
		UpdateStatus(id string, status string) error
		CreateCommandOutput(cmd CommandOutput) (string, error)
		GetByStepID(stepID string) ([]CommandOutput, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Nodes:             &NodeStore{db},
		Repos:             &RepoStore{db},
		PipelineRunsStore: &PipelineRunStore{db},
		WorkflowRunsStore: &WorkflowRunStore{db},
		JobRunsStore:      &JobRunStore{db},
		StepRunsStore:     &StepRunStore{db},
	}
}
