package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v65/github"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/sharithg/siphon/internal/protogen"
)

type WorkflowManager struct {
	store          *repository.Queries
	githubClientId string
	cache          *redis.Client
	service        *service.Service
	Client         *github.Client
}

type WorkflowRun struct {
	Id     uuid.UUID                     `json:"id"`
	Status repository.WorkflowStatusEnum `json:"status"`
	Type   string                        `json:"type"`
}

func (w WorkflowRun) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

func New(s *repository.Queries, c string, r *redis.Client, serv *service.Service, ghClient *github.Client) *WorkflowManager {
	return &WorkflowManager{store: s, githubClientId: c, cache: r, service: serv, Client: ghClient}
}

func (wm *WorkflowManager) TriggerWorkflow(ctx context.Context, workflowId uuid.UUID) error {
	jobMap, workflow, dag, err := wm.service.Job.GetJobsAndStepsByWorkflowId(ctx, workflowId)

	if err != nil {
		return err
	}

	executionOrder, err := dag.GetExecutionGroups()

	if err != nil {
		return err
	}

	nodes, err := wm.store.GetAllNodes(ctx)

	if err != nil {
		return err
	}

	if err = wm.updateGithubStatus(ctx, workflow, repository.WorkflowStatusEnumRunning); err != nil {
		slog.Error("error updating git status", "err", err)
	}

	if err = wm.executeJobs(ctx, executionOrder, jobMap, nodes[0]); err != nil {
		if errUpdate := wm.updateGithubStatus(ctx, workflow, repository.WorkflowStatusEnumFailed); errUpdate != nil {
			slog.Error("error updating git status", "err", errUpdate)
		}
		slog.Error("error executing jobs", "err", err)
	}

	if err = wm.updateGithubStatus(ctx, workflow, repository.WorkflowStatusEnumSuccess); err != nil {
		slog.Error("error updating git status", "err", err)
	}

	return nil
}

func (wm *WorkflowManager) executeJobs(ctx context.Context, executionOrder [][]uuid.UUID, jobMap map[uuid.UUID][]repository.GetJobsAndStepsByWorkflowIdRow, node repository.GetAllNodesRow) error {
	for _, jobs := range executionOrder {
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		errChan := make(chan error, 1)

		for _, jobId := range jobs {
			steps, ok := jobMap[jobId]

			if !ok {
				return fmt.Errorf("unexpected error, job id %s not found in map", jobId)
			}

			wg.Add(1)

			go func(ctx context.Context, jobId uuid.UUID, steps []repository.GetJobsAndStepsByWorkflowIdRow) {
				defer wg.Done()

				fmt.Printf("Processing job: %s\n", jobId)

				select {
				case <-ctx.Done():
					slog.Error("context canceled for job", "jobId", jobId)
					if err := wm.updateJobStatus(jobId, repository.JobStatusEnumFailed); err != nil {
						slog.Error("error updating job status 1", "err", err)
					}
					return
				default:
					if err := wm.updateJobStatus(jobId, repository.JobStatusEnumRunning); err != nil {
						slog.Error("error updating job status 2", "err", err)
						errChan <- err
						cancel()
						return
					}

					if err := wm.executeJob(ctx, steps, node); err != nil {
						if errUpdate := wm.updateJobStatus(jobId, repository.JobStatusEnumFailed); errUpdate != nil {
							slog.Error("error updating job status 3", "err", errUpdate)
						}
						slog.Error("error executing job", "err", err)
						errChan <- err
						cancel()
						return
					}

					if errUpdate := wm.updateJobStatus(jobId, repository.JobStatusEnumSuccess); errUpdate != nil {
						slog.Error("error updating job status 4", "err", errUpdate)
						errChan <- errUpdate
						cancel()
						return
					}
				}

			}(ctx, jobId, steps)
		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

		if err := <-errChan; err != nil {
			if je, ok := err.(*JobError); ok {
				if errUpdate := wm.updateJobStatus(je.JobId, repository.JobStatusEnumFailed); errUpdate != nil {
					return errUpdate
				}
				for _, stepId := range je.StepIds {
					if stepId == je.ErrorStep {
						if errUpdate := wm.updateStepStatus(stepId.String(), repository.StepStatusEnumFailed); errUpdate != nil {
							return errUpdate
						}
					} else {
						if errUpdate := wm.updateStepStatus(stepId.String(), repository.StepStatusEnumCancelled); errUpdate != nil {
							return errUpdate
						}
					}
				}
			} else {
				// err is not of type *JobError
				fmt.Println("Not a JobError")
			}
			return err
		}
	}

	return nil
}

func NewGrpcConnection(host string) (pb.StreamingServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:8888", host), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, nil, err
	}

	client := pb.NewStreamingServiceClient(conn)

	return client, conn, nil
}

func (wm *WorkflowManager) updateGithubStatus(ctx context.Context, workflow *repository.GetJobsAndStepsByWorkflowIdRow, status repository.WorkflowStatusEnum) error {
	var ghStatus string

	switch status {
	case "running":
		ghStatus = "pending"
	case "failed":
		ghStatus = "failure"
	case "success":
		ghStatus = "success"
	}

	statusUrl := fmt.Sprintf("http://localhost:5173/dashboard/pipelines/workflows/%s", workflow.WorkflowID)
	desc := "Ciphon running"

	ghStatusObj := github.RepoStatus{
		// error, failure, pending, success
		State:       &ghStatus,
		URL:         &statusUrl,
		TargetURL:   &statusUrl,
		Description: &desc,
	}

	_, _, err := wm.Client.Repositories.CreateStatus(ctx, workflow.Owner, workflow.RepoName, workflow.CommitSha, &ghStatusObj)

	if err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) updateJobStatus(jobId uuid.UUID, status repository.JobStatusEnum) error {
	if err := wm.store.UpdateJobRunStatus(context.Background(), repository.UpdateJobRunStatusParams{
		ID:     jobId,
		Status: status,
	}); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) updateStepStatus(stepId string, status repository.StepStatusEnum) error {
	ctx := context.Background()
	id, err := uuid.Parse(stepId)

	if err != nil {
		return err
	}

	if err := wm.store.UpdateStepRunStatus(ctx, repository.UpdateStepRunStatusParams{
		Status: status,
		ID:     id,
	}); err != nil {
		return err
	}
	return nil
}

func printCmdOutput(output *pb.CommandOutput) {
	slog.Info("received output", "cmdType", output.CmdType, "outputType", output.OutputType, "isUserCmd", output.IsUserCmd)
}

func (wm *WorkflowManager) executeJob(ctx context.Context, steps []repository.GetJobsAndStepsByWorkflowIdRow, node repository.GetAllNodesRow) error {

	var host string

	if len(steps) == 0 {
		return nil
	}

	if os.Getenv("GOENV") == "local" {
		host = "localhost"
	} else {
		host = node.Host
	}

	client, conn, err := NewGrpcConnection(host)

	defer conn.Close()

	if err != nil {
		return fmt.Errorf("error connecting to grpc: %s", err)
	}

	commands, err := wm.getSteps(steps, node)
	if err != nil {
		return err
	}

	var stepIds []uuid.UUID
	jobId := steps[0].JobID

	jobError := JobError{
		JobId:   jobId,
		StepIds: stepIds,
	}

	for _, step := range steps {
		stepIds = append(stepIds, step.StepID)
	}

	md := metadata.New(map[string]string{
		"authorization": node.AgentToken,
	})

	grpcctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.RunCommands(grpcctx, commands)

	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			slog.Error("job context cancelled")
			return fmt.Errorf("job context cancelled")
		default:

			output, err := stream.Recv()

			printCmdOutput(output)

			if err != nil {
				slog.Error("error reading grpc message", "err", err)
				return err
			}

			if output.CmdType == "done" {
				slog.Info("received done message, exiting")
				return nil
			} else if output.CmdType == "error" && output.IsUserCmd {
				if err = wm.saveCommandOutput(output.Id, output.OutputType, output.Output); err != nil {
					slog.Error("error updating command output status", "err", err)
				}
				if err := wm.updateStepStatus(output.Id, "failed"); err != nil {
					slog.Error("error updating step status", "err", err, "type", output.CmdType)
				}
				stId, err := uuid.Parse(output.Id)

				if err != nil {
					return err
				}

				jobError.ErrorStep = stId

				return &jobError
			} else if output.CmdType == "running" && output.IsUserCmd {
				if err := wm.updateStepStatus(output.Id, "running"); err != nil {
					slog.Error("error updating step status", "err", err, "type", output.CmdType)
				}
			} else if output.CmdType == "doneCmd" && output.IsUserCmd {
				if output.Id != "" {
					if err := wm.updateStepStatus(output.Id, "success"); err != nil {
						slog.Error("error updating step status", "err", err, "type", output.CmdType)
						return nil
					}
				}
			} else if output.CmdType == "cmd" {
				if err = wm.saveCommandOutput(output.Id, output.OutputType, output.Output); err != nil {
					slog.Error("error updating command output status", "err", err)
					return nil
				}
			} else {
				slog.Warn("received unknown command type", "CmdType", output.CmdType)
			}

		}

	}

}

func (wm *WorkflowManager) getSteps(steps []repository.GetJobsAndStepsByWorkflowIdRow, node repository.GetAllNodesRow) (*pb.Commands, error) {
	var commands []*pb.Command
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

			commands = append(commands, &pb.Command{
				Id:    step.StepID.String(),
				Cmd:   fmt.Sprintf("git clone %s && cd %s && git fetch origin && git checkout %s && git log -1", cloneUrl, step.RepoName, step.Branch),
				Order: step.StepOrder,
			})

			workDir = fmt.Sprintf("/%s", step.RepoName)

		case "restore_cache":
			commands = append(commands, &pb.Command{
				Id:      step.StepID.String(),
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		case "save_cache":
			commands = append(commands, &pb.Command{
				Id:      step.StepID.String(),
				Cmd:     "echo 'restore_cache'",
				Order:   step.StepOrder,
				WorkDir: workDir,
			})
		default:
			cmd := step.Command
			if cmd != nil {
				commands = append(commands, &pb.Command{
					Id:      step.StepID.String(),
					Cmd:     *cmd,
					Order:   step.StepOrder,
					WorkDir: workDir,
				})
			}
		}
	}

	payload := &pb.Commands{
		BaseEvent: &pb.BaseEvent{
			Type: "run_command",
		},
		Image:    steps[0].Docker,
		Commands: commands,
	}

	return payload, nil
}

func (wm *WorkflowManager) saveCommandOutput(stepId string, streamType, output string) error {

	ctx := context.Background()

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
