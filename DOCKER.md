# Docker and OpenShift Deployment Guide

## Building the Docker Image

### Prerequisites
- Docker installed and running
- Access to a container registry (for OpenShift deployment)

### Build Process

1. **Build the image locally:**
   ```bash
   ./build.sh
   ```

2. **Build with custom tag:**
   ```bash
   ./build.sh v1.0.0
   ```

3. **Build and push to registry:**
   ```bash
   REGISTRY=your-registry.com/namespace ./build.sh v1.0.0
   ```

## OpenShift Deployment

### Quick Deployment

1. **Create a new project:**
   ```bash
   oc new-project k8s-monitor
   ```

2. **Deploy using the template:**
   ```bash
   oc process -f openshift-template.yaml \
     -p IMAGE_NAME=k8s-monitor \
     -p IMAGE_TAG=latest \
     -p NAMESPACE=k8s-monitor | oc apply -f -
   ```

3. **Expose the service:**
   ```bash
   oc expose service k8s-monitor
   ```

4. **Get the route URL:**
   ```bash
   oc get route k8s-monitor
   ```

### Security Features

The Dockerfile and deployment template follow OpenShift security best practices:

- ✅ **Non-root user**: Runs as user ID 1001
- ✅ **Read-only root filesystem**: Container filesystem is read-only except for data volume
- ✅ **No privilege escalation**: `allowPrivilegeEscalation: false`
- ✅ **Dropped capabilities**: All Linux capabilities are dropped
- ✅ **Resource limits**: CPU and memory limits are defined
- ✅ **Health checks**: Liveness and readiness probes configured
- ✅ **Minimal base image**: Uses UBI minimal for smaller attack surface
- ✅ **RBAC**: Least-privilege access to Kubernetes resources

### Configuration

The application configuration can be customized by editing the ConfigMap in the template:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: k8s-monitor-config
data:
  config.json: |
    {
      "webPort": 8080,
      "persistence": {
        "enabled": true,
        "filePath": "/app/data/changes.json",
        "autoSave": true,
        "saveInterval": 30
      },
      "resources": [
        # Add or remove resources to monitor
      ]
    }
```

### Persistent Storage

If you need persistent storage for the changes data:

```bash
oc set volume deployment/k8s-monitor \
  --add --name=data-volume \
  --type=persistentVolumeClaim \
  --claim-size=1Gi \
  --mount-path=/app/data \
  --overwrite
```

### Environment Variables

Available environment variables for runtime configuration:
- `WEB_PORT`: Port for the web server (default: 8080)
- `PERSISTENCE_FILE_PATH`: Path for the changes JSON file (default: changes.json)
- `KUBECONFIG`: Path to kubeconfig file (uses in-cluster config by default)

**Example with custom environment variables:**
```bash
docker run -d --name k8s-monitor \
  -p 9090:9090 \
  -v ~/.kube:/home/k8s-monitor/.kube:ro \
  -v ./data:/app/data \
  -e KUBECONFIG=/home/k8s-monitor/.kube/config \
  -e WEB_PORT=9090 \
  -e PERSISTENCE_FILE_PATH=/app/data/my-k8s-changes.json \
  --user $(id -u):$(id -g) \
  k8s-monitor:latest
```

**Note:** Environment variables take precedence over config.json settings for supported parameters.

### Troubleshooting

1. **Check pod logs:**
   ```bash
   oc logs deployment/k8s-monitor
   ```

2. **Check pod status:**
   ```bash
   oc get pods -l app=k8s-monitor
   ```

3. **Check service account permissions:**
   ```bash
   oc auth can-i list pods --as=system:serviceaccount:k8s-monitor:k8s-monitor
   ```

## Local Development with Docker

### Run locally with Docker:

```bash
# Build the image
docker build -t k8s-monitor .

# Run with local kubeconfig
docker run -it --rm \
  -p 8080:8080 \
  -v ~/.kube/config:/tmp/kubeconfig:ro \
  -e KUBECONFIG=/tmp/kubeconfig \
  k8s-monitor
```

### Run with docker-compose (optional):

Create a `docker-compose.yml`:

```yaml
version: '3.8'
services:
  k8s-monitor:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ~/.kube/config:/tmp/kubeconfig:ro
      - ./data:/app/data
    environment:
      - KUBECONFIG=/tmp/kubeconfig
    restart: unless-stopped
```

Then run:
```bash
docker-compose up
```
