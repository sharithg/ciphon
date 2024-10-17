package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/repository"
)

type JobService struct {
	repository *repository.Queries
}

func (j *JobService) GetByWorkflowId(ctx context.Context, id uuid.UUID) ([]repository.GetJobsByWorkflowIdRow, *DAG, error) {
	jobs, err := j.repository.GetJobsByWorkflowId(ctx, id)

	if err != nil {
		return nil, nil, err
	}

	jobDag := setupJobDAG(jobs)

	return jobs, jobDag, nil
}

func (j *JobService) GetJobsAndStepsByWorkflowId(ctx context.Context, id uuid.UUID) (map[uuid.UUID][]repository.GetJobsAndStepsByWorkflowIdRow, *DAG, error) {
	workflows, err := j.repository.GetJobsAndStepsByWorkflowId(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	jobMap := j.partitionWorkflowsByJob(workflows)

	_, dag, err := j.GetByWorkflowId(ctx, id)

	if err != nil {
		return nil, nil, err
	}

	return jobMap, dag, nil
}

func (j *JobService) partitionWorkflowsByJob(workflows []repository.GetJobsAndStepsByWorkflowIdRow) map[uuid.UUID][]repository.GetJobsAndStepsByWorkflowIdRow {
	jobMap := make(map[uuid.UUID][]repository.GetJobsAndStepsByWorkflowIdRow)
	for _, workflow := range workflows {
		jobMap[workflow.JobID] = append(jobMap[workflow.JobID], workflow)
	}
	return jobMap
}

func getJobRequires(jobs []repository.GetJobsByWorkflowIdRow) map[uuid.UUID][]uuid.UUID {
	jobMap := make(map[uuid.UUID][]uuid.UUID)
	for _, job := range jobs {
		jobMap[job.ID] = append(jobMap[job.ID], job.Requires...)
	}
	return jobMap
}

func setupJobDAG(jobs []repository.GetJobsByWorkflowIdRow) *DAG {
	jobDag := NewDAG()

	jobRequires := getJobRequires(jobs)

	for jobId := range jobRequires {
		jobDag.AddNode(jobId)
	}
	for jobId, requires := range jobRequires {
		for _, req := range requires {
			jobDag.AddEdge(req, jobId)
		}
	}

	return jobDag
}
