# GitHub Actions CI/CD Pipeline

This repository uses GitHub Actions for continuous integration and deployment with semantic versioning.

## Pipeline Overview

The CI/CD pipeline consists of several jobs that run automatically on push and pull request events:

### 1. Test Job
- Runs Go tests and static analysis
- Caches Go modules for faster builds
- Validates code quality with `go vet`

### 2. Build and Push Job
- Builds multi-platform Docker images (linux/amd64, linux/arm64)
- Pushes images to GitHub Container Registry (ghcr.io)
- Tags images with branch name, SHA, and `latest` for main branch

### 3. Semantic Release Job (main branch only)
- Analyzes commit messages using conventional commits
- Generates semantic version numbers automatically
- Creates GitHub releases with changelog
- Only runs on pushes to main branch

### 4. Release Docker Image Job
- Tags and pushes Docker images with semantic version
- Creates both `v1.2.3` and `1.2.3` format tags
- Only runs when a new release is created

### 5. Security Scan Job
- Scans Docker images for vulnerabilities using Trivy
- Uploads results to GitHub Security tab

## Semantic Versioning

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automatic semantic versioning:

- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `BREAKING CHANGE:` - Breaking changes (major version bump)
- `chore:`, `docs:`, `style:`, `refactor:`, `test:` - No version bump

### Commit Examples

```bash
# Patch release (1.0.0 → 1.0.1)
git commit -m "fix: resolve memory leak in monitor"

# Minor release (1.0.0 → 1.1.0)
git commit -m "feat: add new monitoring endpoint"

# Major release (1.0.0 → 2.0.0)
git commit -m "feat: redesign API structure

BREAKING CHANGE: API endpoints have been restructured"
```

## Docker Images

Images are automatically built and pushed to:
- `ghcr.io/michaeltrip/k8s-monitor:latest` (main branch)
- `ghcr.io/michaeltrip/k8s-monitor:main` (main branch)
- `ghcr.io/michaeltrip/k8s-monitor:v1.2.3` (semantic version)
- `ghcr.io/michaeltrip/k8s-monitor:1.2.3` (semantic version without 'v')

## Using the Images

```bash
# Pull latest version
docker pull ghcr.io/michaeltrip/k8s-monitor:latest

# Pull specific version
docker pull ghcr.io/michaeltrip/k8s-monitor:v1.0.0

# Run container
docker run -p 8080:8080 ghcr.io/michaeltrip/k8s-monitor:latest
```

## Required Secrets

The pipeline uses the built-in `GITHUB_TOKEN` which is automatically provided by GitHub Actions. No additional secrets are required for basic functionality.

## Manual Release

If you need to trigger a release manually (not recommended), you can push an empty commit:

```bash
git commit --allow-empty -m "chore: trigger release"
git push origin main
```

## Pipeline Triggers

- **Pull Requests**: Runs tests and builds (no release)
- **Push to main**: Full pipeline including semantic release
- **Manual workflow dispatch**: Can be enabled if needed

## Monitoring the Pipeline

1. Check the Actions tab in your GitHub repository
2. View build logs and status
3. Monitor security scan results in the Security tab
4. Check releases in the Releases section

## Troubleshooting

- **Tests failing**: Check the test job logs for specific errors
- **Build failing**: Verify Dockerfile and dependencies
- **No release created**: Ensure commits follow conventional format
- **Permission errors**: Check that Actions have necessary permissions
