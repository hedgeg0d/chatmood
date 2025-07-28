# Contributing to ChatMood

Thank you for your interest in contributing to ChatMood! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Community](#community)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful**: Treat everyone with respect and kindness
- **Be inclusive**: Welcome newcomers and help them succeed
- **Be collaborative**: Work together and share knowledge
- **Be patient**: Remember that everyone is learning
- **Be constructive**: Provide helpful feedback and suggestions

## Getting Started

### Prerequisites

Before contributing, make sure you have:

- Go 1.24.5 or higher
- Docker and Docker Compose
- Git
- A Telegram bot token (for testing)
- Basic knowledge of Go, JavaScript, and Telegram Web Apps

### First-time Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/chatmood.git
   cd chatmood
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/hedgeg0d/chatmood.git
   ```

3. **Set up environment**
   ```bash
   make setup
   # Edit .env with your bot token and settings
   ```

4. **Install development tools**
   ```bash
   make install-tools
   ```

5. **Run tests**
   ```bash
   make test
   ```

6. **Start development server**
   ```bash
   make dev
   ```

## Development Setup

### Local Development

```bash
# Start development server with hot reload
make dev

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Format code
make format

# Security check
make security
```

### Docker Development

```bash
# Start development environment with Docker
make docker-dev

# View logs
make docker-logs

# Stop containers
make docker-stop
```

## Contributing Process

### 1. Choose an Issue

- Browse [open issues](https://github.com/hedgeg0d/chatmood/issues)
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to indicate you're working on it

### 2. Create a Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

### 3. Make Changes

- Write clean, documented code
- Follow coding standards (see below)
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes

```bash
# Run all tests
make test

# Check code quality
make lint

# Run security scan
make security

# Test manually in browser
make dev
```

### 5. Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add sticker animation effects

- Add bounce, fade, and slide animations
- Update canvas rendering for smoother transitions
- Add animation controls in UI
- Update tests for new animation features

Closes #123"
```

### 6. Push and Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create pull request on GitHub
```

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `goimports` for formatting
- Write clear, self-documenting code
- Include comments for exported functions
- Handle errors appropriately

**Example:**
```go
// GenerateSticker creates a custom mood sticker with the provided configuration.
// It returns the sticker ID on success or an error if generation fails.
func GenerateSticker(config StickerConfig) (string, error) {
    if err := validateConfig(config); err != nil {
        return "", fmt.Errorf("invalid config: %w", err)
    }
    
    // Implementation...
    return stickerID, nil
}
```

### JavaScript Code Style

- Use ES6+ features
- Follow consistent naming conventions
- Add JSDoc comments for functions
- Use async/await for promises
- Handle errors gracefully

**Example:**
```javascript
/**
 * Generates a sticker with the provided configuration
 * @param {Object} config - Sticker configuration
 * @param {string} config.emoji - Emoji character
 * @param {string} config.text - Custom text
 * @returns {Promise<Object>} Generation result
 */
async function generateSticker(config) {
    try {
        const response = await fetch('/api/generate-sticker', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('Sticker generation failed:', error);
        throw error;
    }
}
```

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
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
feat(stickers): add animation effects support

fix(api): handle invalid base64 image data

docs(readme): update installation instructions

test(bot): add unit tests for message handling
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/bot/...

# Run with verbose output
go test -v ./...
```

### Writing Tests

- Write unit tests for all new functions
- Use table-driven tests where appropriate
- Include edge cases and error conditions
- Aim for high test coverage (>80%)

**Example:**
```go
func TestGenerateSticker(t *testing.T) {
    tests := []struct {
        name    string
        config  StickerConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: StickerConfig{
                Emoji: "😊",
                Text:  "Happy",
                Mood:  "happy",
            },
            wantErr: false,
        },
        {
            name: "invalid mood",
            config: StickerConfig{
                Emoji: "😊",
                Mood:  "invalid",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := GenerateSticker(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateSticker() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Document all exported functions and types
- Include usage examples
- Explain complex algorithms
- Update README.md for user-facing changes

### API Documentation

- Update API.md for API changes
- Include request/response examples
- Document error codes and messages

### Comments

```go
// Package bot provides Telegram bot functionality for ChatMood.
package bot

// Bot represents a Telegram bot instance with webhook support.
type Bot struct {
    Token      string // Telegram bot token
    Username   string // Bot username from API
    WebhookURL string // Webhook URL for receiving updates
}

// NewBot creates a new bot instance and configures it with Telegram.
// It validates the token, sets up webhooks, and configures commands.
func NewBot(token, webhookURL string) (*Bot, error) {
    // Implementation...
}
```

## Pull Request Process

### Before Submitting

- [ ] Tests pass locally
- [ ] Code is formatted and linted
- [ ] Documentation is updated
- [ ] Commit messages follow convention
- [ ] Branch is up to date with main

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other (describe)

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing performed
- [ ] All tests pass

## Screenshots
(If applicable)

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] No breaking changes (or documented)

## Related Issues
Closes #123
```

### Review Process

1. **Automated Checks**: CI/CD pipeline runs tests and linting
2. **Code Review**: Maintainers review code for quality and design
3. **Feedback**: Address review comments and suggestions
4. **Approval**: Once approved, PR will be merged

## Issue Reporting

### Bug Reports

Use the bug report template:

```markdown
**Bug Description**
Clear description of the bug

**Steps to Reproduce**
1. Go to...
2. Click on...
3. See error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., Ubuntu 20.04]
- Go version: [e.g., 1.24.5]
- Browser: [e.g., Chrome 95]

**Additional Context**
Screenshots, logs, etc.
```

### Feature Requests

```markdown
**Feature Description**
Clear description of the feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should it work?

**Alternatives**
Other approaches considered

**Additional Context**
Mockups, examples, etc.
```

## Community

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions and reviews

### Getting Help

- Check existing issues and documentation
- Search closed issues for similar problems
- Ask questions in GitHub Discussions
- Be specific and provide context

### Recognition

Contributors are recognized in:
- README.md contributors section
- Release notes
- GitHub contributor statistics

## Development Tips

### Debugging

```bash
# Debug with Delve
dlv debug cmd/server/main.go

# View logs
docker-compose logs -f chatmood

# Database access
docker-compose exec postgres psql -U chatmood -d chatmood
```

### Performance

- Profile with `go tool pprof`
- Monitor memory usage
- Optimize database queries
- Cache frequently accessed data

### Security

- Never commit secrets or tokens
- Validate all user inputs
- Use HTTPS in production
- Follow security best practices

## Release Process

Releases are handled by maintainers:

1. Version bump in go.mod
2. Update CHANGELOG.md
3. Create release tag
4. GitHub Actions builds and publishes
5. Docker images updated
6. Documentation updated

## Questions?

If you have questions about contributing:

1. Check this document first
2. Search existing issues
3. Create a new issue with the `question` label
4. Join GitHub Discussions

Thank you for contributing to ChatMood! 🎭