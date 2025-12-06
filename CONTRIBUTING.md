# Contributing to Auth Service

Thank you for your interest in contributing to Auth Service! We welcome contributions from the community to help improve this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Development Setup](#development-setup)
  - [Prerequisites](#prerequisites)
  - [Fork and Clone the Repository](#fork-and-clone-the-repository)
  - [Building the Project](#building-the-project)
  - [Running Tests](#running-tests)
- [Coding Standards](#coding-standards)
  - [Go Code Style](#go-code-style)
  - [Commit Messages](#commit-messages)
  - [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [License](#license)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the [existing issues](https://github.com/Glitch-guy0/authService/issues) to see if the problem has already been reported.

When creating a bug report, please include the following information:

1. A clear and descriptive title
2. Steps to reproduce the issue
3. Expected behavior
4. Actual behavior
5. Environment details (OS, Go version, etc.)
6. Any relevant logs or error messages

### Suggesting Enhancements

We welcome suggestions for new features and improvements. Before suggesting a new feature:

1. Check if a similar feature has already been suggested
2. Explain why this enhancement would be useful
3. Provide as much detail as possible about the proposed changes

### Your First Code Contribution

Looking for your first contribution? Look for issues labeled `good first issue` in the [issue tracker](https://github.com/Glitch-guy0/authService/issues).

### Pull Requests

1. Fork the repository and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure the test suite passes.
4. Make sure your code lints.
5. Update the documentation if necessary.
6. Issue a Pull Request.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional but recommended)
- Docker (for containerized development)

### Fork and Clone the Repository

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/authService.git
   cd authService
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/Glitch-guy0/authService.git
   ```

### Building the Project

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Build the project:
   ```bash
   make build
   ```

### Running Tests

Run the test suite:

```bash
make test
```

For integration tests:

```bash
make test-integration
```

## Coding Standards

### Go Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` or `goimports` to format your code
- Keep functions small and focused on a single responsibility
- Write clear and concise comments for exported functions and types
- Follow the project's existing code style

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally
- Consider starting the commit message with an applicable emoji:
  - ✨ `:sparkles:` When adding a new feature
  - 🐛 `:bug:` When fixing a bug
  - 🔧 `:wrench:` When changing configuration
  - 📝 `:memo:` When writing docs
  - ♻️ `:recycle:` When refactoring code
  - 🚀 `:rocket:` When improving performance
  - 🎨 `:art:` When improving the format/structure of the code
  - 🔥 `:fire:` When removing code or files
  - ✅ `:white_check_mark:` When adding tests
  - 🔒 `:lock:` When dealing with security
  - ⬆️ `:arrow_up:` When upgrading dependencies

### Documentation

- Update the relevant documentation when making changes
- Add comments to explain complex logic
- Keep API documentation up to date
- Document any breaking changes

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations, and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. You may merge the Pull Request once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## License

By contributing, you agree that your contributions will be licensed under its MIT License.
