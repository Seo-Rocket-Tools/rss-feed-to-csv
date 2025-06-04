# RSS Feed to CSV Converter

A fast and reliable web service that converts RSS feeds to CSV format. Built with Go for high performance and easy deployment.

## Features

- 🚀 Fast RSS parsing and CSV generation
- 🔒 Input validation and security measures
- 🧹 Optional HTML sanitization
- ⚙️ Configurable via environment variables
- 📊 Structured logging
- 🛡️ Graceful shutdown
- 🧪 Well-tested codebase

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go        # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── errors/            # Custom error types
│   ├── handlers/          # HTTP request handlers
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middleware (rate limiting)
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   │   ├── rss_fetcher.go # RSS fetching logic
│   │   └── csv_exporter.go # CSV export logic
│   ├── utils/             # Utility functions
│   └── validator/         # Input validation
├── web/
│   └── index.html         # Web UI
├── build/
│   └── docker/
│       ├── Dockerfile     # Docker build configuration
│       └── docker-compose.yml # Docker compose setup
├── scripts/
│   └── dev.sh             # Development helper scripts
├── docs/                  # Documentation
├── bin/                   # Build output directory
├── Makefile               # Build and development commands
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/rss-feed-to-csv.git
cd rss-feed-to-csv

# Build and run
make build
make run

# Or use the development mode with hot reload
make dev
```

The service will start on `http://localhost:8080`

## Docker Usage

### Using Docker Compose
```bash
# Build and run with Docker Compose
docker-compose -f build/docker/docker-compose.yml up

# Run in detached mode
docker-compose -f build/docker/docker-compose.yml up -d

# Stop the service
docker-compose -f build/docker/docker-compose.yml down
```

### Using Docker directly
```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run
```

## Usage

### Web Interface
Navigate to `http://localhost:8080` and enter an RSS feed URL.

### API Endpoint
```bash
curl "http://localhost:8080/export?url=https://example.com/feed.rss&sanitize=true" -o feed.csv
```

Parameters:
- `url` (required): The RSS feed URL
- `sanitize` (optional): Set to "true" to strip HTML from content

## Configuration

Configure the application using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `:8080` |
| `READ_TIMEOUT` | HTTP read timeout | `15s` |
| `WRITE_TIMEOUT` | HTTP write timeout | `15s` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `30s` |
| `RSS_FETCH_TIMEOUT` | RSS fetch timeout | `30s` |
| `MAX_RSS_SIZE` | Maximum RSS size in bytes | `10485760` (10MB) |
| `USER_AGENT` | User agent for RSS requests | `RSS-to-CSV-Exporter/1.0` |
| `MAX_URL_LENGTH` | Maximum URL length | `2048` |
| `RATE_LIMIT_PER_MIN` | Rate limit per minute | `60` |
| `DEFAULT_SANITIZE` | Default HTML sanitization | `false` |
| `LOG_LEVEL` | Logging level | `INFO` |

## Development

### Prerequisites
- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Available Commands

```bash
make build       # Build the application
make run         # Run the application
make test        # Run tests
make coverage    # Generate test coverage report
make lint        # Run linters
make fmt         # Format code
make dev         # Run with hot reload
make clean       # Clean build artifacts
make help        # Show all available commands
```

### Development with Hot Reload

For development with automatic reloading on file changes:
```bash
# Using make
make dev

# Or directly run the development script
./scripts/dev.sh
```

## Testing

Run the test suite:
```bash
make test
```

Generate coverage report:
```bash
make coverage
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.