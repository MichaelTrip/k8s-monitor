# Conventional Commits Guide

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated semantic versioning.

## Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Commit Types

| Type | Description | Version Impact |
|------|-------------|----------------|
| `feat` | New feature | Minor (1.0.0 → 1.1.0) |
| `fix` | Bug fix | Patch (1.0.0 → 1.0.1) |
| `BREAKING CHANGE` | Breaking change | Major (1.0.0 → 2.0.0) |
| `docs` | Documentation only | No release |
| `style` | Code style changes | No release |
| `refactor` | Code refactoring | No release |
| `test` | Adding/updating tests | No release |
| `chore` | Maintenance tasks | No release |
| `ci` | CI/CD changes | No release |
| `perf` | Performance improvements | Patch |
| `revert` | Revert previous commit | Patch |

## Examples

### Patch Release (Bug Fix)
```bash
git commit -m "fix: resolve kubernetes connection timeout"
git commit -m "fix(monitor): handle nil pointer in event processing"
```

### Minor Release (New Feature)
```bash
git commit -m "feat: add health check endpoint"
git commit -m "feat(web): implement real-time event streaming"
```

### Major Release (Breaking Change)
```bash
git commit -m "feat: redesign configuration format

BREAKING CHANGE: configuration file format has changed from JSON to YAML"
```

### No Release
```bash
git commit -m "docs: update README with installation instructions"
git commit -m "test: add unit tests for monitor package"
git commit -m "chore: update dependencies"
```

## Multi-paragraph Commits

```bash
git commit -m "feat: add persistent storage for events

This change introduces a new file-based persistence layer that
allows events to be stored across application restarts.

- Add FileWriter utility for atomic file operations
- Implement event serialization/deserialization
- Add configuration option for storage file path

Closes #123"
```

## Breaking Changes

You can indicate breaking changes in two ways:

1. **In the footer:**
```bash
git commit -m "feat: update API response format

BREAKING CHANGE: API responses now include timestamp in ISO format"
```

2. **With exclamation mark:**
```bash
git commit -m "feat!: update API response format"
```

## Best Practices

1. **Use lowercase** for the type and description
2. **Use present tense** ("add feature" not "added feature")
3. **Be descriptive** but concise in the description
4. **Include scope** when it adds clarity (`fix(docker):`, `feat(api):`)
5. **Reference issues** in the footer when applicable (`Closes #123`)
6. **Explain the why** in the body for complex changes

## Scopes (Optional)

Common scopes for this project:
- `api` - API related changes
- `web` - Web interface changes
- `monitor` - Monitoring logic changes
- `config` - Configuration changes
- `docker` - Docker/containerization changes
- `ci` - CI/CD pipeline changes
- `docs` - Documentation changes

Example: `feat(api): add new monitoring endpoint`
