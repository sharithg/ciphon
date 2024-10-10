package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/storage"
	"github.com/sharithg/siphon/ws"
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
	workflows, err := wm.getWorkflowsByID(ctx, workflowId)

	if err != nil {
		return err
	}

	jobMap := wm.partitionWorkflowsByJob(workflows)

	nodes, err := wm.store.Nodes.All(ctx)

	if err != nil {
		return err
	}

	for jobId, steps := range jobMap {
		fmt.Printf("Processing job: %s\n", jobId)

		if err := wm.updateJobStatus(ctx, jobId, "running"); err != nil {
			slog.Error("error updating workflow status: %w", "err", err)
		}

		if err := wm.executeJob(ctx, steps, nodes[0]); err != nil {
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
	if err := wm.store.JobRunsStore.UpdateStatus(ctx, jobId, status); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) updateStepStatus(ctx context.Context, stepId, status string) error {
	eventPayload := WorkflowRun{Id: stepId, Status: status, Type: "step"}
	if err := wm.cache.Publish(ctx, "step_run", eventPayload).Err(); err != nil {
		return err
	}

	if err := wm.store.StepRunsStore.UpdateStatus(ctx, stepId, status); err != nil {
		return err
	}
	return nil
}

func (wm *WorkflowManager) getWorkflowsByID(ctx context.Context, workflowId string) ([]storage.WorkflowRunSteps, error) {
	workflows, err := wm.store.WorkflowRunsStore.GetById(ctx, workflowId)
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

func (wm *WorkflowManager) executeJob(ctx context.Context, steps []storage.WorkflowRunSteps, node storage.Node) error {

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

	i := 0

	for {
		i = i + 1
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

			var output ws.CommandOutput
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

func (wm *WorkflowManager) getSteps(steps []storage.WorkflowRunSteps, node storage.Node) (*ws.Commands, error) {
	var commands []ws.Command
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

			commands = append(commands, ws.Command{
				Id:    step.StepID,
				Cmd:   fmt.Sprintf("git clone %s && cd %s && git fetch origin && git checkout %s && git log -1", cloneUrl, step.RepoName, step.Branch),
				Order: step.StepOrder,
			})

			workDir = fmt.Sprintf("/%s", step.RepoName)

		case "restore_cache":
			commands = append(commands, ws.Command{
				Id:      step.StepID,
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		case "save_cache":
			commands = append(commands, ws.Command{
				Id:      step.StepID,
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		default:
			commands = append(commands, ws.Command{
				Id:      step.StepID,
				Cmd:     step.Command,
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		}
	}

	payload := &ws.Commands{
		BaseEvent: ws.BaseEvent{
			Type: "run_command",
		},
		Image:    steps[0].Docker,
		Commands: commands,
	}

	return payload, nil
}

func (wm *WorkflowManager) saveCommandOutput(ctx context.Context, stepId string, streamType, output string) error {

	cmd := storage.CommandOutput{
		StepID: stepId,
		Type:   &streamType,
		Stdout: output,
	}
	_, err := wm.store.StepRunsStore.CreateCommandOutput(ctx, cmd)

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
