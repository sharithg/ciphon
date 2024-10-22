package grpcserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sort"

	"github.com/sharithg/siphon/agent/cli"
	"github.com/sharithg/siphon/agent/docker"
	pb "github.com/sharithg/siphon/internal/protogen"
	"github.com/sharithg/siphon/internal/utils"
	"google.golang.org/grpc"
)

type RunId string

const RunIdKey RunId = "runId"

type GrpcServer struct {
	Cli    *cli.Cli
	Docker *docker.Docker
	pb.UnimplementedStreamingServiceServer
}

func New(cli *cli.Cli, dock *docker.Docker) *GrpcServer {
	return &GrpcServer{
		Cli:    cli,
		Docker: dock,
	}
}

func (s GrpcServer) RunCommands(req *pb.Commands, srv grpc.ServerStreamingServer[pb.CommandOutput]) error {

	ctx := context.Background()
	runId := utils.RandStringBytes(10)

	runId = fmt.Sprintf("ciphon-run-%s", runId)

	logger := slog.With("runId", runId)

	logger.Info("Starting Run")

	outputChan := make(chan *pb.CommandOutput)

	stdoutSetupFunc := stdoutHandler(outputChan, "setup", "")
	stderrSetupFunc := stderrHandler(outputChan, "setup", "")

	defer s.teardown(srv, runId, logger)

	go func() {
		defer close(outputChan)

		running := pb.CommandOutput{
			CmdType:   "running",
			IsUserCmd: false,
		}

		if err := sendOutput(srv, &running); err != nil {
			logger.Error("error sending output over rpc", "err", err)
			return
		}

		if err := s.Cli.DockerPullImageAndStreamOutput(ctx, req.Image, stdoutSetupFunc, stderrSetupFunc); err != nil {
			logger.Error("error pulling docker image", "err", err)
			sendError(srv, err, false, "", logger)
			return
		}

		if err := sendOutput(srv, &running); err != nil {
			logger.Error("error sending output over rpc", "err", err)
			return
		}

		if err := s.Cli.DockerRunBackgroundContainer(runId, req.Image, stdoutSetupFunc, stderrSetupFunc); err != nil {
			logger.Error("error running background container", "err", err)
			sendError(srv, err, false, "", logger)
			return
		}

		doneCmd := pb.CommandOutput{
			CmdType:   "doneCmd",
			IsUserCmd: false,
		}
		if err := sendOutput(srv, &doneCmd); err != nil {
			logger.Error("error sending output over rpc", "err", err)
			return
		}

		commands := req.Commands

		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Order < commands[j].Order
		})

		for _, cmd := range commands {

			running := pb.CommandOutput{
				CmdType:   "running",
				Id:        cmd.Id,
				IsUserCmd: true,
			}

			stdoutCmdFunc := stdoutHandler(outputChan, "cmd", cmd.Id)
			stderrCmdFunc := stderrHandler(outputChan, "cmd", cmd.Id)

			workDir := "/"

			if cmd.WorkDir != "" {
				workDir = cmd.WorkDir
			}

			if err := sendOutput(srv, &running); err != nil {
				logger.Error("error sending output over rpc", "err", err)
				return
			}

			if err := s.Cli.DockerExecAndStreamLogs(runId, workDir, cmd.Cmd, stdoutCmdFunc, stderrCmdFunc); err != nil {
				logger.Error("error running command", "err", err, "runId", runId)
				sendError(srv, err, true, cmd.Id, logger)
				return
			}
			doneCmd := pb.CommandOutput{
				CmdType:   "doneCmd",
				Id:        cmd.Id,
				IsUserCmd: true,
			}
			if err := sendOutput(srv, &doneCmd); err != nil {
				logger.Error("error sending output over rpc", "err", err)
				return
			}
		}

	}()

	for output := range outputChan {
		if err := sendOutput(srv, output); err != nil {
			logger.Error("error sending output over rpc", "err", err)
			return err
		}
	}

	return nil
}

func (s *GrpcServer) teardown(srv pb.StreamingService_RunCommandsServer, runId string, logger *slog.Logger) {

	logger.Info("teardown operation started")
	running := pb.CommandOutput{
		CmdType:   "running",
		IsUserCmd: false,
	}
	if err := sendOutput(srv, &running); err != nil {
		logger.Error("error sending output over rpc", "err", err)
	}

	containerId := []string{runId}

	if err := s.Docker.CleanUpContainers(containerId); err != nil {
		logger.Error("error stopping and deleting containers", "err", err)
		sendError(srv, err, false, "", logger)
		return
	}

	doneCmd := pb.CommandOutput{
		CmdType:   "doneCmd",
		IsUserCmd: false,
	}
	if err := sendOutput(srv, &doneCmd); err != nil {
		logger.Error("error sending output over rpc", "err", err)
	}

	output := pb.CommandOutput{
		CmdType: "done",
	}

	if err := sendOutput(srv, &output); err != nil {
		logger.Error("error sending output over rpc", "err", err)
	}
}

func stdoutHandler(outputChan chan *pb.CommandOutput, cmdType string, cmdId string) func(string) {
	return func(message string) {
		outputChan <- &pb.CommandOutput{
			OutputType: "stdout",
			Output:     message,
			CmdType:    cmdType,
			Id:         cmdId,
		}
	}
}

func stderrHandler(outputChan chan *pb.CommandOutput, cmdType string, cmdId string) func(string) {
	return func(message string) {
		outputChan <- &pb.CommandOutput{
			OutputType: "stderr",
			Output:     message,
			CmdType:    cmdType,
			Id:         cmdId,
		}
	}
}

func sendOutput(srv pb.StreamingService_RunCommandsServer, output *pb.CommandOutput) error {
	if err := srv.Send(output); err != nil {
		log.Println("error generating response")
		return err
	}
	return nil
}

func sendError(srv pb.StreamingService_RunCommandsServer, errMsg error, isUserCmd bool, cmdId string, logger *slog.Logger) error {
	err := srv.Send(&pb.CommandOutput{
		CmdType:   "error",
		Output:    errMsg.Error(),
		IsUserCmd: isUserCmd,
		Id:        cmdId,
	})
	if err != nil {
		return fmt.Errorf("can't send: %s", err.Error())
	}

	output := pb.CommandOutput{
		CmdType: "done",
	}

	if err := sendOutput(srv, &output); err != nil {
		logger.Error("error sending output over rpc", "err", err)
	}

	return nil
}
