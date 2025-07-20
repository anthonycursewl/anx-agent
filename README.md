<div align="center">
  <img src="https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-brightgreen?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/github/actions/workflow/status/anthonycursewl/anx-agent/go.yml?style=for-the-badge" alt="Build Status">
  <img src="https://img.shields.io/github/go-mod/go-version/anthonycursewl/anx-agent?style=for-the-badge" alt="Go Version">
  
  <h1>🤖 ANX Agent</h1>
  <h3>Build Smarter CLI Tools with AI-Powered Automation</h3>
  
  <p align="center">
    <a href="#key-features">Features</a> •
    <a href="#-getting-started">Quick Start</a> •
    <a href="#-project-architecture">Architecture</a> •
    <a href="#-development">Development</a>
  </p>
  
</div>

## 🚀 Introduction

**ANX Agent** is a powerful, modular framework for building intelligent command-line applications with AI capabilities. Built in Go, it provides a solid foundation for creating CLI tools that can understand, process, and respond to complex commands using state-of-the-art AI models.

```bash
# Example usage (soon)
anx-agent analyze --path ./my-project --report-format markdown
```

## ✨ Key Features

<div align="center">
  <table>
    <tr>
      <td align="center">
        <img src="https://img.icons8.com/color/48/000000/expand-arrow--v1.png" width="42" height="42"/>
        <h4>Extensible Architecture</h4>
        <p>Easily add new AI providers, commands, and functionalities</p>
      </td>
      <td align="center">
        <img src="https://img.icons8.com/color/48/000000/artificial-intelligence.png" width="42" height="42"/>
        <h4>Multi-AI Support</h4>
        <p>Designed to work with multiple AI services, starting with Gemini</p>
      </td>
      <td align="center">
        <img src="https://img.icons8.com/color/48/000000/settings-3.png" width="42" height="42"/>
        <h4>Simple Configuration</h4>
        <p>Configure via YAML or environment variables</p>
      </td>
    </tr>
  </table>
</div>

## 🛠 Installation

### Prerequisites

- Go 1.20 or higher
- A Gemini API key for AI functionality (or any other AI provider)

### Quick Install

```bash
# Clone the repository
git clone https://github.com/anthonycursewl/anx-agent.git
cd anx-agent

# Install dependencies
go mod download

# Build the binary
go build -o bin/anx-agent ./cmd/agentcli

# Add to PATH (optional)
export PATH=$PATH:$(pwd)/bin
```

## ⚙️ Configuration

### Option 1: YAML Configuration

Create a `config.yaml` file:

```yaml
# config.yaml
gemini_api_key: "your-api-key-here"
log_level: "info"
max_retries: 3
timeout: "30s"
```

### Option 2: Environment Variables

```bash
export GEMINI_API_KEY="your-api-key-here"
export ANX_LOG_LEVEL="debug"
```

## 🚀 Usage Examples

### Basic Usage
```bash
# Start interactive mode
anx-agent

# Analyze a specific directory
anx-agent analyze --path ./src --output report.md

# Get help
anx-agent --help
```

### Advanced Usage
```bash
# Use a custom config file
anx-agent --config /path/to/config.yaml

# Enable debug logging
anx-agent --log-level debug

# Process specific file types
anx-agent analyze --path . --extensions go,md,txt
```

## 🏗 Project Structure

```
.
├── cmd/                 # CLI entry points
│   └── agentcli/        # Main CLI application
├── config/              # Configuration templates
├── internal/
│   ├── agent/          # Core agent logic
│   ├── ai/             # AI integrations
│   ├── cli/            # CLI components
│   ├── config/         # Configuration handling
│   └── reporting/      # Output formatters
└── testdata/           # Test fixtures
```

## 🧩 Extending ANX Agent

### Adding a New AI Provider

1. Create a new package in `internal/ai/`
2. Implement the `AIClient` interface
3. Register your provider in the factory

```go
// Example: internal/ai/myprovider/client.go
type MyAIClient struct {
    // Implementation
}

func NewMyAIClient(cfg Config) *MyAIClient {
    // Initialization
}
```

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📬 Contact

Project Link: [https://github.com/anthonycursewl/anx-agent](https://github.com/anthonycursewl/anx-agent)

## 🙏 Acknowledgments

- [Go](https://golang.org/) - The programming language
- [Cobra](https://github.com/spf13/cobra) - CLI library for Go
- [Viper](https://github.com/spf13/viper) - Go configuration with fangs

<div align="center">
  Made with ❤️ by the ANX Team
</div>