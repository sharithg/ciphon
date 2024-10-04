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
		for _, jobName := range workflow.Jobs {
			if _, exists := c.Jobs[jobName]; !exists {
				return fmt.Errorf("workflow '%s' references undefined job '%s'", workflowName, jobName)
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
	for _, jobName := range workflow.Jobs {
		job, exists := c.Jobs[jobName]
		if !exists {
			return nil, fmt.Errorf("job '%s' not found in workflow '%s'", jobName, workflowName)
		}
		jobs = append(jobs, JobWithName{
			Job:  job,
			Name: jobName,
		})
	}

	return jobs, nil
}
