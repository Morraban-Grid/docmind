# Contributing to DocMind

Thank you for your interest in contributing to DocMind! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect differing viewpoints and experiences

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/yourusername/docmind.git
   cd docmind
   ```
3. **Set up development environment**
   ```bash
   ./scripts/bootstrap.sh
   ```
4. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### 1. Make Your Changes

- Follow the existing code style
- Write clear, descriptive commit messages
- Add tests for new features
- Update documentation as needed

### 2. Code Style

**Go:**
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add comments for exported functions

**Python:**
- Follow [PEP 8](https://www.python.org/dev/peps/pep-0008/)
- Use `black` for formatting
- Use `pylint` for linting
- Add docstrings for all functions

### 3. Testing

**Run tests before committing:**
```bash
# Go tests
make test-go

# Python tests
make test-python

# All tests
make test
```

### 4. Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(auth): add JWT token refresh endpoint

Implement token refresh functionality to allow users to
extend their session without re-authenticating.

Closes #123
```

```
fix(upload): validate file size before processing

Add file size validation to prevent memory issues with
large files. Maximum size is now 50MB.

Fixes #456
```

### 5. Pull Request Process

1. **Update your branch with latest main**
   ```bash
   git fetch origin
   git rebase origin/main
   ```

2. **Push your changes**
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create Pull Request**
   - Use a clear, descriptive title
   - Reference related issues
   - Describe what changed and why
   - Include screenshots for UI changes
   - Ensure all tests pass
   - Request review from maintainers

4. **Address Review Feedback**
   - Make requested changes
   - Push updates to your branch
   - Respond to comments

5. **Merge**
   - Maintainers will merge once approved
   - Delete your feature branch after merge

## Security

### Reporting Security Issues

**DO NOT** open public issues for security vulnerabilities.

Instead, email security concerns to: [security@docmind.example.com]

### Security Guidelines

- Never commit secrets, API keys, or credentials
- Always use environment variables for sensitive data
- Validate all user inputs
- Use prepared statements for database queries
- Keep dependencies up to date

## Project Structure

```
docmind/
├── .github/              # GitHub workflows and templates
├── .kiro/specs/          # Specification documents
├── deployments/          # Deployment configurations
├── docs/                 # Documentation
├── proto/                # gRPC protocol definitions
├── scripts/              # Utility scripts
├── services/
│   ├── go-user-service/  # Go microservice
│   └── python-rag-service/ # Python microservice
└── README.md
```

## Adding New Features

### New File Format Support

1. Add parser in `services/python-rag-service/app/utils/file_parsers.py`
2. Update `TextExtractor` class
3. Add tests in `tests/test_file_parsers.py`
4. Update documentation

### New API Endpoint

1. Define route in appropriate handler
2. Implement business logic in service layer
3. Add validation
4. Add tests
5. Update API documentation

### New gRPC Method

1. Update `proto/rag/rag.proto`
2. Regenerate code: `./scripts/generate-proto.sh`
3. Implement in Python service
4. Update Go client
5. Add tests

## Documentation

- Update README.md for user-facing changes
- Update docs/architecture.md for architectural changes
- Add inline comments for complex logic
- Update CHANGELOG.md

## Questions?

- Open a discussion on GitHub
- Check existing issues and PRs
- Read the documentation in `/docs`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to DocMind! 🚀
