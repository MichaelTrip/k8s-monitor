# K8s Monitor

K8s Monitor is a Go application that monitors changes in Kubernetes API objects and displays them in a beautiful web interface. This project provides real-time tracking of changes in your Kubernetes cluster.

## Features

- 🚀 **Real-time monitoring** of all major Kubernetes API objects (Pods, Deployments, Services, ConfigMaps, Secrets)
- ⚙️ **Configurable resources** - choose which Kubernetes resources to monitor via config.json
- 🎨 **Beautiful web interface** with auto-refresh functionality
- 📊 **Live statistics** showing total changes and uptime
- 🔍 **Detailed change tracking** with timestamps, event types, and resource details
- 💾 **Memory-efficient** - only tracks changes since program start (no historical dumps)
- 🎯 **Resource-specific watchers** for different Kubernetes objects
- 🌐 **RESTful API** for integration with other tools
- ✅ **Mark as read functionality** to manage change notifications
- 🔄 **Advanced filtering and sorting** by date, resource type, event type, and read status
- 💾 **Persistent storage** - save changes to JSON file with configurable auto-save

## Configuration

The application uses a `config.json` file to configure which resources to monitor and other settings:

```json
{
  "webPort": 8080,
  "persistence": {
    "enabled": true,
    "filePath": "changes.json",
    "autoSave": true,
    "saveInterval": 30
  },
  "resources": [
    {
      "name": "pods",
      "enabled": true,
      "namespace": "",
      "description": "Kubernetes Pods"
    },
    {
      "name": "deployments", 
      "enabled": true,
      "description": "Kubernetes Deployments"
    }
  ]
}
```

### Supported Resources:
- **pods** - Kubernetes Pods
- **deployments** - Kubernetes Deployments  
- **services** - Kubernetes Services
- **configmaps** - Kubernetes ConfigMaps
- **secrets** - Kubernetes Secrets
- **replicasets** - Kubernetes ReplicaSets
- **daemonsets** - Kubernetes DaemonSets
- **statefulsets** - Kubernetes StatefulSets
- **jobs** - Kubernetes Jobs
- **cronjobs** - Kubernetes CronJobs
- **persistentvolumes** - Kubernetes PersistentVolumes
- **persistentvolumeclaims** - Kubernetes PersistentVolumeClaims
- **ingresses** - Kubernetes Ingresses
- **networkpolicies** - Kubernetes NetworkPolicies

### Configuration Options:
- `webPort`: Port for the web interface (default: 8080)
- `persistence.enabled`: Enable/disable saving changes to file
- `persistence.filePath`: Path to the JSON file for saving changes
- `persistence.autoSave`: Automatically save changes at regular intervals
- `persistence.saveInterval`: Auto-save interval in seconds
- `logging.enabled`: Master switch for all logging (default: false)
- `logging.logChanges`: Log individual change events to stdout (default: false)
- `logging.logOperations`: Log save/load operations to stdout (default: false)
- `resources[].enabled`: Whether to monitor this resource type
- `resources[].namespace`: Specific namespace to monitor (empty = all namespaces)

### Environment Variables:
The following environment variables can override configuration settings:
- `PERSISTENCE_FILE_PATH`: Override the path for the changes JSON file (e.g., `/app/data/changes.json`)
- `WEB_PORT`: Override the web server port (e.g., `8080`)
- `KUBECONFIG`: Path to Kubernetes configuration file

## Screenshots

The web interface provides:
- Real-time change feed with color-coded event types (Added/Modified/Deleted)
- Statistics dashboard showing uptime and total changes
- Filterable and sortable change history
- Auto-refresh functionality with manual refresh option

## Prerequisites

- Go 1.16 or later
- Access to a Kubernetes cluster
- `kubectl` configured to communicate with your cluster

## Installation

### Option 1: Local Go Installation

1. Clone the repository:

   ```
   git clone <repository-url>
   cd k8s-monitor
   ```

2. Install the dependencies:

   ```
   go mod tidy
   ```

### Option 2: Docker Installation (Recommended)

1. Build the Docker image:

   ```bash
   ./build.sh
   ```

   Or build manually:
   ```bash
   docker build -t k8s-monitor .
   ```

2. Run with Docker (choose one option):

   **Basic run with kubeconfig mounting:**
   ```bash
   sudo docker run --rm -p 8080:8080 \
     -v ~/.kube:/home/k8s-monitor/.kube:ro \
     -v $(pwd)/data:/app/data \
     -e KUBECONFIG=/home/k8s-monitor/.kube/config \
     -e PERSISTENCE_FILE_PATH=/app/data/changes.json \
     --user $(id -u):$(id -g) \
     k8s-monitor:latest
   ```

   **With host networking (better for local Kubernetes access):**
   ```bash
   sudo docker run --rm --network host \
     -v ~/.kube:/home/k8s-monitor/.kube:ro \
     -v $(pwd)/data:/app/data \
     -e KUBECONFIG=/home/k8s-monitor/.kube/config \
     -e PERSISTENCE_FILE_PATH=/app/data/changes.json \
     --user $(id -u):$(id -g) \
     k8s-monitor:latest
   ```

   **For secure kubeconfig handling:**
   ```bash
   # Copy kubeconfig to temporary location with proper permissions
   mkdir -p /tmp/k8s-monitor-kube
   cp ~/.kube/config /tmp/k8s-monitor-kube/
   chmod 644 /tmp/k8s-monitor-kube/config

   # Run container
   sudo docker run --rm -p 8080:8080 \
     -v /tmp/k8s-monitor-kube:/home/k8s-monitor/.kube:ro \
     -v $(pwd)/data:/app/data \
     -e KUBECONFIG=/home/k8s-monitor/.kube/config \
     -e PERSISTENCE_FILE_PATH=/app/data/changes.json \
     --user $(id -u):$(id -g) \
     k8s-monitor:latest
   ```

### Option 3: OpenShift Deployment

For production deployment on OpenShift:

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

See `DOCKER.md` for detailed Docker and OpenShift deployment instructions.

## Usage

### Local Go Development

1. Build the application:

   ```
   go build -o k8s-monitor ./cmd/main.go
   ```

2. Run the application:

   ```bash
   # Basic run
   ./k8s-monitor
   
   # With custom persistence file path
   PERSISTENCE_FILE_PATH=/path/to/custom/changes.json ./k8s-monitor
   
   # With custom web port
   WEB_PORT=9090 ./k8s-monitor
   
   # With both custom settings
   PERSISTENCE_FILE_PATH=/tmp/k8s-changes.json WEB_PORT=9090 ./k8s-monitor
   ```

   The application will automatically create a default `config.json` file if one doesn't exist.

### Docker Usage

1. **Run with Docker (after building the image):**

   ```bash
   sudo docker run --rm -p 8080:8080 \
     -v ~/.kube:/home/k8s-monitor/.kube:ro \
     -v $(pwd)/data:/app/data \
     -e KUBECONFIG=/home/k8s-monitor/.kube/config \
     -e PERSISTENCE_FILE_PATH=/app/data/changes.json \
     --user $(id -u):$(id -g) \
     k8s-monitor:latest
   ```

   **Alternative with custom settings:**
   ```bash
   sudo docker run --rm -p 9090:9090 \
     -v ~/.kube:/home/k8s-monitor/.kube:ro \
     -v $(pwd)/data:/app/data \
     -e KUBECONFIG=/home/k8s-monitor/.kube/config \
     -e PERSISTENCE_FILE_PATH=/app/data/my-k8s-changes.json \
     -e WEB_PORT=9090 \
     --user $(id -u):$(id -g) \
     k8s-monitor:latest
   ```

2. **Access the web interface:**

   Open your web browser and navigate to `http://localhost:8080` to view the monitoring dashboard.

3. **Using the interface:**
   - The application will start monitoring the Kubernetes API objects and display changes in real-time
   - Use the "⚙️ Configuration" button in the web UI to see which resources are being monitored
   - Changes are automatically saved to the persistent data volume

### Production Deployment

For production environments, use the OpenShift template provided in `openshift-template.yaml` which includes:
- Proper RBAC configuration
- Resource limits and requests
- Health checks
- Security contexts
- Persistent storage options

## Customizing Configuration

Edit the `config.json` file to customize which resources to monitor:

```json
{
  "webPort": 8080,
  "persistence": {
    "enabled": true,
    "filePath": "changes.json",
    "autoSave": true,
    "saveInterval": 30
  },
  "logging": {
    "enabled": true,
    "logChanges": true,
    "logOperations": true
  },
  "resources": [
    {
      "name": "pods",
      "enabled": true,
      "namespace": "default",
      "description": "Kubernetes Pods"
    }
  ]
}
```

After editing the configuration, restart the application for changes to take effect.

## Troubleshooting

### Docker Issues

**Permission denied when accessing kubeconfig:**
```bash
# Solution: Run with proper user mapping
sudo docker run --rm -p 8080:8080 \
  -v ~/.kube:/home/k8s-monitor/.kube:ro \
  -v $(pwd)/data:/app/data \
  -e KUBECONFIG=/home/k8s-monitor/.kube/config \
  --user $(id -u):$(id -g) \
  k8s-monitor:latest
```

**Connection refused to Kubernetes API:**
- Use `--network host` for local Kubernetes clusters
- Ensure your kubeconfig is valid: `kubectl cluster-info`
- Check if your cluster is accessible from Docker containers

**Changes not persisting:**
- Ensure the data directory exists: `mkdir -p data`
- Check volume mount permissions
- Verify the container can write to `/app/data`

### General Issues

**Web interface not loading:**
- Check if the application started successfully
- Verify port 8080 is not in use by another application
- Check firewall settings

**No changes showing:**
- Verify RBAC permissions for the service account
- Check if resources exist in the monitored namespaces
- Review application logs for connection errors

## Files

- `cmd/main.go` - Main application entry point
- `pkg/config/` - Configuration management
- `pkg/monitor/` - Kubernetes monitoring logic
- `pkg/web/` - Web interface and API
- `pkg/utils/` - Utility functions
- `config.json` - Application configuration
- `Dockerfile` - Container image definition
- `openshift-template.yaml` - OpenShift deployment template
- `build.sh` - Build script for Docker images
- `DOCKER.md` - Detailed Docker and OpenShift deployment guide

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.