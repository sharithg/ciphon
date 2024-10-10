package remote

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sharithg/siphon/internal/config"
)

func generateYamlConfig(token string) (string, error) {
	cfg := config.AgentConfig{
		Token: token,
	}

	val, err := json.Marshal(cfg)

	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (s *SshConn) InstallTools(token string) error {
	start := time.Now()
	client, err := s.Dial()

	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	config, err := generateYamlConfig(token)

	if err != nil {
		return fmt.Errorf("failed to make agent config: %w", err)
	}

	command := fmt.Sprintf(`
		sudo apt-get update && \
		sudo apt-get install -y docker.io git && \
		
		if ! getent group docker >/dev/null; then
			sudo groupadd docker
		fi

		if ! groups $USER | grep &>/dev/null '\bdocker\b'; then
			sudo usermod -aG docker $USER
		fi

		newgrp docker

		docker pull sharith/ciphon-agent && \
		mkdir -p ~/.ciphon && \
		echo '%s' > ~/.ciphon/agent.json

		IMAGE_NAME="sharith/ciphon-agent"

		CONTAINER_ID=$(docker ps -q --filter "name=ciphon-agent")

		if [ -n "$CONTAINER_ID" ]; then
		echo "A container with image $IMAGE_NAME is already running (Container ID: $CONTAINER_ID)."
		docker stop $CONTAINER_ID
		docker run --rm -d -v ~/.ciphon/agent.json:/app/agent.json \
				   --name ciphon-agent \
				   -v /var/run/docker.sock:/var/run/docker.sock \
				   -e AGENT_CONFIG_PATH=/app/agent.json \
				   -p 8888:8888 \
				   sharith/ciphon-agent
		else
		echo "No running container found for image $IMAGE_NAME. Starting a new container..."
		docker run --rm -d -v ~/.ciphon/agent.json:/app/agent.json \
		            -v /var/run/docker.sock:/var/run/docker.sock \
					--name ciphon-agent \
					-e AGENT_CONFIG_PATH=/app/agent.json \
					-p 8888:8888 \
					$IMAGE_NAME
		fi
	`, config)

	fmt.Printf("dial took %v\n", time.Since(start))

	if err := RunCommand(session, command, func(streamType string, buf []byte) {
		fmt.Printf("[%s] %s", streamType, string(buf))
	}); err != nil {
		return err
	}

	log.Printf("Tools installed successfully on %s", s.Host)
	return nil
}
