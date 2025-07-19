<div align="center">
<pre>
  _   _    _    _   _ ____    _    ____ _____ 
 / \ | \ | |  / \  | \ | | __ )  / \  / ___| ____|
/ _ \|  \| | / _ \ |  \| |  _ \ / _ \ \___ \|  _|
/ ___ \ | |\  |/ ___ \| |\  | |_) / ___ \ ___) | |___
/_/   \_\_| \_|_/   \_\_| \_\|____/_/   \_\____/|_____|
</pre>
<h1>ANX Agent</h1>
<p><strong>A flexible and extensible AI agent framework for building command-line tools in Go.</strong></p>

<p>
    <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.20%2B-blue.svg" alt="Go Version"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License"></a>
    <a href="https://github.com/anthonycursewl/anx-agent/actions/workflows/go.yml"><img src="https://github.com/anthonycursewl/anx-agent/actions/workflows/go.yml/badge.svg" alt="Build Status"></a>
    <a href="https://goreportcard.com/report/github.com/anthonycursewl/anx-agent"><img src="https://goreportcard.com/badge/github.com/anthonycursewl/anx-agent" alt="Go Report Card"></a>
</p>
</div>

**ANX Agent** is a robust and modular framework written in Go for building powerful AI-driven command-line interface (CLI) tools.

---

## ✨ Key Features

- **Extensible Architecture**: Easily add new AI providers, commands, and functionalities.
- **Multi-AI Support**: Initially implemented with Gemini, but designed to support multiple AI services.
- **Simple Configuration**: Configure your agent using a `config.yaml` file or environment variables.
- **CLI-Focused**: Built from the ground up to power interactive and intelligent command-line applications.
- **Centralized Logic**: The agent core coordinates all components for a cohesive workflow.

## 📂 Project Structure

```
.
├── cmd/
│   └── agentcli/         # Main CLI application entry point
├── config/               # Configuration templates and examples
├── internal/
│   ├── agent/           # Core agent logic and interfaces
│   ├── ai/              # AI client implementations (e.g., Gemini)
│   ├── cli/             # Command-line interface components
│   ├── config/          # Configuration management
│   └── reporting/       # Reporting and output formatting
├── pkg/                 # Reusable packages (if any)
├── scripts/             # Utility scripts
└── testdata/            # Test data and fixtures
```

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or higher
- A Gemini API key for AI functionality

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/anthonycursewl/anx-agent.git
    cd anx-agent
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

### Configuration

You can provide your API key in two ways:

**Option 1: `config.yaml` file**

Create a `config.yaml` file in the project root with the following content:

```yaml
gemini_api_key: "your-gemini-api-key-here"
```

**Option 2: Environment Variable**

Export the API key as an environment variable. This method takes precedence over the configuration file.

```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

### Usage

Run the agent from the project root:

```bash
# Build and run with default configuration
go run ./cmd/agentcli

# Run with a custom configuration file
go run ./cmd/agentcli --config /path/to/your/config.yaml

# Analyze a specific file or directory
go run ./cmd/agentcli --input /path/to/analyze
```

## 🔧 Project Architecture

### Core Components

-   **CLI (Command Line Interface)**: Manages all user interaction (inputs, arguments, flags), controls the command execution flow, and provides a clear, interactive user experience.
-   **AI Client**: Abstracts communication with AI services like Gemini. It handles API requests, responses, and potential errors or retries.
-   **Configuration**: Loads settings from YAML files or environment variables and validates them to ensure the agent runs correctly.
-   **Agent Core**: Implements the main business logic, coordinates communication between the CLI, AI client, and other components, and manages state and context during execution.

## 🛠️ Development

### Building

To create an executable binary:

```bash
go build -o bin/anx-agent ./cmd/agentcli
```

You can now run `./bin/anx-agent`.

### Testing

```bash
# Run all tests in the project
go test ./...

# Run tests and generate a coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Adding New Features

-   **Add a new AI Provider**:
    1.  Create a new package under `internal/ai/`.
    2.  Implement the `AIClient` interface.
    3.  Register the new provider in the corresponding factory so the agent can use it.

-   **Add a new Command**:
    1.  Add the command logic in `internal/cli/commands/`.
    2.  Register the new command in the CLI setup.
    3.  Remember to add help text and documentation for the new command.

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are always welcome! If you want to improve the project, feel free to open a Pull Request.