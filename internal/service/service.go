package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/repository"
)

type Service struct {
	Job interface {
		GetByWorkflowId(ctx context.Context, id uuid.UUID) ([]repository.GetJobsByWorkflowIdRow, *DAG, error)
		GetJobsAndStepsByWorkflowId(ctx context.Context, id uuid.UUID) (map[uuid.UUID][]repository.GetJobsAndStepsByWorkflowIdRow, *DAG, error)
	}
}

func NewService(r *repository.Queries) *Service {
	return &Service{
		Job: &JobService{repository: r},
	}
}
