package docker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Docker struct{}

func New() *Docker {
	return &Docker{}
}

func (d *Docker) GetAllContainers(prefix string) ([]string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	var matchingContainers []string
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.HasPrefix(strings.TrimPrefix(name, "/"), prefix) {
				matchingContainers = append(matchingContainers, strings.TrimPrefix(name, "/"))
				break
			}
		}
	}

	return matchingContainers, nil
}

func (d *Docker) CleanUpContainers(containerIDs []string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	for _, containerID := range containerIDs {
		slog.Info("Stopping container...", "containerID", containerID)
		if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
			if client.IsErrNotFound(err) {
				slog.Warn("Container not found, skipping stop", "containerID", containerID)
			} else {
				slog.Error("Error stopping container %s: %v", "containerID", containerID, "err", err)
			}
		} else {
			slog.Info("Container stopped", "containerID", containerID)
		}

		slog.Info("Removing container", "containerID", containerID)
		removeOptions := container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}
		if err := cli.ContainerRemove(ctx, containerID, removeOptions); err != nil {
			slog.Error("Error removing container", "containerID", containerID, "err", err)
		} else {
			slog.Info("Container removed", "containerID", containerID)
		}
	}

	slog.Info("Cleaning up dangling volumes...")
	pruneFilters := filters.NewArgs()
	volumes, err := cli.VolumesPrune(ctx, pruneFilters)
	if err != nil {
		return fmt.Errorf("error pruning volumes: %v", err)
	}
	slog.Info("Reclaimed volumes and space", "numVolumes", len(volumes.VolumesDeleted), "space", volumes.SpaceReclaimed)

	return nil
}
