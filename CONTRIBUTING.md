# Contributing to Toutā DataMapper MySQL Adapter

Thank you for your interest in contributing to the MySQL adapter! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker
- Check if the issue already exists
- Provide detailed information:
  - Go version
  - MySQL version
  - Operating system
  - Steps to reproduce
  - Expected vs actual behavior

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write or update tests
5. Ensure tests pass (`go test ./...`)
6. Ensure code is formatted (`go fmt ./...`)
7. Run linter (`golangci-lint run`)
8. Commit with conventional commit format
9. Push to your fork
10. Open a Pull Request

### Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `perf`: Performance improvement
- `refactor`: Code restructuring
- `test`: Test additions or modifications
- `docs`: Documentation changes
- `chore`: Build, CI, or tooling changes

**Examples:**
```
feat(adapter): add bulk update support
fix(connection): handle connection timeout correctly
perf(query): optimize parameter binding
docs(readme): add connection pooling examples
test(adapter): add transaction tests
```

## Development Setup

### Prerequisites

- Go 1.22 or higher
- MySQL 5.7+ or MariaDB 10.2+
- Git
- golangci-lint (for linting)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/toutago-datamapper-mysql
cd toutago-datamapper-mysql

# Install dependencies
go mod download

# Set up test database
export MYSQL_TEST_DSN="user:password@tcp(localhost:3306)/testdb"

# Run tests
go test ./...
```

## Testing Requirements

- All new code must include tests
- Maintain minimum 80% code coverage
- Tests must pass with race detector: `go test -race ./...`
- Test against MySQL 5.7+ and MariaDB 10.2+

## Code Quality Standards

- Follow Go best practices and idioms
- Use meaningful variable and function names
- Keep functions focused and small
- Document exported types and functions
- Pass golangci-lint without errors
- Follow MySQL best practices for queries

## Documentation

- Update README.md for user-facing changes
- Update doc.go for API changes
- Add examples for new features
- Keep CHANGELOG.md current

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
