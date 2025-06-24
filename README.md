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

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
