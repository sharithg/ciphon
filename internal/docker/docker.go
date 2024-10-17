package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

type Docker struct{}

func New() (*Docker, error) {
	return &Docker{}, nil
}

func (d *Docker) RunBackgroundContainer(containerName, imageName string, stdoutHandler, stderrHandler func(string)) error {
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, imageName, "tail", "-f", "/dev/null")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	streamOutput(stdout, stdoutHandler)
	streamOutput(stderr, stderrHandler)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	return nil
}

func (d *Docker) ExecAndStreamLogs(containerName, workingDir, command string, stdoutHandler, stderrHandler func(string)) error {
	cmd := exec.Command("docker", "exec", "-w", workingDir, containerName, "sh", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	streamOutput(stdout, stdoutHandler)
	streamOutput(stderr, stderrHandler)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	return nil
}

func (d *Docker) PullImageAndStreamOutput(ctx context.Context, imageName string, stdoutHandler, stderrHandler func(string)) error {
	cmd := exec.CommandContext(ctx, "docker", "pull", imageName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start pulling image: %v", err)
	}

	streamOutput(stdout, stdoutHandler)
	streamOutput(stderr, stderrHandler)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	return nil
}

func (d *Docker) StopAndRemoveContainer(ctx context.Context, containerName string, stdoutHandler, stderrHandler func(string)) error {
	cmdStop := exec.CommandContext(ctx, "docker", "stop", containerName)

	stdoutStop, err := cmdStop.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout for stopping container: %v", err)
	}

	stderrStop, err := cmdStop.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr for stopping container: %v", err)
	}

	if err := cmdStop.Start(); err != nil {
		return fmt.Errorf("failed to start stopping container: %v", err)
	}

	streamOutput(stdoutStop, stdoutHandler)
	streamOutput(stderrStop, stderrHandler)

	if err := cmdStop.Wait(); err != nil {
		return fmt.Errorf("error while stopping container: %v", err)
	}

	cmdRemove := exec.CommandContext(ctx, "docker", "rm", "-v", containerName)

	stdoutRemove, err := cmdRemove.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout for removing container: %v", err)
	}

	stderrRemove, err := cmdRemove.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr for removing container: %v", err)
	}

	if err := cmdRemove.Start(); err != nil {
		return fmt.Errorf("failed to start removing container: %v", err)
	}

	streamOutput(stdoutRemove, stdoutHandler)
	streamOutput(stderrRemove, stderrHandler)

	if err := cmdRemove.Wait(); err != nil {
		return fmt.Errorf("error while removing container: %v", err)
	}

	return nil
}

func streamOutput(reader io.ReadCloser, handler func(string)) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		handler(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		handler(fmt.Sprintf("error reading output: %v", err))
	}
}
