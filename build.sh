#!/bin/bash

# Build script for OpenShift deployment
set -e

APP_NAME="k8s-monitor"
IMAGE_NAME="${APP_NAME}"
IMAGE_TAG="${1:-latest}"
REGISTRY="${REGISTRY:-}"

echo "Building ${APP_NAME} application..."

# Build the Docker image
if [ -n "$REGISTRY" ]; then
    FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
else
    FULL_IMAGE_NAME="${IMAGE_NAME}:${IMAGE_TAG}"
fi

echo "Building image: ${FULL_IMAGE_NAME}"
docker build -t ${FULL_IMAGE_NAME} .

echo "Build completed successfully!"
echo "Image: ${FULL_IMAGE_NAME}"

# Optional: Push to registry if REGISTRY is set
if [ -n "$REGISTRY" ]; then
    echo "Pushing to registry..."
    docker push ${FULL_IMAGE_NAME}
    echo "Push completed!"
fi

echo ""
echo "To deploy to OpenShift:"
echo "1. Create a new project: oc new-project k8s-monitor"
echo "2. Process template: oc process -f openshift-template.yaml -p IMAGE_NAME=${FULL_IMAGE_NAME%:*} -p IMAGE_TAG=${IMAGE_TAG} | oc apply -f -"
echo "3. Create route: oc expose service k8s-monitor"
