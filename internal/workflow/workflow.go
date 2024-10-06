package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/storage"
	"golang.org/x/crypto/ssh"
)

type WorkflowManager struct {
	store          *storage.Storage
	githubClientId string
	cache          *redis.Client
}

type WorkflowRun struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

func (w WorkflowRun) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

func New(s *storage.Storage, c string, r *redis.Client) *WorkflowManager {
	return &WorkflowManager{store: s, githubClientId: c, cache: r}
}

func (wm *WorkflowManager) TriggerWorkflow(ctx context.Context, workflowId string) error {
	workflows, err := wm.getWorkflowsByID(workflowId)
	if err != nil {
		return err
	}

	jobMap := wm.partitionWorkflowsByJob(workflows)

	nodes, err := wm.store.Nodes.All()
	if err != nil {
		return err
	}

	client, err := wm.initializeSSH(nodes[0])
	if err != nil {
		return err
	}

	for jobId, steps := range jobMap {
		fmt.Printf("Processing job: %s\n", jobId)

		if err := wm.updateJobStatus(ctx, jobId, "running"); err != nil {
			slog.Error("error updating workflow status: %w", "err", err)
		}

		if err := wm.executeJob(ctx, client, steps, nodes[0]); err != nil {
			if err := wm.updateJobStatus(ctx, jobId, "failed"); err != nil {
				slog.Error("error updating workflow status: %w", "err", err)
			}
			return err
		}
		if err := wm.updateJobStatus(ctx, jobId, "success"); err != nil {
			slog.Error("error updating workflow status: %w", "err", err)
		}
	}

	return nil
}

func (wm *WorkflowManager) updateJobStatus(ctx context.Context, jobId, status string) error {
	eventPayload := WorkflowRun{Id: jobId, Status: status, Type: "job"}
	if err := wm.cache.Publish(ctx, "job_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := wm.store.JobRunsStore.UpdateStatus(jobId, status); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) updateStepStatus(ctx context.Context, stepId, status string) error {
	eventPayload := WorkflowRun{Id: stepId, Status: status, Type: "step"}
	if err := wm.cache.Publish(ctx, "step_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := wm.store.StepRunsStore.UpdateStatus(stepId, status); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) getWorkflowsByID(workflowId string) ([]storage.WorkflowRunSteps, error) {
	workflows, err := wm.store.WorkflowRunsStore.GetById(workflowId)
	if err != nil {
		return nil, err
	}

	if len(workflows) == 0 {
		return nil, fmt.Errorf("no workflows found for workflow ID: %s", workflowId)
	}

	return workflows, nil
}

func (wm *WorkflowManager) partitionWorkflowsByJob(workflows []storage.WorkflowRunSteps) map[string][]storage.WorkflowRunSteps {
	jobMap := make(map[string][]storage.WorkflowRunSteps)
	for _, workflow := range workflows {
		jobMap[workflow.JobID] = append(jobMap[workflow.JobID], workflow)
	}
	return jobMap
}

func (wm *WorkflowManager) initializeSSH(node storage.Node) (*ssh.Client, error) {
	sshConn, err := remote.New(node.Host, node.User, []byte(node.PemFile), true)
	if err != nil {
		return nil, err
	}

	client, err := sshConn.Dial()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (wm *WorkflowManager) executeJob(ctx context.Context, client *ssh.Client, steps []storage.WorkflowRunSteps, node storage.Node) error {
	// Shared data for the job
	gitUrl := steps[0].Url
	repoName := steps[0].RepoName
	ref := steps[0].Branch

	defer func() {
		go wm.stopDockerImage(client, repoName)
	}()

	for _, step := range steps {
		if err := wm.updateStepStatus(ctx, step.StepID, "running"); err != nil {
			slog.Error("error updating step status: %w", "err", err)
		}
		if err := wm.executeStep(client, step, gitUrl, repoName, ref, node); err != nil {
			if err := wm.updateStepStatus(ctx, step.StepID, "failed"); err != nil {
				slog.Error("error updating step status: %w", "err", err)
			}
			return err
		}
		if err := wm.updateStepStatus(ctx, step.StepID, "success"); err != nil {
			slog.Error("error updating step status: %w", "err", err)
		}
	}

	return nil
}

func (wm *WorkflowManager) executeStep(client *ssh.Client, step storage.WorkflowRunSteps, gitUrl, repoName, ref string, node storage.Node) error {
	switch step.Type {
	case "checkout":
		if err := wm.pullDockerImage(client, step.Docker); err != nil {
			return err
		}

		if err := wm.runBackgroundDockerImage(client, step.Docker, repoName); err != nil {
			return err
		}

		if err := wm.checkout(client, node.PemFile, gitUrl, repoName, ref); err != nil {
			return err
		}

	case "restore_cache":
		fmt.Printf("step %s (restore_cache) not implemented\n", step.StepID)

	case "save_cache":
		fmt.Printf("step %s (save_cache) not implemented\n", step.StepID)

	default:
		cmd := fmt.Sprintf(`
			docker exec -w /%s %s sh -c "%s"
		`, repoName, repoName, step.Command)
		if err := wm.runStepCommand(client, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (wm *WorkflowManager) saveCommandOutput(streamType string, output []byte) {
	fmt.Printf("[%s] %s", streamType, output)
}

func (wm *WorkflowManager) runStepCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker pull: %w", err)
	}
	defer session.Close()

	if err := remote.RunCommand(session, command, wm.saveCommandOutput); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) checkout(client *ssh.Client, pemContent, gitUrl, name, ref string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for checkout: %w", err)
	}
	defer session.Close()

	token, err := remote.GenerateJWTToken([]byte(pemContent), wm.githubClientId)
	if err != nil {
		return err
	}

	cloneUrl, err := convertGitHubURL(gitUrl, token)
	if err != nil {
		return err
	}

	command := fmt.Sprintf(`
    docker exec %s sh -c "
        pwd
        git clone %s
        cd %s && git fetch origin && git checkout %s
		git log -1
    "
	`, name, cloneUrl, name, ref)

	if err := remote.RunCommand(session, command, wm.saveCommandOutput); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) pullDockerImage(client *ssh.Client, docker string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker pull: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
        docker pull %s
    `, docker)

	if err := remote.RunCommand(session, command, wm.saveCommandOutput); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) runBackgroundDockerImage(client *ssh.Client, imageName, repoName string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker run: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
        docker run -d -v /home/ubuntu/%s:/app/%s --name %s %s tail -f /dev/null
    `, repoName, repoName, repoName, imageName)

	if err := remote.RunCommand(session, command, wm.saveCommandOutput); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) stopDockerImage(client *ssh.Client, repoName string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker stop: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
		if [ "$(docker ps -q -f name=%s)" ]; then
			docker stop %s
		fi
		if [ "$(docker ps -aq -f name=%s)" ]; then
			docker rm %s
		fi
    `, repoName, repoName, repoName, repoName)

	if err := remote.RunCommand(session, command, wm.saveCommandOutput); err != nil {
		return err
	}

	return nil
}

func convertGitHubURL(originalURL string, token string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitHub URL: %v", err)
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", fmt.Errorf("invalid GitHub URL format")
	}
	owner, repo := pathParts[0], pathParts[1]

	convertedURL := fmt.Sprintf("https://x-access-token:%s@%s/%s/%s.git", token, parsedURL.Host, owner, repo)
	return convertedURL, nil
}
