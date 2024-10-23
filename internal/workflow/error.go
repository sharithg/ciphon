package workflow

import (
	"fmt"

	"github.com/google/uuid"
)

type JobError struct {
	JobId     uuid.UUID
	StepIds   []uuid.UUID
	ErrorStep uuid.UUID
}

func (m *JobError) Error() string {
	return fmt.Sprintf("error execting job: %s", m.JobId)
}

type StepCancelled struct {
	StepId string
}

func (m *StepCancelled) Error() string {
	return fmt.Sprintf("step cancelled: %s", m.StepId)
}
