# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the k8s-monitor application.

## Files

- `rbac.yaml` - ServiceAccount, ClusterRole, and ClusterRoleBinding (read-only access)
- `deployment.yaml` - Application deployment with health checks and security contexts
- `service.yaml` - ClusterIP service to expose the application internally
- `ingress.yaml` - Ingress resource for external access

## Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured to access your cluster
- Ingress controller installed (e.g., nginx-ingress) for external access

## Quick Deploy

### Option 1: All Resources (Recommended)
Deploy all resources including RBAC at once:

```bash
kubectl apply -f k8s/
```

### Option 2: Selective Deployment
Deploy individually for more control:

```bash
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml  # Optional
```

## Configuration

### Ingress

Before deploying the ingress, update the host in `ingress.yaml`:

```yaml
rules:
- host: k8s-monitor.example.com  # Change this to your actual domain
```

### TLS (Optional)

To enable HTTPS, uncomment the TLS section in `ingress.yaml` and create a TLS secret:

```bash
kubectl create secret tls k8s-monitor-tls --cert=path/to/tls.crt --key=path/to/tls.key
```

Or use cert-manager:

```bash
# Uncomment the cert-manager annotation in ingress.yaml
# cert-manager.io/cluster-issuer: "letsencrypt-prod"
```

### Custom Namespace

Deploy to a custom namespace:

```bash
kubectl create namespace k8s-monitor
kubectl apply -f k8s/ -n k8s-monitor
```

**Note**: Update the `ClusterRoleBinding` namespace in `rbac.yaml` if using a custom namespace.

## RBAC Permissions

The application requires **read-only** access to monitor resources across the cluster:

### ServiceAccount
- Name: `k8s-monitor`
- Namespace: `default` (or custom namespace)

### ClusterRole Permissions
```yaml
rules:
  # Read access to all API resources for monitoring
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["get", "list", "watch"]
  
  # Specific access to namespaces
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get", "list", "watch"]
```

### Why These Permissions?
- **`get`, `list`, `watch` on all resources**: Required to monitor changes across all resource types
- **No write access**: Application is read-only and cannot modify any cluster resources
- **Cluster-wide access**: Needed to monitor resources across all namespaces
- **Watch permissions**: Required for real-time change detection

### Security Note
While the permissions are cluster-wide, they are **read-only**. The application:
- ✅ Cannot create, update, or delete any resources
- ✅ Cannot access secrets' values (only lists that they exist)
- ✅ Runs as non-root user (UID 1000)
- ✅ Has dropped all Linux capabilities

## Accessing the Application

1. **Port Forward (for testing)**:
   ```bash
   kubectl port-forward service/k8s-monitor 8080:80
   ```
   Then visit `http://localhost:8080`

2. **Via Ingress** (production):
   - Ensure your ingress controller is installed
   - Update your DNS to point to the ingress controller's external IP
   - Visit `http://your-domain.com`

## Monitoring

Check the deployment status:

```bash
kubectl get deployments
kubectl get pods -l app=k8s-monitor
kubectl get services
kubectl get ingress
```

View logs:

```bash
kubectl logs -l app=k8s-monitor
kubectl logs -l app=k8s-monitor --tail=100 -f
```

## Scaling

Scale the deployment (though typically one instance is sufficient):

```bash
kubectl scale deployment k8s-monitor --replicas=2
```

## Resource Requirements

The application is designed to be lightweight:

- **CPU**: 100m requests, 500m limits
- **Memory**: 128Mi requests, 512Mi limits
- **Storage**: Uses emptyDir for temporary data (consider PVC for persistence)

For production, consider:
- Adding resource quotas
- Setting up monitoring and alerting
- Configuring log aggregation

## Troubleshooting

### Pods not starting
```bash
kubectl describe pod -l app=k8s-monitor
kubectl logs -l app=k8s-monitor
```

### Permission errors
Ensure the `ServiceAccount`, `ClusterRole`, and `ClusterRoleBinding` are properly created:
```bash
kubectl get serviceaccount k8s-monitor
kubectl get clusterrole k8s-monitor
kubectl get clusterrolebinding k8s-monitor
```

### Ingress not working
Check ingress controller logs and ensure your DNS is properly configured:
```bash
kubectl get ingress k8s-monitor
kubectl describe ingress k8s-monitor
```

### No changes showing in the UI
- Verify RBAC permissions are correctly applied
- Check if there are resources in the cluster to monitor
- Review application logs for connection errors

## Health Checks

The deployment includes health checks:

- **Liveness Probe**: `/health` endpoint on port 8080
- **Readiness Probe**: `/health` endpoint on port 8080

Health check configuration:
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

## Security Context

The application runs with a secure configuration:

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: false
  capabilities:
    drop:
      - ALL
```

## Cleanup

Remove all resources:

```bash
kubectl delete -f k8s/
```

Or individually:

```bash
kubectl delete deployment k8s-monitor
kubectl delete service k8s-monitor
kubectl delete ingress k8s-monitor
kubectl delete serviceaccount k8s-monitor
kubectl delete clusterrole k8s-monitor
kubectl delete clusterrolebinding k8s-monitor
```

## Production Considerations

For production deployments:

1. **Persistence**: Consider using a PersistentVolumeClaim for data storage
2. **Monitoring**: Set up Prometheus metrics and Grafana dashboards
3. **Logging**: Configure centralized logging (ELK, Fluentd, etc.)
4. **Security**: Regular security scans and updates
5. **Backup**: Regular backups of configuration and data
6. **Resource Limits**: Adjust based on cluster size and monitoring needs
7. **Network Policies**: Implement network policies for additional security

## Examples

### Custom Configuration with PersistentVolume

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: k8s-monitor-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

Then update the deployment to use the PVC instead of emptyDir.

### Multi-Namespace Monitoring

For monitoring specific namespaces, create Role and RoleBinding instead of ClusterRole:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: monitoring
  name: k8s-monitor
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
```