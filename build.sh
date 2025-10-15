#!/bin/bash

# Build script for Kubernetes Monitor
set -e

APP_NAME="k8s-monitor"
IMAGE_NAME="${APP_NAME}"
IMAGE_TAG="${1:-latest}"
REGISTRY="${REGISTRY:-}"

# Version information
APP_VERSION="${APP_VERSION:-v1.0.0}"
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "🚀 Building ${APP_NAME} application..."
echo "   Version: ${APP_VERSION}"
echo "   Build Date: ${BUILD_DATE}"
echo "   Git Commit: ${GIT_COMMIT}"
echo ""

# Build the Docker image
if [ -n "$REGISTRY" ]; then
    FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
else
    FULL_IMAGE_NAME="${IMAGE_NAME}:${IMAGE_TAG}"
fi

echo "Building image: ${FULL_IMAGE_NAME}"
docker build \
    --build-arg APP_VERSION="${APP_VERSION}" \
    --build-arg BUILD_DATE="${BUILD_DATE}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    -t ${FULL_IMAGE_NAME} .

echo "✅ Build completed successfully!"
echo "📦 Image: ${FULL_IMAGE_NAME}"

# Optional: Push to registry if REGISTRY is set
if [ -n "$REGISTRY" ]; then
    echo ""
    echo "📤 Pushing to registry..."
    docker push ${FULL_IMAGE_NAME}
    echo "✅ Push completed!"
fi

echo ""
echo "🚀 Deployment Options:"
echo ""
echo "🐳 Docker:"
echo "   docker run -p 8080:8080 -v ~/.kube/config:/home/app/.kube/config:ro ${FULL_IMAGE_NAME}"
echo ""
echo "🐙 Docker Compose:"
echo "   docker-compose up -d"
echo ""
echo "☸️  Kubernetes:"
echo "   kubectl apply -f k8s/"
echo "   kubectl port-forward service/k8s-monitor 8080:80"
echo ""
echo "🌐 Access the application:"
echo "   http://localhost:8080"
