package generator

import (
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// generateDockerFiles generates Docker-related files
func (g *Generator) generateDockerFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate Dockerfile
	dockerfileTemplate := `# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy config files if they exist
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
`

	if err := g.templateEngine.WriteFile(filepath.Join(projectPath, "Dockerfile"), dockerfileTemplate); err != nil {
		return err
	}

	// Generate docker-compose.yml
	dockerComposeTemplate := `version: '3.8'

services:
  {{.ProjectName}}:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      {{- if ne .Database ""}}
      {{- if eq .Database "postgresql"}}
      - DB_HOST=postgres
      {{- else if eq .Database "mysql"}}
      - DB_HOST=mysql
      {{- else if eq .Database "mongodb"}}
      - DB_HOST=mongodb
      {{- end}}
      {{- end}}
      {{- if .Config.Features.Caching}}
      - REDIS_HOST=redis
      {{- end}}
    depends_on:
      {{- if eq .Database "postgresql"}}
      - postgres
      {{- else if eq .Database "mysql"}}
      - mysql
      {{- else if eq .Database "mongodb"}}
      - mongodb
      {{- end}}
      {{- if .Config.Features.Caching}}
      - redis
      {{- end}}
    volumes:
      - ./logs:/root/logs
    networks:
      - {{.ProjectName}}_network

  {{- if eq .Database "postgresql"}}
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: {{.ProjectName}}_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - {{.ProjectName}}_network

  {{- else if eq .Database "mysql"}}
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: {{.ProjectName}}_db
      MYSQL_ROOT_PASSWORD: root
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - {{.ProjectName}}_network

  {{- else if eq .Database "mongodb"}}
  mongodb:
    image: mongo:6.0
    environment:
      MONGO_INITDB_DATABASE: {{.ProjectName}}_db
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    networks:
      - {{.ProjectName}}_network
  {{- end}}

  {{- if .Config.Features.Caching}}
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{.ProjectName}}_network
  {{- end}}

  {{- if .Config.Features.Metrics}}
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./deployments/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - {{.ProjectName}}_network

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - {{.ProjectName}}_network
  {{- end}}

volumes:
  {{- if eq .Database "postgresql"}}
  postgres_data:
  {{- else if eq .Database "mysql"}}
  mysql_data:
  {{- else if eq .Database "mongodb"}}
  mongodb_data:
  {{- end}}
  {{- if .Config.Features.Caching}}
  redis_data:
  {{- end}}
  {{- if .Config.Features.Metrics}}
  prometheus_data:
  grafana_data:
  {{- end}}

networks:
  {{.ProjectName}}_network:
    driver: bridge
`

	if err := g.templateEngine.RenderToFile(dockerComposeTemplate, filepath.Join(projectPath, "docker-compose.yml"), data); err != nil {
		return err
	}

	// Generate .dockerignore
	dockerignoreContent := `# Ignore everything
*

# Allow files and directories
!cmd/
!internal/
!pkg/
!api/
!docs/
!go.mod
!go.sum
!main.go
!.env.example

# Ignore specific files and directories
.git
.env
.env.local
*.log
*.db
*.sqlite
*.sqlite3
tmp/
temp/
vendor/
node_modules/
.vscode/
.idea/
*.swp
*.swo
.DS_Store
Thumbs.db
`

	if err := g.templateEngine.WriteFile(filepath.Join(projectPath, ".dockerignore"), dockerignoreContent); err != nil {
		return err
	}

	return nil
}

// generateCICDFiles generates CI/CD pipeline files
func (g *Generator) generateCICDFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	if cfg.CICD == config.CICDGitHub {
		// Generate GitHub Actions workflow - use raw string to avoid template conflicts
		githubWorkflowContent := `name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'

    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'

    - name: Build
      run: go build -v ./...

    - name: Build for Linux
      run: GOOS=linux GOARCH=amd64 go build -o ` + cfg.ProjectName + `-linux-amd64 .

    - name: Build for Windows
      run: GOOS=windows GOARCH=amd64 go build -o ` + cfg.ProjectName + `-windows-amd64.exe .

    - name: Build for macOS
      run: GOOS=darwin GOARCH=amd64 go build -o ` + cfg.ProjectName + `-darwin-amd64 .
`

		workflowDir := filepath.Join(projectPath, ".github/workflows")
		if err := g.templateEngine.WriteFile(filepath.Join(workflowDir, "ci-cd.yml"), githubWorkflowContent); err != nil {
			return err
		}

	} else if cfg.CICD == config.CICDGitLab {
		// Generate GitLab CI configuration
		gitlabCITemplate := `stages:
  - test
  - build

variables:
  GO_VERSION: "1.22"

before_script:
  - apt-get update -qq && apt-get install -y -qq git ca-certificates
  - go version

test:
  stage: test
  image: golang:$GO_VERSION
  script:
    - go mod download
    - go test -v -race -coverprofile=coverage.out ./...
  coverage: '/coverage: \d+\.\d+% of statements/'

build:
  stage: build
  image: golang:$GO_VERSION
  script:
    - go mod download
    - go build -o {{.ProjectName}} .
  artifacts:
    paths:
      - {{.ProjectName}}
    expire_in: 1 week
  only:
    - main
    - develop
`

		if err := g.templateEngine.RenderToFile(gitlabCITemplate, filepath.Join(projectPath, ".gitlab-ci.yml"), data); err != nil {
			return err
		}
	}

	// Generate Makefile
	makefileTemplate := `# {{.ProjectName}} Makefile

.PHONY: build run test clean docker-build docker-run deps fmt vet lint help

# Variables
APP_NAME={{.ProjectName}}
GO_VERSION=1.22
DOCKER_IMAGE={{.ProjectName}}:latest

# Build the application
build:
	@echo "Building {{.ProjectName}}..."
	@go build -o bin/$(APP_NAME) .

# Run the application
run:
	@echo "Running {{.ProjectName}}..."
	@go run .

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker compose up
docker-up:
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping services with docker-compose..."
	@docker-compose down

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@cp .env.example .env
	@go mod download
	@echo "Setup complete! Edit .env file with your configuration."

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start services with docker-compose"
	@echo "  docker-down   - Stop services with docker-compose"
	@echo "  setup         - Setup development environment"
	@echo "  help          - Show this help message"
`

	if err := g.templateEngine.RenderToFile(makefileTemplate, filepath.Join(projectPath, "Makefile"), data); err != nil {
		return err
	}

	return nil
}
