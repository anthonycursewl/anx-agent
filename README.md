# ANX Agent

ANX Agent is a Go-based AI agent framework that provides a flexible and extensible platform for building AI-powered command-line tools.

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
- Gemini API key (for AI functionality)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/anthonycursewl/anx-agent.git
   cd anx-agent
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a configuration file (`config.yaml`):
   ```yaml
   gemini_api_key: "your-gemini-api-key-here"
   ```

   Or set the API key via environment variable:
   ```bash
   export GEMINI_API_KEY="your-gemini-api-key-here"
   ```

### Usage

```bash
# Run with default configuration
./anx-agent

# Specify custom config file
./anx-agent --config /path/to/config.yaml

# Analyze specific file or directory
./anx-agent --input /path/to/analyze
```

## 🔧 Project Architecture

### Core Components

1. **CLI (Command Line Interface)**
   - Handles user input/output
   - Manages command execution flow
   - Provides interactive prompts

2. **AI Client**
   - Manages connections to AI services (e.g., Gemini)
   - Handles API requests and responses
   - Implements error handling and retries

3. **Configuration**
   - Loads settings from YAML files
   - Supports environment variable overrides
   - Validates configuration values

4. **Agent Core**
   - Implements the main agent logic
   - Coordinates between components
   - Manages state and context

## 🛠 Development

### Building

```bash
go build -o bin/anx-agent ./cmd/agentcli
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Adding New Features

1. **New AI Provider**
   - Create a new package under `internal/ai/`
   - Implement the `AIClient` interface
   - Register the provider in the factory

2. **New Commands**
   - Add command logic in `internal/cli/commands/`
   - Register the command in the CLI setup
   - Add help text and documentation

## 📝 License

[MIT License](LICENSE)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
