# Multi-stage build for Go application
# Stage 1: Build stage
FROM registry.access.redhat.com/ubi9/go-toolset:1.21 AS builder

# Set working directory
WORKDIR /workspace

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with static linking for better portability
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-monitor ./cmd/main.go

# Stage 2: Runtime stage
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Install necessary packages for running the application
RUN microdnf update -y && \
    microdnf install -y ca-certificates && \
    microdnf clean all

# Create a non-root user (OpenShift best practice)
RUN useradd -r -u 1001 -g 0 -s /sbin/nologin \
    -c "k8s-monitor user" k8s-monitor

# Create directories with proper permissions
RUN mkdir -p /app/data /home/k8s-monitor && \
    chown -R 1001:0 /app /home/k8s-monitor && \
    chmod -R g=u /app /home/k8s-monitor

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /workspace/k8s-monitor .

# Copy configuration files
COPY --chown=1001:0 config.json .

# Ensure the binary is executable
RUN chmod +x k8s-monitor

# Switch to non-root user
USER 1001

# Expose port (configurable via environment)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables with defaults
ENV WEB_PORT=8080
ENV PERSISTENCE_FILE_PATH=/app/data/changes.json
ENV HOME=/home/k8s-monitor

# Define volume for persistent data
VOLUME ["/app/data"]

# Labels for better metadata (OpenShift best practice)
LABEL name="k8s-monitor" \
      vendor="k8s-monitor" \
      version="1.0" \
      summary="Kubernetes monitoring application" \
      description="A Go application that monitors changes in Kubernetes API objects and displays them in a web interface. Supports environment variable configuration for WEB_PORT and PERSISTENCE_FILE_PATH." \
      io.k8s.display-name="K8s Monitor" \
      io.k8s.description="Kubernetes monitoring application with web interface" \
      io.openshift.expose-services="8080:http" \
      io.openshift.tags="monitoring,kubernetes,go"

# Command to run the application
CMD ["./k8s-monitor"]
