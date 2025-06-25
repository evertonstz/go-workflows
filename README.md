<div align="center">
  <img src="https://github.com/user-attachments/assets/e9288a0d-4c15-4543-a745-822fa58f2b13" alt="gopher-workflows" width="400" />
</div>

## Overview

`go-workflows` is a terminal-based application designed to simplify your daily workflow by providing a collection of snippets and commands that you use frequently. It features a simple yet elegant Text User Interface (TUI) and stores all data in a JSON file, making it easy to version and sync using tools like `yadm` and `chezmoi`.

## Features

- **Command and Snippet Management**: Quickly access and manage your frequently used commands and snippets.
- **Terminal-Based TUI**: A user-friendly interface that runs directly in the terminal.
- **JSON Storage**: All data is stored in a JSON file, enabling easy versioning and syncing.
- **Integration with Sync Tools**: Compatible with tools like `yadm` and `chezmoi` for seamless synchronization across devices.

## Installation (Homebrew)

1. Add the tap for `go-workflows`:
   ```bash
   brew tap evertonstz/go-workflows
   ```
2. Install the application:
   ```bash
   brew install go-workflows
   ```

## Usage

To start the application, run the following command:

```bash
go-workflows
```

Open your terminal to interact with the TUI and manage your snippets and commands.

## Development

### Prerequisites

- Go 1.24.4 or higher
- Make (for using the Makefile)

### Building from Source

```bash
git clone https://github.com/evertonstz/go-workflows.git
cd go-workflows
make build
```

### Testing

This project has comprehensive test coverage including unit tests and Bubble Tea integration tests using `teatest`.

#### Quick Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Show coverage summary in terminal
make test-cover-summary

# Run tests with verbose output
make test-verbose

# Run tests with race detection
make test-race

# Run integration tests (Bubble Tea with teatest)
make test-integration

# Update golden files for integration tests
make test-integration-update

# Generate HTML coverage report
make test-cover-html
```

#### Test Categories

- **Unit Tests**: Business logic, services, data models
- **Integration Tests**: Complete Bubble Tea application behavior using `teatest`
- **Golden File Tests**: UI output verification and regression detection

#### Bubble Tea Testing with teatest

We use the experimental `teatest` package to test the complete Bubble Tea application:

```bash
# Run Bubble Tea integration tests
make test-integration

# Update golden files when UI changes
make test-integration-update

# Or use go test directly
go test -v -run "TestApp" ./
go test -v -run "TestApp_FullOutput" ./ -update
```

### CI/CD Pipeline

The project uses GitHub Actions for continuous integration with **native GitHub coverage reporting** (no external services like Codecov needed):

- **Pull Request CI**: Runs on every PR

  - Linting and formatting checks (using `goimports`)
  - Complete test suite with race detection
  - Cross-platform builds (Linux, Windows, macOS)
  - Golden file verification
  - **Automated coverage reporting**: Creates/updates PR comments with coverage statistics
  - **Coverage artifacts**: Downloadable HTML and profile reports (30-day retention)

- **Main Branch CI**: Runs on main branch pushes and daily
  - Comprehensive test suite
  - Code formatting validation with `goimports`
  - Dependency vulnerability checks
  - Multi-version Go compatibility testing
  - Coverage summaries: Integrated into GitHub Actions workflow summaries
  - Coverage artifacts: Long-term storage (90-day retention)

- **Release CI**: Runs on version tags
  - Full test suite before release
  - Cross-compilation for multiple platforms
  - Automated releases with GoReleaser

#### Coverage Visualization

The project uses GitHub's built-in features for coverage visualization:

- **PR Comments**: Automated coverage reports with package-by-package breakdown
- **Workflow Summaries**: Coverage statistics in GitHub Actions summary pages
- **Artifacts**: Downloadable HTML reports and coverage profiles
- **Security**: All coverage data stays within GitHub's secure infrastructure
- **Zero Cost**: No external service subscriptions or API limits

### Code Quality

- **Linting**: Uses `golangci-lint` with comprehensive rules
- **Formatting**: Enforced `goimports` formatting with local package prioritization
- **Dependencies**: Vulnerability scanning with `govulncheck`

### Make Commands

```bash
# Code quality
make format          # Format code with goimports
make format-check    # Check if code is properly formatted
make lint           # Run golangci-lint
make ci             # Run full CI pipeline (format-check, lint, test-race, test-cover)
```

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
