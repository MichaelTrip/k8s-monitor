#!/bin/bash

# Set user and group IDs to match current user
export UID=$(id -u)
export GID=$(id -g)

# Create data directory if it doesn't exist
mkdir -p data

# Make sure data directory has correct permissions
chmod 755 data

# Ensure kubeconfig is readable
if [ -f ~/.kube/config ]; then
    echo "Making kubeconfig readable..."
    chmod 644 ~/.kube/config
    chmod 755 ~/.kube
fi

echo "Starting k8s-monitor with user ID: $UID, group ID: $GID"
docker compose up "$@"
