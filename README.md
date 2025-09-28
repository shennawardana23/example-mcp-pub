# example-mcp-pub

A simple Go web application with hot reload development support using Air.

## Features

- Simple HTTP server with health check endpoint
- Hot reload development with Air
- Comprehensive Makefile for development tasks
- Test coverage reporting
- Cross-platform building

## Getting Started

### Prerequisites

- Go 1.25+ (automatically downloaded if needed)
- Make

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/shennawardana23/example-mcp-pub.git
   cd example-mcp-pub
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

## Development

### Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

### Common Development Tasks

- **Start development with hot reload**: `make dev`
- **Build the application**: `make build`
- **Run tests**: `make test`
- **Run all checks**: `make check` (formats code, runs vet, and tests)
- **Clean build artifacts**: `make clean`

### Hot Reload Development

The project includes Air for hot reload development. Simply run:

```bash
make dev
```

This will:
1. Install Air if not already installed
2. Initialize Air configuration if needed
3. Start the development server with hot reload on http://localhost:8080

Any changes to Go files will automatically rebuild and restart the server.

### Testing

- Run tests: `make test`
- Run tests with coverage: `make test-coverage`
- Generate HTML coverage report: `make test-coverage-html`

### Building

- Build for current platform: `make build`
- Build for Linux: `make build-linux`

## API Endpoints

- `GET /` - Hello World with current timestamp
- `GET /health` - Health check endpoint (returns JSON)

## Configuration

Air configuration is stored in `.air.toml`. You can modify this file to customize hot reload behavior.