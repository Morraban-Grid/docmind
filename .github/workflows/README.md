# GitHub Actions Workflows

This directory contains CI/CD workflows for the DocMind project.

## Current Workflows

### security-scan.yml
Runs security scans on every push and pull request:
- **Secret Scanning**: Detects exposed secrets using Gitleaks
- **Dependency Scanning**: Scans for vulnerabilities using Trivy

## Future Workflows (Coming in Later Iterations)

### ci-go.yml (Iteration 2+)
Will run Go tests and linting when Go service is implemented.

### ci-python.yml (Iteration 4+)
Will run Python tests and linting when Python service is implemented.

## Permissions

The security-scan workflow requires:
- `contents: read` - To checkout code
- `security-events: write` - To upload security results
- `actions: read` - To read workflow status

## Notes

- Workflows are triggered on push to `main` and `develop` branches
- Pull requests also trigger workflows
- Security scan results appear in the Security tab
