# Contributing to claude-mux

Thank you for your interest in contributing to claude-mux! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be kind, respectful, and considerate to others. We're all here to make claude-mux better.

## How to Contribute

### Reporting Issues

- Check if the issue already exists
- Include steps to reproduce the issue
- Include your OS, Go version, and claude-mux version
- Include any relevant error messages or logs

### Suggesting Features

- Open an issue with the "enhancement" label
- Describe the use case and why it would be valuable
- Be open to discussion and alternative approaches

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Set up your development environment**:
   ```bash
   git clone https://github.com/yourusername/claude-mux.git
   cd claude-mux
   make dev-deps  # Install development dependencies
   ```

3. **Make your changes**:
   - Write clear, concise commit messages
   - Add tests for new functionality
   - Update documentation as needed
   - Follow the existing code style

4. **Ensure quality**:
   ```bash
   make check  # Runs fmt, vet, lint, security, and tests
   ```

5. **Submit your PR**:
   - Provide a clear description of the changes
   - Reference any related issues
   - Be responsive to feedback

## Development Workflow

### Project Structure

```
claude-mux/
├── cmd/claude-mux/    # CLI entry point
├── internal/          # Private packages
│   ├── git/          # Git operations
│   ├── worktree/     # Worktree management
│   └── config/       # Configuration
└── pkg/              # Public packages (future)
```

### Common Commands

```bash
# Install dev dependencies
make dev-deps

# Build the binary
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Run security scan
make security

# Run all checks
make check

# Test release process
make release-dry
```

### Testing

- Write unit tests for new functions
- Use table-driven tests where appropriate
- Mock external dependencies (git commands)
- Aim for >80% code coverage

Example test:
```go
func TestFeature(t *testing.Testimport os.Args) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "expected", false},
        {"invalid input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Feature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Feature() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("Feature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Code Style

- Use `goimports` for formatting (via `make fmt`)
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Keep functions small and focused
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors explicitly

### Commit Messages

Follow conventional commits format:
```
type: subject

body (optional)

footer (optional)
```

Types:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test changes
- `refactor:` Code refactoring
- `chore:` Build/tooling changes

Example:
```
feat: add session persistence

Implement basic session management to allow switching between
active Claude sessions without losing state.

Closes #42
```

## Release Process

Releases are automated via GitHub Actions when a tag is pushed:

1. Ensure all tests pass on `main`
2. Update version if needed
3. Create and push a tag:
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```
4. GitHub Actions will create the release with binaries

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions in existing issues
- Check the README for usage information

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
