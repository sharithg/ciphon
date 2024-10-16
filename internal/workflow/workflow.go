package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/runner"
)

type WorkflowManager struct {
	store          *repository.Queries
	githubClientId string
	cache          *redis.Client
}

type WorkflowRun struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"`
	Type   string    `json:"type"`
}

func (w WorkflowRun) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

func New(s *repository.Queries, c string, r *redis.Client) *WorkflowManager {
	return &WorkflowManager{store: s, githubClientId: c, cache: r}
}

func (wm *WorkflowManager) TriggerWorkflow(ctx context.Context, workflowId uuid.UUID) error {
	workflows, err := wm.getWorkflowsByID(ctx, workflowId)

	if err != nil {
		return err
	}

	jobMap := wm.partitionWorkflowsByJob(workflows)

	nodes, err := wm.store.GetAllNodes(ctx)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for jobId, steps := range jobMap {
		wg.Add(1)

		go func(jobId uuid.UUID, steps []repository.GetWorkflowRunByIdRow) {
			defer wg.Done()

			fmt.Printf("Processing job: %s\n", jobId)

			if err := wm.updateJobStatus(ctx, jobId, "running"); err != nil {
				slog.Error("error updating workflow status: %w", "err", err)
				return
			}

			if err := wm.executeJob(ctx, steps, nodes[0]); err != nil {
				if err := wm.updateJobStatus(ctx, jobId, "failed"); err != nil {
					slog.Error("error updating workflow status: %w", "err", err)
				}
				slog.Error("error executing job: %w", "err", err)
				return
			}

			if err := wm.updateJobStatus(ctx, jobId, "success"); err != nil {
				slog.Error("error updating workflow status: %w", "err", err)
			}
		}(jobId, steps)
	}

	wg.Wait()

	return nil
}

func (wm *WorkflowManager) updateJobStatus(ctx context.Context, jobId uuid.UUID, status string) error {
	eventPayload := WorkflowRun{Id: jobId, Status: status, Type: "job"}
	if err := wm.cache.Publish(ctx, "job_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := wm.store.UpdateJobRunStatus(ctx, repository.UpdateJobRunStatusParams{
		ID:     jobId,
		Status: &status,
	}); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) updateStepStatus(ctx context.Context, stepId string, status string) error {

	id, err := uuid.Parse(stepId)

	if err != nil {
		return err
	}

	eventPayload := WorkflowRun{Id: id, Status: status, Type: "step"}
	if err := wm.cache.Publish(ctx, "step_run", eventPayload).Err(); err != nil {
		return err
	}

	if err := wm.store.UpdateStepRunStatus(ctx, repository.UpdateStepRunStatusParams{
		Status: &status,
		ID:     id,
	}); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) getWorkflowsByID(ctx context.Context, workflowId uuid.UUID) ([]repository.GetWorkflowRunByIdRow, error) {
	workflows, err := wm.store.GetWorkflowRunById(ctx, workflowId)
	if err != nil {
		return nil, err
	}

	if len(workflows) == 0 {
		return nil, fmt.Errorf("no workflows found for workflow ID: %s", workflowId)
	}

	return workflows, nil
}

func (wm *WorkflowManager) partitionWorkflowsByJob(workflows []repository.GetWorkflowRunByIdRow) map[uuid.UUID][]repository.GetWorkflowRunByIdRow {
	jobMap := make(map[uuid.UUID][]repository.GetWorkflowRunByIdRow)
	for _, workflow := range workflows {
		jobMap[workflow.JobID] = append(jobMap[workflow.JobID], workflow)
	}
	return jobMap
}

func (wm *WorkflowManager) executeJob(ctx context.Context, steps []repository.GetWorkflowRunByIdRow, node repository.GetAllNodesRow) error {

	headers := http.Header{}
	headers.Set("X-Ciphon-Auth", node.AgentToken)

	fmt.Println("connecting: ", fmt.Sprintf("ws://%s:8888/ws", node.Host), node.AgentToken)
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:8888/ws", node.Host), headers)
	if err != nil {
		slog.Error("error connecting to agent", "err", err)
		return err
	}
	defer wsConn.Close()

	commands, err := wm.getSteps(steps, node)
	if err != nil {
		return err
	}

	cmdData, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	if err = wsConn.WriteMessage(websocket.TextMessage, cmdData); err != nil {
		return err
	}

	// Create a timeout context to avoid infinite looping if no messages arrive
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			slog.Error("timeout exceeded waiting for WebSocket messages")
			return fmt.Errorf("timeout exceeded")
		default:

			// Read a message from websocket connection
			_, msg, err := wsConn.ReadMessage()

			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Error("websocket connection closed")
				return fmt.Errorf("websocket connection closed")
			}

			if err != nil {
				slog.Error("error reading WebSocket message", "err", err)
				return err
			}

			var output runner.CommandOutput
			if err = json.Unmarshal(msg, &output); err != nil {
				slog.Error("error unmarshalling output", "error", err)
				return err
			}

			fmt.Printf("Received message: %s\n", string(msg))

			switch output.CmdType {
			case "done":
				slog.Info("received done message, exiting")
				return nil
			case "error":
				if err := wm.updateStepStatus(ctx, output.Id, "failed"); err != nil {
					slog.Error("error updating step status", "err", err)
				}
			case "running":
				if err := wm.updateStepStatus(ctx, output.Id, "running"); err != nil {
					slog.Error("error updating step status", "err", err)
				}
			case "doneCmd":
				if output.Id != "" {
					if err := wm.updateStepStatus(ctx, output.Id, "success"); err != nil {
						slog.Error("error updating step status", "err", err)
						return nil
					}
				}
			case "cmd":
				if err = wm.saveCommandOutput(ctx, output.Id, output.OutputType, output.Output); err != nil {
					slog.Error("error updating step status", "err", err)
					return nil
				}
			default:
				slog.Warn("received unknown command type", "CmdType", output.CmdType)
			}

		}

	}

}

func (wm *WorkflowManager) getSteps(steps []repository.GetWorkflowRunByIdRow, node repository.GetAllNodesRow) (*runner.Commands, error) {
	var commands []runner.Command
	workDir := "/"

	for _, step := range steps {

		switch step.Type {
		case "checkout":
			token, err := remote.GenerateJWTToken([]byte(node.PemFile), wm.githubClientId)
			if err != nil {
				return nil, err
			}
			cloneUrl, err := convertGitHubURL(step.Url, token)
			if err != nil {
				return nil, err
			}

			commands = append(commands, runner.Command{
				Id:    step.StepID,
				Cmd:   fmt.Sprintf("git clone %s && cd %s && git fetch origin && git checkout %s && git log -1", cloneUrl, step.RepoName, step.Branch),
				Order: step.StepOrder,
			})

			workDir = fmt.Sprintf("/%s", step.RepoName)

		case "restore_cache":
			commands = append(commands, runner.Command{
				Id:      step.StepID,
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		case "save_cache":
			commands = append(commands, runner.Command{
				Id:      step.StepID,
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		default:
			cmd := step.Command
			if cmd != nil {
				commands = append(commands, runner.Command{
					Id:      step.StepID,
					Cmd:     *cmd,
					Order:   step.StepOrder,
					WorkDir: workDir,
				})
			}
		}
	}

	payload := &runner.Commands{
		BaseEvent: runner.BaseEvent{
			Type: "run_command",
		},
		Image:    steps[0].Docker,
		Commands: commands,
	}

	return payload, nil
}

func (wm *WorkflowManager) saveCommandOutput(ctx context.Context, stepId string, streamType, output string) error {

	id, err := uuid.Parse(stepId)

	if err != nil {
		return err
	}

	cmd := repository.CreateCommandOutputParams{
		StepID: id,
		Type:   streamType,
		Stdout: output,
	}
	_, err = wm.store.CreateCommandOutput(ctx, cmd)

	if err != nil {
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
