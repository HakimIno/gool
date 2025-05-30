package generator

import (
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// generateCoreFiles generates the main files for the project
func (g *Generator) generateCoreFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate go.mod
	goModTemplate := `module {{.ModulePath}}

go 1.22

require (
{{- if eq .Framework "gin"}}
	github.com/gin-gonic/gin v1.9.1
{{- else if eq .Framework "echo"}}
	github.com/labstack/echo/v4 v4.11.4
{{- else if eq .Framework "fiber"}}
	github.com/gofiber/fiber/v2 v2.52.0
{{- else if eq .Framework "revel"}}
	github.com/revel/revel v1.1.0
{{- end}}
{{- if eq .ORM "gorm"}}
	gorm.io/gorm v1.25.5
	{{- if eq .Database "postgresql"}}
	gorm.io/driver/postgres v1.5.4
	{{- else if eq .Database "mysql"}}
	gorm.io/driver/mysql v1.5.2
	{{- else if eq .Database "sqlite"}}
	gorm.io/driver/sqlite v1.5.4
	{{- end}}
{{- else if eq .ORM "sqlx"}}
	github.com/jmoiron/sqlx v1.3.5
	{{- if eq .Database "postgresql"}}
	github.com/lib/pq v1.10.9
	{{- else if eq .Database "mysql"}}
	github.com/go-sql-driver/mysql v1.7.1
	{{- else if eq .Database "sqlite"}}
	github.com/mattn/go-sqlite3 v1.14.18
	{{- end}}
{{- end}}
{{- if eq .Config.Logging "zap"}}
	go.uber.org/zap v1.26.0
{{- else if eq .Config.Logging "logrus"}}
	github.com/sirupsen/logrus v1.9.3
{{- else if eq .Config.Logging "charm"}}
	github.com/charmbracelet/log v0.4.0
	github.com/charmbracelet/lipgloss v0.9.1
{{- end}}
{{- if eq .Config.Auth "jwt"}}
	github.com/golang-jwt/jwt/v5 v5.2.0
{{- end}}
{{- if .Config.Features.Swagger}}
	github.com/swaggo/swag v1.16.2
	{{- if eq .Framework "gin"}}
	github.com/swaggo/gin-swagger v1.6.0
	{{- else if eq .Framework "echo"}}
	github.com/swaggo/echo-swagger v1.4.1
	{{- else if eq .Framework "fiber"}}
	github.com/swaggo/fiber-swagger v1.3.0
	{{- end}}
{{- end}}
{{- if .Config.Features.HealthCheck}}
	github.com/heptiolabs/healthcheck v0.0.0-20211123025425-613501dd5deb
{{- end}}
{{- if .Config.Features.Metrics}}
	github.com/prometheus/client_golang v1.17.0
{{- end}}
	github.com/spf13/viper v1.18.2
	github.com/joho/godotenv v1.5.1
)
`

	if err := g.templateEngine.RenderToFile(goModTemplate, filepath.Join(projectPath, "go.mod"), data); err != nil {
		return err
	}

	// Generate main.go
	mainTemplate := `package main

import (
	"log"
	"{{.ModulePath}}/internal/app"
)

// @title {{.ProjectName}} API
// @version 1.0
// @description A {{.Framework}} web service
// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
`

	if err := g.templateEngine.RenderToFile(mainTemplate, filepath.Join(projectPath, "main.go"), data); err != nil {
		return err
	}

	// Generate .env file
	envTemplate := `# Application Configuration
APP_NAME={{.ProjectName}}
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database Configuration
{{- if ne .Database ""}}
{{- if eq .Database "postgresql"}}
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME={{.ProjectName}}_db
DB_SSLMODE=disable
{{- else if eq .Database "mysql"}}
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME={{.ProjectName}}_db
{{- else if eq .Database "sqlite"}}
DB_PATH=./{{.ProjectName}}.db
{{- end}}
{{- end}}

# JWT Configuration
{{- if eq .Config.Auth "jwt"}}
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
{{- end}}

# Redis Configuration (if using cache)
{{- if .Config.Features.Caching}}
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
{{- end}}

# Log Configuration
LOG_LEVEL=info
{{- if eq .Config.Logging "zap"}}
LOG_FORMAT=json
{{- else if eq .Config.Logging "charm"}}
LOG_FORMAT=text
{{- end}}
`

	if err := g.templateEngine.RenderToFile(envTemplate, filepath.Join(projectPath, ".env"), data); err != nil {
		return err
	}

	// Generate .env.example
	if err := g.templateEngine.RenderToFile(envTemplate, filepath.Join(projectPath, ".env.example"), data); err != nil {
		return err
	}

	// Generate .gitignore
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Application files
*.log
*.pid
tmp/
temp/

# Database files
*.db
*.sqlite
*.sqlite3

# Build output
build/
dist/
`

	if err := g.templateEngine.WriteFile(filepath.Join(projectPath, ".gitignore"), gitignoreContent); err != nil {
		return err
	}

	return nil
}

// generateFrameworkFiles generates framework-specific files
func (g *Generator) generateFrameworkFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	switch cfg.Framework {
	case config.FrameworkGin:
		return g.generateGinFiles(cfg, projectPath, data)
	case config.FrameworkEcho:
		return g.generateEchoFiles(cfg, projectPath, data)
	case config.FrameworkFiber:
		return g.generateFiberFiles(cfg, projectPath, data)
	default:
		return g.generateGinFiles(cfg, projectPath, data) // Default to Gin
	}
}

// generateGinFiles generates Gin-specific files
func (g *Generator) generateGinFiles(cfg *config.ProjectConfig, projectPath string, data *templates.TemplateData) error {
	// Generate app/app.go
	appTemplate := `package app

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"{{.ModulePath}}/api/routes"
	"{{.ModulePath}}/pkg/config"
	{{- if ne .ORM "none"}}
	"{{.ModulePath}}/pkg/database"
	{{- end}}
	pkgLogger "{{.ModulePath}}/pkg/logger"
	"{{.ModulePath}}/pkg/startup"
	{{- if .Config.Features.Swagger}}
	"{{.ModulePath}}/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	{{- end}}
	{{- if eq .Config.Logging "zap"}}
	"go.uber.org/zap"
	{{- end}}
)

type App struct {
	router *gin.Engine
	config *config.Config
}

func New() *App {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	pkgLogger.Init(cfg.Log.Level, cfg.Log.Format)

	{{- if ne .ORM "none"}}
	// Initialize database
	if err := database.Init(cfg); err != nil {
		{{- if eq .Config.Logging "zap"}}
		pkgLogger.Fatal("Failed to initialize database", zap.Error(err))
		{{- else}}
		pkgLogger.Fatal("Failed to initialize database", "error", err)
		{{- end}}
	}
	{{- end}}

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := &App{
		router: gin.New(),
		config: cfg,
	}

	app.setupMiddleware()
	app.setupRoutes()

	return app
}

func (a *App) setupMiddleware() {
	// Recovery middleware
	a.router.Use(gin.Recovery())

	{{- if .Config.Middleware.Logging}}
	// Logging middleware
	a.router.Use(gin.Logger())
	{{- end}}

	{{- if .Config.Middleware.CORS}}
	// CORS middleware
	a.router.Use(gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}))
	{{- end}}
}

func (a *App) setupRoutes() {
	{{- if .Config.Features.Swagger}}
	// Swagger documentation
	docs.SwaggerInfo.Title = "{{.ProjectName}} API"
	docs.SwaggerInfo.Description = "{{.ProjectName}} API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{{- end}}

	// Setup API routes
	routes.SetupRoutes(a.router)
}

func (a *App) Run() error {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Show beautiful startup message
	startup.ShowWelcome("{{.ProjectName}}", "1.0", port)
	
	return a.router.Run(":" + port)
}
`

	if err := g.templateEngine.RenderToFile(appTemplate, filepath.Join(projectPath, "internal/app/app.go"), data); err != nil {
		return err
	}

	return nil
}

// generateEchoFiles generates Echo-specific files
func (g *Generator) generateEchoFiles(cfg *config.ProjectConfig, projectPath string, data *templates.TemplateData) error {
	// Generate app/app.go for Echo
	appTemplate := `package app

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"{{.ModulePath}}/api/routes"
	"{{.ModulePath}}/pkg/config"
	{{- if ne .ORM "none"}}
	"{{.ModulePath}}/pkg/database"
	{{- end}}
	pkgLogger "{{.ModulePath}}/pkg/logger"
	"{{.ModulePath}}/pkg/startup"
	{{- if .Config.Features.Swagger}}
	"{{.ModulePath}}/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
	{{- end}}
	"go.uber.org/zap"
)

type App struct {
	echo   *echo.Echo
	config *config.Config
}

func New() *App {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	pkgLogger.Init(cfg.Log.Level, cfg.Log.Format)

	{{- if ne .ORM "none"}}
	// Initialize database
	if err := database.Init(cfg); err != nil {
		pkgLogger.Fatal("Failed to initialize database", zap.Error(err))
	}
	{{- end}}

	e := echo.New()

	app := &App{
		echo:   e,
		config: cfg,
	}

	app.setupMiddleware()
	app.setupRoutes()

	return app
}

func (a *App) setupMiddleware() {
	// Recovery middleware
	a.echo.Use(middleware.Recover())

	{{- if .Config.Middleware.Logging}}
	// Logging middleware
	a.echo.Use(middleware.Logger())
	{{- end}}

	{{- if .Config.Middleware.CORS}}
	// CORS middleware
	a.echo.Use(middleware.CORS())
	{{- end}}

	{{- if .Config.Middleware.RateLimit}}
	// Rate limiting middleware
	a.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	{{- end}}
}

func (a *App) setupRoutes() {
	{{- if .Config.Features.Swagger}}
	// Swagger documentation
	docs.SwaggerInfo.Title = "{{.ProjectName}} API"
	docs.SwaggerInfo.Description = "{{.ProjectName}} API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	a.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	{{- end}}

	// Setup API routes
	routes.SetupRoutes(a.echo)
}

func (a *App) Run() error {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Show beautiful startup message
	startup.ShowWelcome("{{.ProjectName}}", "1.0", port)
	
	return a.echo.Start(":" + port)
}
`

	if err := g.templateEngine.RenderToFile(appTemplate, filepath.Join(projectPath, "internal/app/app.go"), data); err != nil {
		return err
	}

	return nil
}

// generateFiberFiles generates Fiber-specific files
func (g *Generator) generateFiberFiles(cfg *config.ProjectConfig, projectPath string, data *templates.TemplateData) error {
	// Generate app/app.go for Fiber
	appTemplate := `package app

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"{{.ModulePath}}/api/routes"
	"{{.ModulePath}}/pkg/config"
	{{- if ne .ORM "none"}}
	"{{.ModulePath}}/pkg/database"
	{{- end}}
	pkgLogger "{{.ModulePath}}/pkg/logger"
	"{{.ModulePath}}/pkg/startup"
	{{- if .Config.Features.Swagger}}
	"{{.ModulePath}}/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	{{- end}}
	"go.uber.org/zap"
)

type App struct {
	fiber  *fiber.App
	config *config.Config
}

func New() *App {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	pkgLogger.Init(cfg.Log.Level, cfg.Log.Format)

	{{- if ne .ORM "none"}}
	// Initialize database
	if err := database.Init(cfg); err != nil {
		pkgLogger.Fatal("Failed to initialize database", zap.Error(err))
	}
	{{- end}}

	fiberApp := fiber.New(fiber.Config{
		AppName: "{{.ProjectName}}",
		DisableStartupMessage: true,
	})

	app := &App{
		fiber:  fiberApp,
		config: cfg,
	}

	app.setupMiddleware()
	app.setupRoutes()

	return app
}

func (a *App) setupMiddleware() {
	// Recovery middleware
	a.fiber.Use(recover.New())

	{{- if .Config.Middleware.Logging}}
	// Logging middleware
	a.fiber.Use(logger.New())
	{{- end}}

	{{- if .Config.Middleware.CORS}}
	// CORS middleware
	a.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))
	{{- end}}
}

func (a *App) setupRoutes() {
	{{- if .Config.Features.Swagger}}
	// Swagger documentation
	docs.SwaggerInfo.Title = "{{.ProjectName}} API"
	docs.SwaggerInfo.Description = "{{.ProjectName}} API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	a.fiber.Get("/swagger/*", fiberSwagger.WrapHandler)
	{{- end}}

	// Setup API routes
	routes.SetupRoutes(a.fiber)
}

func (a *App) Run() error {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Show beautiful startup message
	startup.ShowWelcome("{{.ProjectName}}", "1.0", port)
	
	return a.fiber.Listen(":" + port)
}
`

	if err := g.templateEngine.RenderToFile(appTemplate, filepath.Join(projectPath, "internal/app/app.go"), data); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateMiddlewareFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate auth middleware
	if cfg.Auth == config.AuthJWT {
		authMiddlewareTemplate := `package middleware

import (
	"net/http"
	"strings"

	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
	"github.com/golang-jwt/jwt/v5"
	"{{.ModulePath}}/pkg/logger"
	"go.uber.org/zap"
)

type Claims struct {
	UserID uint   ` + "`" + `json:"user_id"` + "`" + `
	Email  string ` + "`" + `json:"email"` + "`" + `
	Role   string ` + "`" + `json:"role"` + "`" + `
	jwt.RegisteredClaims
}

{{- if eq .Framework "gin"}}
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // TODO: Use environment variable
		})

		if err != nil || !token.Valid {
			logger.Error("Invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
{{- else if eq .Framework "echo"}}
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := extractToken(c.Request().Header.Get("Authorization"))
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization token"})
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("your-secret-key"), nil // TODO: Use environment variable
			})

			if err != nil || !token.Valid {
				logger.Error("Invalid JWT token", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)
			return next(c)
		}
	}
}
{{- else if eq .Framework "fiber"}}
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c.Get("Authorization"))
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization token"})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // TODO: Use environment variable
		})

		if err != nil || !token.Valid {
			logger.Error("Invalid JWT token", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		return c.Next()
	}
}
{{- end}}

func extractToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	
	return parts[1]
}
`

		if err := g.templateEngine.RenderToFile(authMiddlewareTemplate, filepath.Join(projectPath, "internal/middleware/auth.go"), data); err != nil {
			return err
		}
	}

	// Generate CORS middleware
	if cfg.Middleware.CORS {
		corsMiddlewareTemplate := `package middleware

import (
	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
)

{{- if eq .Framework "gin"}}
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
{{- else if eq .Framework "echo"}}
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	})
}
{{- else if eq .Framework "fiber"}}
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	})
}
{{- end}}
`

		if err := g.templateEngine.RenderToFile(corsMiddlewareTemplate, filepath.Join(projectPath, "internal/middleware/cors.go"), data); err != nil {
			return err
		}
	}

	// Generate request logging middleware
	if cfg.Middleware.Logging {
		loggingMiddlewareTemplate := `package middleware

import (
	"time"

	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
	"{{.ModulePath}}/pkg/logger"
	{{- if eq .Config.Logging "zap"}}
	"go.uber.org/zap"
	{{- end}}
)

{{- if eq .Framework "gin"}}
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		{{- if eq .Config.Logging "zap"}}
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
		)
		{{- else}}
		logger.Info("HTTP Request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
		)
		{{- end}}
		return ""
	})
}
{{- else if eq .Framework "echo"}}
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			{{- if eq .Config.Logging "zap"}}
			logger.Info("HTTP Request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("client_ip", c.RealIP()),
			)
			{{- else}}
			logger.Info("HTTP Request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency", time.Since(start),
				"client_ip", c.RealIP(),
			)
			{{- end}}
			
			return err
		}
	}
}
{{- else if eq .Framework "fiber"}}
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		err := c.Next()
		
		{{- if eq .Config.Logging "zap"}}
		logger.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.IP()),
		)
		{{- else}}
		logger.Info("HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", time.Since(start),
			"client_ip", c.IP(),
		)
		{{- end}}
		
		return err
	}
}
{{- end}}
`

		if err := g.templateEngine.RenderToFile(loggingMiddlewareTemplate, filepath.Join(projectPath, "internal/middleware/logging.go"), data); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateTestFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate main test file
	mainTestTemplate := `package main

import (
	"testing"
	"{{.ModulePath}}/internal/app"
)

func TestAppInitialization(t *testing.T) {
	app := app.New()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
}
`

	if err := g.templateEngine.RenderToFile(mainTestTemplate, filepath.Join(projectPath, "main_test.go"), data); err != nil {
		return err
	}

	// Generate test utilities
	testUtilsTemplate := `package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
)

// TestServer represents a test server for making HTTP requests
type TestServer struct {
	{{- if eq .Framework "gin"}}
	Router *gin.Engine
	{{- else if eq .Framework "echo"}}
	Echo *echo.Echo
	{{- else if eq .Framework "fiber"}}
	Fiber *fiber.App
	{{- end}}
}

// NewTestServer creates a new test server
func NewTestServer() *TestServer {
	{{- if eq .Framework "gin"}}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return &TestServer{Router: router}
	{{- else if eq .Framework "echo"}}
	e := echo.New()
	return &TestServer{Echo: e}
	{{- else if eq .Framework "fiber"}}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	return &TestServer{Fiber: app}
	{{- end}}
}

{{- if eq .Framework "gin"}}
// MakeRequest makes an HTTP request to the test server
func (ts *TestServer) MakeRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	
	w := httptest.NewRecorder()
	ts.Router.ServeHTTP(w, req)
	return w
}
{{- else if eq .Framework "echo"}}
// MakeRequest makes an HTTP request to the test server
func (ts *TestServer) MakeRequest(method, path string, body interface{}) (*http.Response, error) {
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	
	rec := httptest.NewRecorder()
	c := ts.Echo.NewContext(req, rec)
	ts.Echo.ServeHTTP(rec, req)
	return rec.Result(), nil
}
{{- else if eq .Framework "fiber"}}
// MakeRequest makes an HTTP request to the test server
func (ts *TestServer) MakeRequest(method, path string, body interface{}) (*http.Response, error) {
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	
	resp, err := ts.Fiber.Test(req)
	return resp, err
}
{{- end}}

// AssertJSON compares expected JSON with actual response body
func AssertJSON(t *testing.T, expected, actual string) {
	var expectedMap, actualMap map[string]interface{}
	
	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	
	if err := json.Unmarshal([]byte(actual), &actualMap); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}
	
	expectedJSON, _ := json.Marshal(expectedMap)
	actualJSON, _ := json.Marshal(actualMap)
	
	if !bytes.Equal(expectedJSON, actualJSON) {
		t.Errorf("JSON mismatch.\nExpected: %s\nActual: %s", expectedJSON, actualJSON)
	}
}
`

	if err := g.templateEngine.RenderToFile(testUtilsTemplate, filepath.Join(projectPath, "test/testutils/utils.go"), data); err != nil {
		return err
	}

	// Generate user handler tests (if models are generated)
	if cfg.ORM != config.ORMNone {
		userHandlerTestTemplate := `package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
	"{{.ModulePath}}/test/testutils"
)

func TestGetUsers(t *testing.T) {
	{{- if eq .Framework "gin"}}
	// Setup
	ts := testutils.NewTestServer()
	ts.Router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []interface{}{},
			"message": "Users retrieved successfully",
		})
	})

	// Test
	w := ts.MakeRequest("GET", "/users", nil)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}
	{{- else if eq .Framework "echo"}}
	// Setup
	ts := testutils.NewTestServer()
	ts.Echo.GET("/users", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"message": "Users retrieved successfully",
		})
	})

	// Test
	resp, err := ts.MakeRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	{{- else if eq .Framework "fiber"}}
	// Setup
	ts := testutils.NewTestServer()
	ts.Fiber.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []interface{}{},
			"message": "Users retrieved successfully",
		})
	})

	// Test
	resp, err := ts.MakeRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	{{- end}}
}
`

		if err := g.templateEngine.RenderToFile(userHandlerTestTemplate, filepath.Join(projectPath, "test/handlers/user_test.go"), data); err != nil {
			return err
		}
	}

	return nil
}

// generateDocsPackage generates Swagger docs package
func (g *Generator) generateDocsPackage(cfg *config.ProjectConfig, projectPath string) error {
	if !cfg.Features.Swagger {
		return nil // Skip if Swagger is not enabled
	}

	data := templates.NewTemplateData(cfg)

	docsTemplate := `// Package docs GENERATED BY SWAG; DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = ` + "`" + `{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "description": "{{.ProjectName}} API documentation",
        "title": "{{.ProjectName}} API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {
        "/health": {
            "get": {
                "description": "Check if the service is running",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HealthCheckResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Get a paginated list of users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Items per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create a new user with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get a user by their ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.APIResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "$ref": "#/definitions/models.APIError"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "models.CreateUserRequest": {
            "type": "object",
            "required": [
                "email",
                "name"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "models.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "services": {},
                "status": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "models.UserResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}` + "`" + `

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "{{.ProjectName}} API",
	Description:      "{{.ProjectName}} API documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
`

	if err := g.templateEngine.RenderToFile(docsTemplate, filepath.Join(projectPath, "docs/docs.go"), data); err != nil {
		return err
	}

	return nil
}
