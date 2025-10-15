# Kubernetes Monitor 🚀

A lightweight, web-based tool for real-time monitoring and tracking of Kubernetes API object changes. Perfect for DevOps engineers, cluster administrators, and anyone who needs to quickly understand what's happening in their Kubernetes clusters.

## Features

- 🔍 **Real-time Monitoring**: Live tracking of all major Kubernetes API objects with instant change detection
- 📊 **Change Analytics**: Comprehensive statistics showing total changes, unread notifications, and session activity
- 🎨 **Modern UI**: Clean, responsive interface with intuitive filtering and advanced search capabilities
- ⚡ **Fast & Efficient**: Memory-optimized tracking with configurable persistence and auto-save functionality
- 🔎 **Advanced Filtering**: Filter by resource type, event type (Added/Modified/Deleted), read status, and timestamps
- 📤 **Data Persistence**: Configurable JSON file storage with automatic backup and recovery
- 🐳 **Container Ready**: Multi-architecture Docker image (AMD64, ARM64) with health checks
- 🔒 **Secure**: Runs as non-root user, minimal RBAC permissions, security contexts
- 🏥 **Production Ready**: Kubernetes manifests with proper health checks and resource limits
- ⚙️ **Configurable**: Extensive configuration options for resources, namespaces, and monitoring behavior
- 🌐 **RESTful API**: Full API access for automation and integration with other tools
- 📋 **Mark as Read**: Efficient change management with read/unread status tracking

## Quick Start

### Docker

```bash
# Run with your kubeconfig
docker run -p 8080:8080 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  k8s-monitor:latest

# Visit in browser
open http://localhost:8080
```

### Docker Compose

```bash
# Start with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f
```

### Kubernetes

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Port forward for testing
kubectl port-forward service/k8s-monitor 8080:80

# Visit
open http://localhost:8080
```

## Usage Examples

### Browser Access
Visit `http://localhost:8080` in your web browser to see:
- 📋 **Real-time change feed** - All Kubernetes object changes with timestamps and details
- 📊 **Statistics dashboard** - Total changes, unread notifications, and session metrics
- 🔍 **Advanced filters** - Filter by resource type, event type, namespace, and read status
- ✅ **Mark as read** - Manage notifications with read/unread status
- ⚙️ **Configuration panel** - View monitored resources and persistence settings

### API Endpoints

| Endpoint | Description | Response Format |
|----------|-------------|-----------------|
| `/api/changes` | List all monitored changes | JSON |
| `/api/stats` | Get monitoring statistics | JSON |
| `/api/config` | Get current configuration | JSON |
| `/api/mark-read` | Mark change as read | JSON |
| `/api/mark-all-read` | Mark all changes as read | JSON |
| `/api/save-now` | Trigger immediate save | JSON |
| `/api/debug` | Debug status and version | JSON |
| `/health` | Health check endpoint | JSON |

### API Examples

```bash
# Get all changes
curl http://localhost:8080/api/changes

# Get statistics
curl http://localhost:8080/api/stats

# Mark all changes as read
curl -X POST http://localhost:8080/api/mark-all-read

# Trigger save to file
curl -X POST http://localhost:8080/api/save-now

# Health check
curl http://localhost:8080/health
```

## Configuration

The application uses a `config.json` file to configure which resources to monitor and other settings:

# Kubernetes Monitor 🚀

A lightweight, web-based tool for real-time monitoring of Kubernetes API object changes. Perfect for DevOps engineers, cluster administrators, and anyone who needs to track resource modifications in their Kubernetes clusters.

## Features

- 🔍 **Real-time Monitoring**: Live tracking of all major Kubernetes API objects with instant change detection
- 📊 **Change Analytics**: Comprehensive statistics showing event types, resource counts, and activity patterns
- 🎨 **Modern UI**: Clean, responsive interface with intuitive navigation and dark/light theme support
- 🚀 **Fast & Efficient**: Memory-efficient monitoring with configurable persistence and auto-save
- 🔎 **Advanced Filtering**: Search and filter changes by resource type, event type, namespace, and read status
- 📤 **Persistent Storage**: Save changes to JSON files with configurable auto-save intervals
- 🐳 **Container Ready**: Multi-architecture Docker image with Alpine base and non-root user
- 🔒 **Secure**: Runs as non-root user with minimal RBAC permissions
- 🏥 **Health Checks**: Built-in health check endpoints for monitoring
- ⚡ **Production Ready**: Go-based server with efficient resource handling and graceful shutdown

## Quick Start

### Docker

```bash
# Run with your kubeconfig
docker run -p 8080:8080 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  k8s-monitor:latest

# Visit in browser
open http://localhost:8080
```

### Docker Compose

```bash
# Start with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f
```

### Kubernetes

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Port forward for testing
kubectl port-forward service/k8s-monitor 8080:80

# Visit
open http://localhost:8080
```

## Usage Examples

### Browser Access
Visit `http://localhost:8080` in your web browser to see:
- 📋 **Real-time Change Feed** - Live updates of all Kubernetes object modifications
- 📊 **Statistics Dashboard** - Total changes, unread count, session info, and uptime
- 🎛️ **Advanced Filters** - Filter by event type, resource type, namespace, and read status
- ⚙️ **Configuration Panel** - View monitored resources and persistence settings
- 💾 **Data Management** - Mark changes as read and save to persistent storage

### API Endpoints

| Endpoint | Description | Response Format |
|----------|-------------|-----------------|
| `/api/changes` | Get all monitored changes | JSON |
| `/api/stats` | Get monitoring statistics | JSON |
| `/api/config` | Get current configuration | JSON |
| `/api/mark-read` | Mark specific change as read | JSON |
| `/api/mark-all-read` | Mark all changes as read | JSON |
| `/api/save-now` | Force save to persistent storage | JSON |
| `/api/debug` | Debug status and version info | JSON |
| `/health` | Health check endpoint | JSON |

### API Examples

```bash
# Get all changes
curl http://localhost:8080/api/changes

# Get monitoring statistics
curl http://localhost:8080/api/stats

# Mark all changes as read
curl -X POST http://localhost:8080/api/mark-all-read

# Force save to file
curl -X POST http://localhost:8080/api/save-now

# Health check
curl http://localhost:8080/health

# Debug information
curl http://localhost:8080/api/debug
```

## Configuration

The application can be configured using environment variables and a `config.json` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Port to listen on |
| `DEBUG` | `false` | Enable debug mode with verbose logging |
| `KUBECONFIG` | `~/.kube/config` | Path to kubeconfig file |
| `PERSISTENCE_FILE_PATH` | `/app/data/changes.json` | Path to persistent storage file |

### Docker Environment Variables

```bash
# Enable debug mode
docker run -p 8080:8080 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  -e DEBUG=true \
  k8s-monitor:latest

# Custom port and persistence path
docker run -p 3000:3000 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  -v ./data:/app/data \
  -e PORT=3000 \
  -e PERSISTENCE_FILE_PATH=/app/data/custom-changes.json \
  k8s-monitor:latest
```

### Configuration File (config.json)

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
    "enabled": false,
    "logChanges": false,
    "logOperations": false
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

## Deployment Examples

### Docker

#### Simple Run
```bash
docker run -d --name k8s-monitor -p 8080:8080 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  k8s-monitor:latest
```

#### With Debug Mode and Persistence
```bash
docker run -d --name k8s-monitor -p 8080:8080 \
  -v ~/.kube/config:/home/app/.kube/config:ro \
  -v ./data:/app/data \
  -e DEBUG=true \
  -e PERSISTENCE_FILE_PATH=/app/data/changes.json \
  k8s-monitor:latest
```

### Docker Compose

See [`docker-compose.yml`](docker-compose.yml) for a complete example with:
- Environment variable configuration
- Volume mounts for kubeconfig and data
- Network configuration
- Health checks
- Optional nginx reverse proxy

```bash
docker-compose up -d
```

### Kubernetes

The `k8s/` directory contains complete Kubernetes manifests:

- **RBAC**: ServiceAccount, ClusterRole, and ClusterRoleBinding with minimal read-only permissions
- **Deployment**: Single replica deployment with health checks and security contexts
- **Service**: ClusterIP service for internal access
- **Ingress**: Optional ingress for external access

#### Basic Deployment
```bash
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

#### With Ingress
```bash
# Update host in k8s/ingress.yaml first
kubectl apply -f k8s/ingress.yaml
```

#### Port Forward for Testing
```bash
kubectl port-forward service/k8s-monitor 8080:80
```

See [k8s/README.md](k8s/README.md) for detailed Kubernetes deployment instructions.

## Development

### Local Development

```bash
# Clone the repository
git clone https://github.com/MichaelTrip/k8s-monitor.git
cd k8s-monitor

# Install dependencies
go mod download

# Run the application
go run cmd/main.go

# Or with debug mode
DEBUG=true go run cmd/main.go
```

### Building the Container

```bash
# Build locally
./build.sh

# Build with custom version
APP_VERSION=v1.1.0 ./build.sh

# Build and push to registry
REGISTRY=your-registry.com ./build.sh
```

### Testing

```bash
# Build and run
go build -o bin/k8s-monitor cmd/main.go
./bin/k8s-monitor

# Test API endpoints
curl http://localhost:8080/api/stats
curl http://localhost:8080/health

# Test with Docker
docker build -t k8s-monitor:test .
docker run -p 8080:8080 k8s-monitor:test
```

## Architecture

The application follows a simple, efficient architecture:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Browser   │    │   Load Balancer  │    │  K8s Monitor    │
│                 │────▶│     /Ingress     │────▶│                 │
│                 │    │                  │    │   Go Server     │
└─────────────────┘    └──────────────────┘    │   + K8s Client  │
                                               │                 │
┌─────────────────┐                            │   Port 8080     │
│   curl/API      │────────────────────────────▶│                 │
│   Client        │                            └─────────────────┘
└─────────────────┘                                       │
                                                          ▼
                                               ┌─────────────────┐
                                               │  Kubernetes API │
                                               │     Server      │
                                               └─────────────────┘
```

### Key Components

- **Frontend**: Modern HTML5/CSS3/JavaScript with no frameworks
- **Backend**: Go with gorilla/mux router and client-go library for Kubernetes interaction
- **Storage**: Configurable persistent storage to JSON files with auto-save
- **API**: RESTful JSON API with comprehensive endpoints

## Use Cases

- **Change Monitoring**: Real-time tracking of all Kubernetes resource modifications
- **Debugging**: Understand what changed when troubleshooting issues
- **Auditing**: Keep track of resource changes for compliance and documentation
- **Learning**: Explore Kubernetes API behavior and resource interactions
- **Development**: Monitor application deployments and configuration changes
- **Operations**: Track system changes during maintenance windows

## RBAC Requirements

For proper operation in Kubernetes, the service account needs:

```yaml
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list", "watch"]
```

This provides **read-only** access to monitor resources, which is required for change detection and monitoring.

## Performance & Monitoring

The application is designed for efficiency:

- **Memory Efficient**: Only tracks changes since program start (configurable retention)
- **Lightweight**: Alpine-based container image under 50MB
- **Fast Response**: In-memory change storage with optional persistence
- **Resource Aware**: Configurable resource monitoring to reduce API load

This allows the application to handle:
- ✅ Clusters with hundreds of resources
- ✅ Multiple concurrent users
- ✅ Long-running monitoring sessions
- ✅ High-frequency change environments

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `chore:` for maintenance tasks

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/MichaelTrip/k8s-monitor/issues)
- 📖 **Documentation**: This README and [k8s/README.md](k8s/README.md)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/MichaelTrip/k8s-monitor/discussions)

## Related Projects

- [k8s-object-explorer](https://github.com/MichaelTrip/k8s-object-explorer) - Kubernetes resource explorer and browser
- [myipcontainer](https://github.com/MichaelTrip/myipcontainer) - Simple IP address display container
- [lmsensors-container](https://github.com/MichaelTrip/lmsensors-container) - Hardware sensor monitoring for Kubernetes

---

Made with ❤️ by [MichaelTrip](https://github.com/MichaelTrip)

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