package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseConfig(data string) (*Config, error) {
	t := Config{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Config) ValidateWorkflows() error {
	for workflowName, workflow := range c.Workflows {
		for _, job := range workflow.Jobs {
			if _, exists := c.Jobs[job.Name]; !exists {
				return fmt.Errorf("workflow '%s' references undefined job '%s'", workflowName, job.Name)
			}
		}
	}
	return nil
}

type JobWithName struct {
	Job
	Name string
}

func (c *Config) GetWorkflowJobs(workflowName string) ([]JobWithName, error) {
	workflow, exists := c.Workflows[workflowName]
	if !exists {
		return nil, fmt.Errorf("workflow '%s' not found", workflowName)
	}

	var jobs []JobWithName
	for _, jobRun := range workflow.Jobs {
		job, exists := c.Jobs[jobRun.Name]
		if !exists {
			return nil, fmt.Errorf("job '%s' not found in workflow '%s'", jobRun.Name, workflowName)
		}
		jobs = append(jobs, JobWithName{
			Job:  job,
			Name: jobRun.Name,
		})
	}

	return jobs, nil
}

func (c *Config) GetJobRequires(jobName string) []string {
	for _, j := range c.Workflows {
		for _, job := range j.Jobs {
			if jobName == job.Name {
				return job.Requires
			}
		}
	}

	return nil
}
