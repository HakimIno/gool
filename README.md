# Gool - Go Project Scaffolding Tool

Gool is a modern, flexible, and comprehensive CLI tool for generating Go projects with customizable structure, pre-configured files, and dependencies. Build production-ready Go applications in seconds!

## 🚀 Features

### Core Features
- **Multiple Web Frameworks**: Choose from Gin, Echo, Fiber, or Revel
- **Database Support**: PostgreSQL, MySQL, SQLite, MongoDB, Redis, or in-memory store
- **ORM/Database Access**: GORM, sqlx, raw SQL, or none
- **Architecture Patterns**: Simple, Clean Architecture, Hexagonal (Ports & Adapters), MVC, or Custom
- **Configuration Formats**: YAML, JSON, or TOML support with environment-specific configs
- **Authentication**: JWT, OAuth2, or Basic Auth with ready-to-use templates
- **Middleware**: CORS, Rate Limiting, Logging, and Authentication middleware
- **Testing**: Unit test and integration test templates
- **Logging & Monitoring**: Standard log, Logrus, or Zap with Prometheus metrics

### Additional Features
- **WebSocket Support**: Real-time application templates
- **Error Handling**: Centralized error handling with custom error types
- **Caching**: In-memory, Redis, or Memcached support
- **Message Queues**: RabbitMQ, Kafka, or NATS for async tasks
- **Security**: HTTPS, secure headers (HSTS, CSP), and CSRF protection
- **API Documentation**: OpenAPI/Swagger generation
- **Docker & Deployment**: Dockerfile and docker-compose.yml
- **CI/CD**: GitHub Actions and GitLab CI templates
- **Health Checks**: Built-in /health endpoint
- **Internationalization**: Multi-language support
- **Cloud Integration**: AWS, GCP, and Azure deployment configs

## 📦 Installation

### Option 1: Install from source
```bash
# Clone repository
git clone https://github.com/HakimIno/gool.git
cd gool

# Build binary
go build -o gool main.go

# Install to system PATH (macOS/Linux)
sudo mv gool /usr/local/bin/

# Or for Windows
# Move gool.exe to a directory in your PATH
```

### Option 2: Go Install (when published)
```bash
go install github.com/HakimIno/gool@latest
```

### Option 3: Download binary (when available)
Download the latest binary from the [releases page](https://github.com/HakimIno/gool/releases)

## 🛠️ Development Setup

### For Developers (Git Repository Setup)

If you want to create your own repository:

```bash
# 1. Create a new repository on GitHub
# 2. Clone or setup local repository
git init
git add .
git commit -m "Initial commit"

# 3. Add remote repository
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git

# 4. Push to GitHub
git branch -M main
git push -u origin main
```

### Build and Testing

```bash
# Install dependencies
go mod tidy

# Build binary
go build -o gool main.go

# Test commands
./gool --help
./gool version
```

## 🎯 Quick Start

### Interactive Mode (Recommended)
```bash
gool init my-awesome-app
```

This will prompt you through a series of questions to customize your project:
- Project name and module path
- Web framework selection
- Database and ORM preferences
- Architecture pattern
- Authentication method
- Additional features and middleware

### Non-Interactive Mode
```bash
gool init my-app --framework=gin --orm=gorm --database=postgresql --arch=simple
```

### Available Options
```bash
# Framework options
--framework=gin|echo|fiber|revel

# ORM options  
--orm=gorm|sqlx|raw|none

# Database options
--database=postgresql|mysql|sqlite|mongodb|redis|memory

# Architecture options
--arch=simple|clean|hexagonal|mvc|custom
```

## 📂 Generated Project Structure

### Simple Architecture
```
my-app/
├── cmd/                    # Application entrypoints
├── internal/
│   ├── app/               # Application initialization
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   └── middleware/        # Custom middleware
├── pkg/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   └── logger/            # Logging utilities
├── api/
│   └── routes/            # Route definitions
├── docs/                   # Documentation
├── scripts/               # Build and deployment scripts
├── deployments/           # Deployment configs
├── test/                  # Test files
├── static/                # Static assets (if enabled)
├── locales/               # Internationalization files
├── main.go                # Application entry point
├── go.mod                 # Go modules
├── .env                   # Environment variables
├── .gitignore            # Git ignore rules
├── Dockerfile            # Docker configuration
└── README.md             # Project documentation
```

### Clean Architecture
```
my-app/
├── cmd/
├── internal/
│   ├── entity/           # Business entities
│   ├── usecase/          # Business use cases
│   ├── controller/       # Controllers (handlers)
│   ├── repository/       # Data access layer
│   └── infrastructure/   # External services
├── pkg/
└── ...
```

### Hexagonal Architecture
```
my-app/
├── cmd/
├── internal/
│   ├── core/
│   │   ├── domain/       # Domain models
│   │   ├── ports/        # Interfaces
│   │   └── services/     # Domain services
│   └── adapters/
│       ├── primary/      # Driving adapters (HTTP, CLI)
│       └── secondary/    # Driven adapters (Database, Cache)
├── pkg/
└── ...
```

## 🔧 Usage Examples

### Generate a simple web API
```bash
gool init my-api \
  --framework=gin \
  --orm=gorm \
  --database=postgresql \
  --arch=simple
```

### Generate a microservice with clean architecture
```bash
gool init my-service \
  --framework=echo \
  --orm=sqlx \
  --database=mysql \
  --arch=clean
```

### Generate a full-stack application
```bash
# Interactive mode will ask about:
# - Static file serving
# - Template engine
# - WebSocket support
# - Authentication
# - And more...
gool init my-fullstack-app
```

## 🏃‍♂️ Running Your Generated Project

After generating your project:

```bash
cd my-app

# Install dependencies
go mod tidy

# Copy and configure environment
cp .env.example .env
# Edit .env with your database credentials

# Run the application
go run main.go

# Or build and run
go build -o app main.go
./app
```

### With Docker
```bash
# Start services (database, redis, etc.)
docker-compose up -d

# Build and run your app
docker build -t my-app .
docker run -p 8080:8080 my-app
```

## 🧪 Testing

Generated projects include test templates:

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./test/...
```

## 📖 API Documentation

If you enabled Swagger documentation, visit:
```
http://localhost:8080/swagger/index.html
```

## 🎛️ Configuration

Generated projects support environment-based configuration:

```bash
# Development
APP_ENV=development go run main.go

# Production  
APP_ENV=production go run main.go

# Test
APP_ENV=test go test ./...
```

## 🔍 Health Checks

Built-in health check endpoint:
```bash
curl http://localhost:8080/api/v1/health
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management
- [Survey](https://github.com/AlecAivazis/survey) for interactive prompts
- The Go community for amazing libraries and tools

## 📞 Support

- 📚 [Documentation](https://github.com/HakimIno/gool/wiki)
- 🐛 [Issue Tracker](https://github.com/HakimIno/gool/issues)
- 💬 [Discussions](https://github.com/HakimIno/gool/discussions)

---

Made with ❤️ by the Gool team
