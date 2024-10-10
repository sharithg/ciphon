#!/bin/bash

# Check if the image name is passed as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <image-name>"
  exit 1
fi

IMAGE_NAME=$1

# Get all container IDs using the specified image
CONTAINERS=$(docker ps -a --filter "ancestor=$IMAGE_NAME" --format "{{.ID}}")

if [ -z "$CONTAINERS" ]; then
  echo "No containers found for image: $IMAGE_NAME"
  exit 0
fi

# Stop all containers using the specified image
echo "Stopping containers using image: $IMAGE_NAME"
docker stop $CONTAINERS

# Remove all containers using the specified image
echo "Removing containers using image: $IMAGE_NAME"
docker rm $CONTAINERS

echo "All containers for image $IMAGE_NAME have been stopped and removed."
