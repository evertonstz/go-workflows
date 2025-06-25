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

- Go 1.23 or higher
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

# Run tests with verbose output
make test-verbose

# Run tests with race detection
make test-race

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
go test -v -run "TestApp" ./

# Update golden files when UI changes
go test -v -run "TestApp_FullOutput" ./ -update
```

For detailed testing information, see [TESTING.md](TESTING.md).

### CI/CD Pipeline

The project uses GitHub Actions for continuous integration:

- **Pull Request CI**: Runs on every PR

  - Linting and formatting checks
  - Complete test suite with race detection
  - Cross-platform builds (Linux, Windows, macOS)
  - Golden file verification

- **Main Branch CI**: Runs on main branch pushes and daily

  - Comprehensive test suite
  - Security scanning with `gosec`
  - Dependency vulnerability checks
  - Multi-version Go compatibility testing

- **Release CI**: Runs on version tags
  - Full test suite before release
  - Cross-compilation for multiple platforms
  - Automated releases with GoReleaser

### Code Quality

- **Linting**: Uses `golangci-lint` with comprehensive rules
- **Formatting**: Enforced `gofmt` formatting
- **Security**: Regular `gosec` security scans
- **Dependencies**: Vulnerability scanning with `govulncheck`

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
