#!/bin/bash

sudo apt-get update && sudo apt-get install -y docker.io git

if ! getent group docker >/dev/null; then
    sudo groupadd docker
fi

if ! groups $USER | grep &>/dev/null '\bdocker\b'; then
    sudo usermod -aG docker $USER
    newgrp docker
fi

docker pull sharith/ciphon-agent

mkdir -p ~/.ciphon
echo '%s' >~/.ciphon/agent.json

IMAGE_NAME="sharith/ciphon-agent"

CONTAINER_ID=$(docker ps -q --filter "name=ciphon-agent")

if [ -n "$CONTAINER_ID" ]; then
    echo "A container with image $IMAGE_NAME is already running (Container ID: $CONTAINER_ID). Stopping it..."
    docker stop "$CONTAINER_ID"
else
    echo "No running container found for image $IMAGE_NAME. Starting a new container..."
fi

docker run --rm -d \
    -v ~/.ciphon/agent.json:/app/agent.json \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ciphon-agent \
    -e AGENT_CONFIG_PATH=/app/agent.json \
    -p 8888:8888 \
    --log-driver json-file \
    --log-opt max-size=10m \
    --log-opt max-file=3 \
    "$IMAGE_NAME"
