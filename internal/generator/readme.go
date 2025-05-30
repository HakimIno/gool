package generator

import (
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// generateREADME generates a comprehensive README for the project
func (g *Generator) generateREADME(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	readmeTemplate := `# {{.ProjectName}}

🚀 A modern Go {{.Framework}} application built with **Gool** - the Go project generator.

## 🛠️ Tech Stack

- **Framework**: {{.Framework | title}}
{{- if ne .ORM "none"}}
- **ORM**: {{.ORM | title}}
{{- end}}
{{- if ne .Database ""}}
- **Database**: {{.Database | title}}
{{- end}}
- **Architecture**: {{.Config.Architecture | title}}
- **Configuration**: {{.Config.Config | upper}}
{{- if ne .Config.Auth "none"}}
- **Authentication**: {{.Config.Auth | upper}}
{{- end}}
- **Logging**: {{.Config.Logging | title}}

## ✨ Features

{{- if .Config.Features.HealthCheck}}
- ❤️  Health check endpoint
{{- end}}
{{- if .Config.Features.Swagger}}
- 📚 API documentation with Swagger
{{- end}}
{{- if .Config.Middleware.CORS}}
- 🌐 CORS support
{{- end}}
{{- if .Config.Middleware.RateLimit}}
- ⏱️  Rate limiting
{{- end}}
{{- if .Config.Features.Caching}}
- ⚡ Caching support
{{- end}}
{{- if .Config.Features.Metrics}}
- 📊 Prometheus metrics
{{- end}}
{{- if .Config.Features.WebSocket}}
- 🔌 WebSocket support
{{- end}}
{{- if .Config.Docker}}
- 🐳 Docker support
{{- end}}
{{- if ne .Config.CICD "none"}}
- 🚀 CI/CD with {{.Config.CICD | title}}
{{- end}}

## 🚀 Quick Start

### Prerequisites

- Go 1.22 or higher
{{- if eq .Database "postgresql"}}
- PostgreSQL
{{- else if eq .Database "mysql"}}
- MySQL
{{- else if eq .Database "mongodb"}}
- MongoDB
{{- end}}
{{- if .Config.Features.Caching}}
- Redis (for caching)
{{- end}}

### Installation

1. **Clone and setup**
   ` + "```bash" + `
   git clone <your-repo-url>
   cd {{.ProjectName}}
   ` + "```" + `

2. **Install dependencies**
   ` + "```bash" + `
   go mod download
   ` + "```" + `

3. **Setup environment**
   ` + "```bash" + `
   cp .env.example .env
   # Edit .env with your configuration
   ` + "```" + `

{{- if .Config.Docker}}
4. **Start services with Docker**
   ` + "```bash" + `
   docker-compose up -d
   ` + "```" + `
{{- end}}

5. **Run the application**
   ` + "```bash" + `
   go run main.go
   ` + "```" + `

Your application will be available at ` + "`http://localhost:8080`" + `

## 🏗️ Project Structure

` + "```" + `
{{.ProjectName}}/
{{- if eq .Config.Architecture "simple"}}
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
{{- else if eq .Config.Architecture "clean"}}
├── cmd/                    # Application entrypoints
├── internal/
│   ├── entity/            # Business entities
│   ├── usecase/           # Business use cases
│   ├── controller/        # Controllers (handlers)
│   ├── repository/        # Data access layer
│   └── infrastructure/    # External services
├── pkg/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   └── logger/            # Logging utilities
{{- else if eq .Config.Architecture "hexagonal"}}
├── cmd/                    # Application entrypoints
├── internal/
│   ├── core/
│   │   ├── domain/        # Domain models
│   │   ├── ports/         # Interfaces
│   │   └── services/      # Domain services
│   └── adapters/
│       ├── primary/       # Driving adapters (HTTP, CLI)
│       └── secondary/     # Driven adapters (Database, Cache)
├── pkg/
│   ├── config/            # Configuration management
│   └── logger/            # Logging utilities
{{- else if eq .Config.Architecture "mvc"}}
├── cmd/                    # Application entrypoints
├── internal/
│   ├── models/            # Data models
│   ├── views/             # View templates
│   ├── controllers/       # Controllers
│   └── middleware/        # Custom middleware
├── pkg/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   └── logger/            # Logging utilities
├── static/                # Static assets
│   ├── css/
│   ├── js/
│   └── images/
└── templates/             # HTML templates
{{- end}}
├── docs/                   # Documentation
├── scripts/               # Build and deployment scripts
{{- if .Config.Docker}}
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose setup
{{- end}}
{{- if .Config.Testing}}
├── test/                  # Test files
{{- end}}
├── main.go                # Application entry point
├── Makefile              # Development commands
└── README.md             # This file
` + "```" + `

## 🔧 Development

### Available Make Commands

` + "```bash" + `
make help                 # Show all available commands
make run                  # Run the application
make build                # Build the application
make test                 # Run tests
make test-coverage        # Run tests with coverage report
make fmt                  # Format code
make lint                 # Lint code
make clean                # Clean build artifacts
{{- if .Config.Docker}}
make docker-build         # Build Docker image
make docker-run           # Run Docker container
make docker-up            # Start services with docker-compose
make docker-down          # Stop services with docker-compose
{{- end}}
make dev                  # Start development server with hot reload
make docs                 # Generate API documentation
` + "```" + `

### Hot Reload Development

Install Air for hot reloading:
` + "```bash" + `
go install github.com/cosmtrek/air@latest
make dev
` + "```" + `

## 📡 API Endpoints

{{- if .Config.Features.HealthCheck}}
### Health Check
- ` + "`GET /api/v1/health`" + ` - Health check endpoint
{{- end}}

### Users
- ` + "`GET /api/v1/users`" + ` - Get all users
- ` + "`GET /api/v1/users/:id`" + ` - Get user by ID
- ` + "`POST /api/v1/users`" + ` - Create new user
- ` + "`PUT /api/v1/users/:id`" + ` - Update user
- ` + "`DELETE /api/v1/users/:id`" + ` - Delete user

{{- if ne .Config.Auth "none"}}
### Authentication
- ` + "`POST /api/v1/auth/login`" + ` - User login
- ` + "`POST /api/v1/auth/register`" + ` - User registration
{{- if eq .Config.Auth "jwt"}}
- ` + "`POST /api/v1/auth/refresh`" + ` - Refresh JWT token
{{- end}}
{{- end}}

{{- if .Config.Features.Swagger}}
### API Documentation
Visit ` + "`http://localhost:8080/swagger/index.html`" + ` for interactive API documentation.
{{- end}}

## 🗄️ Database

{{- if ne .Database ""}}
### {{.Database | title}} Configuration

Update your ` + "`.env`" + ` file with your database credentials:

` + "```env" + `
{{- if eq .Database "postgresql"}}
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME={{.ProjectName}}_db
DB_SSLMODE=disable
{{- else if eq .Database "mysql"}}
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME={{.ProjectName}}_db
{{- else if eq .Database "sqlite"}}
DB_PATH=./{{.ProjectName}}.db
{{- end}}
` + "```" + `

{{- if eq .ORM "gorm"}}
### Migrations

Database migrations are handled automatically by GORM. To add new models:

1. Create your model in ` + "`internal/models/`" + `
2. Add the migration in ` + "`pkg/database/database.go`" + `
3. Restart the application
{{- end}}
{{- end}}

## 🧪 Testing

Run tests:
` + "```bash" + `
make test
` + "```" + `

Run tests with coverage:
` + "```bash" + `
make test-coverage
` + "```" + `

{{- if .Config.Docker}}
## 🐳 Docker

### Build and run with Docker
` + "```bash" + `
make docker-build
make docker-run
` + "```" + `

### Development with Docker Compose
` + "```bash" + `
make docker-up    # Start all services
make docker-down  # Stop all services
` + "```" + `
{{- end}}

{{- if ne .Config.CICD "none"}}
## 🚀 Deployment

This project includes {{.Config.CICD | title}} CI/CD configuration. The pipeline will:

- Run tests on every push
- Build the application
- Create Docker images
- Deploy to your chosen environment (configure as needed)

{{- if eq .Config.CICD "github"}}
Make sure to set up the following GitHub secrets:
- ` + "`DOCKERHUB_USERNAME`" + ` - Your Docker Hub username
- ` + "`DOCKERHUB_TOKEN`" + ` - Your Docker Hub access token
{{- end}}
{{- end}}

## 📊 Monitoring

{{- if .Config.Features.Metrics}}
### Prometheus Metrics
Metrics are available at ` + "`http://localhost:8080/metrics`" + `

### Grafana Dashboard
If using Docker Compose, Grafana is available at ` + "`http://localhost:3000`" + `
- Username: admin
- Password: admin
{{- end}}

{{- if .Config.Features.HealthCheck}}
### Health Checks
Monitor application health at ` + "`http://localhost:8080/api/v1/health`" + `
{{- end}}

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (` + "`git checkout -b feature/amazing-feature`" + `)
3. Commit your changes (` + "`git commit -m 'Add amazing feature'`" + `)
4. Push to the branch (` + "`git push origin feature/amazing-feature`" + `)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Gool](https://github.com/gool-cli/gool) - Go Project Generator
- Powered by {{.Framework | title}} framework
{{- if ne .ORM "none"}}
- Database access with {{.ORM | title}}
{{- end}}

---

**Generated with ❤️ by Gool** - The modern Go project scaffolding tool
`

	if err := g.templateEngine.RenderToFile(readmeTemplate, filepath.Join(projectPath, "README.md"), data); err != nil {
		return err
	}

	return nil
}
