package generator

import (
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// generateConfigFiles generates configuration files
func (g *Generator) generateConfigFiles(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate pkg/config/config.go
	configTemplate := `package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      ` + "`yaml:\"app\" json:\"app\"`" + `
	Database DatabaseConfig ` + "`yaml:\"database\" json:\"database\"`" + `
	{{- if eq .Config.Auth "jwt"}}
	JWT      JWTConfig      ` + "`yaml:\"jwt\" json:\"jwt\"`" + `
	{{- end}}
	Log      LogConfig      ` + "`yaml:\"log\" json:\"log\"`" + `
	{{- if .Config.Features.Caching}}
	Redis    RedisConfig    ` + "`yaml:\"redis\" json:\"redis\"`" + `
	{{- end}}
}

type AppConfig struct {
	Name  string ` + "`yaml:\"name\" json:\"name\"`" + `
	Env   string ` + "`yaml:\"env\" json:\"env\"`" + `
	Port  string ` + "`yaml:\"port\" json:\"port\"`" + `
	Debug bool   ` + "`yaml:\"debug\" json:\"debug\"`" + `
}

type DatabaseConfig struct {
	{{- if eq .Database "postgresql" "mysql"}}
	Host     string ` + "`yaml:\"host\" json:\"host\"`" + `
	Port     string ` + "`yaml:\"port\" json:\"port\"`" + `
	User     string ` + "`yaml:\"user\" json:\"user\"`" + `
	Password string ` + "`yaml:\"password\" json:\"password\"`" + `
	Name     string ` + "`yaml:\"name\" json:\"name\"`" + `
	{{- if eq .Database "postgresql"}}
	SSLMode  string ` + "`yaml:\"sslmode\" json:\"sslmode\"`" + `
	{{- end}}
	{{- else if eq .Database "sqlite"}}
	Path     string ` + "`yaml:\"path\" json:\"path\"`" + `
	{{- end}}
}

{{- if eq .Config.Auth "jwt"}}
type JWTConfig struct {
	Secret string ` + "`yaml:\"secret\" json:\"secret\"`" + `
	Expiry string ` + "`yaml:\"expiry\" json:\"expiry\"`" + `
}
{{- end}}

type LogConfig struct {
	Level  string ` + "`yaml:\"level\" json:\"level\"`" + `
	Format string ` + "`yaml:\"format\" json:\"format\"`" + `
}

{{- if .Config.Features.Caching}}
type RedisConfig struct {
	Host     string ` + "`yaml:\"host\" json:\"host\"`" + `
	Port     string ` + "`yaml:\"port\" json:\"port\"`" + `
	Password string ` + "`yaml:\"password\" json:\"password\"`" + `
	DB       int    ` + "`yaml:\"db\" json:\"db\"`" + `
}
{{- end}}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	viper.AutomaticEnv()

	config := &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "{{.ProjectName}}"),
			Env:   getEnv("APP_ENV", "development"),
			Port:  getEnv("APP_PORT", "8080"),
			Debug: getEnv("APP_DEBUG", "true") == "true",
		},
		{{- if ne .Database ""}}
		Database: DatabaseConfig{
			{{- if eq .Database "postgresql" "mysql"}}
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "{{if eq .Database "postgresql"}}5432{{else}}3306{{end}}"),
			User:     getEnv("DB_USER", "{{if eq .Database "postgresql"}}postgres{{else}}root{{end}}"),
			Password: getEnv("DB_PASSWORD", "{{if eq .Database "postgresql"}}postgres{{else}}root{{end}}"),
			Name:     getEnv("DB_NAME", "{{.ProjectName}}_db"),
			{{- if eq .Database "postgresql"}}
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			{{- end}}
			{{- else if eq .Database "sqlite"}}
			Path:     getEnv("DB_PATH", "./{{.ProjectName}}.db"),
			{{- end}}
		},
		{{- end}}
		{{- if eq .Config.Auth "jwt"}}
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		{{- end}}
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "{{if eq .Config.Logging "zap"}}json{{else if eq .Config.Logging "charm"}}text{{else}}text{{end}}"),
		},
		{{- if .Config.Features.Caching}}
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		{{- end}}
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`

	if err := g.templateEngine.RenderToFile(configTemplate, filepath.Join(projectPath, "pkg/config/config.go"), data); err != nil {
		return err
	}

	return nil
}

// generateDatabaseFiles generates database-related files
func (g *Generator) generateDatabaseFiles(cfg *config.ProjectConfig, projectPath string) error {
	if cfg.ORM == config.ORMNone {
		return nil
	}

	data := templates.NewTemplateData(cfg)

	if cfg.ORM == config.ORMGorm {
		// Generate GORM database setup
		dbTemplate := `package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	{{- if eq .Database "postgresql"}}
	"gorm.io/driver/postgres"
	{{- else if eq .Database "mysql"}}
	"gorm.io/driver/mysql"
	{{- else if eq .Database "sqlite"}}
	"gorm.io/driver/sqlite"
	{{- end}}
	"{{.ModulePath}}/pkg/config"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	var err error
	
	{{- if eq .Database "postgresql"}}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.Port, cfg.Database.SSLMode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	{{- else if eq .Database "mysql"}}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.Name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	{{- else if eq .Database "sqlite"}}
	DB, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	{{- end}}
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	
	// Auto-migrate tables here
	// err = DB.AutoMigrate(&models.User{})
	// if err != nil {
	//     return fmt.Errorf("failed to migrate database: %w", err)
	// }

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
`

		if err := g.templateEngine.RenderToFile(dbTemplate, filepath.Join(projectPath, "pkg/database/database.go"), data); err != nil {
			return err
		}
	}

	return nil
}

// generateRoutes generates API routes
func (g *Generator) generateRoutes(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	routesTemplate := `package routes

import (
	{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	"{{.ModulePath}}/internal/handlers"
	{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	"{{.ModulePath}}/internal/handlers"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/internal/handlers"
	{{- end}}
)

{{- if eq .Framework "gin"}}
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		{{- if .Config.Features.HealthCheck}}
		api.GET("/health", handlers.HealthCheck)
		{{- end}}
		
		// Example routes
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)
		api.POST("/users", handlers.CreateUser)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
		
		{{- if ne .Config.Auth "none"}}
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			{{- if eq .Config.Auth "jwt"}}
			auth.POST("/refresh", handlers.RefreshToken)
			{{- end}}
		}
		{{- end}}
	}
}
{{- else if eq .Framework "echo"}}
func SetupRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")
	
	{{- if .Config.Features.HealthCheck}}
	api.GET("/health", handlers.HealthCheck)
	{{- end}}
	
	// Example routes
	api.GET("/users", handlers.GetUsers)
	api.GET("/users/:id", handlers.GetUser)
	api.POST("/users", handlers.CreateUser)
	api.PUT("/users/:id", handlers.UpdateUser)
	api.DELETE("/users/:id", handlers.DeleteUser)
	
	{{- if ne .Config.Auth "none"}}
	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/login", handlers.Login)
	auth.POST("/register", handlers.Register)
	{{- if eq .Config.Auth "jwt"}}
	auth.POST("/refresh", handlers.RefreshToken)
	{{- end}}
	{{- end}}
}
{{- else if eq .Framework "fiber"}}
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	
	{{- if .Config.Features.HealthCheck}}
	api.Get("/health", handlers.HealthCheck)
	{{- end}}
	
	// Example routes
	api.Get("/users", handlers.GetUsers)
	api.Get("/users/:id", handlers.GetUser)
	api.Post("/users", handlers.CreateUser)
	api.Put("/users/:id", handlers.UpdateUser)
	api.Delete("/users/:id", handlers.DeleteUser)
	
	{{- if ne .Config.Auth "none"}}
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)
	{{- if eq .Config.Auth "jwt"}}
	auth.Post("/refresh", handlers.RefreshToken)
	{{- end}}
	{{- end}}
}
{{- end}}
`

	if err := g.templateEngine.RenderToFile(routesTemplate, filepath.Join(projectPath, "api/routes/routes.go"), data); err != nil {
		return err
	}

	return nil
}

// generateHandlers generates basic handlers
func (g *Generator) generateHandlers(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	handlersTemplate := `package handlers

import (
	"strconv"

	{{- if eq .Framework "gin"}}
	"net/http"
	"github.com/gin-gonic/gin"
	{{- else if eq .Framework "echo"}}
	"net/http"
	"github.com/labstack/echo/v4"
	{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
	{{- end}}
)

{{- if .Config.Features.HealthCheck}}
// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
{{- if eq .Framework "gin"}}
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "{{.ProjectName}} is running",
	})
}
{{- else if eq .Framework "echo"}}
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"message": "{{.ProjectName}} is running",
	})
}
{{- else if eq .Framework "fiber"}}
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"message": "{{.ProjectName}} is running",
	})
}
{{- end}}
{{- end}}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /users [get]
{{- if eq .Framework "gin"}}
func GetUsers(c *gin.Context) {
	// TODO: Implement get users logic
	c.JSON(http.StatusOK, gin.H{
		"users": []map[string]interface{}{
			{"id": 1, "name": "John Doe", "email": "john@example.com"},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		},
	})
}
{{- else if eq .Framework "echo"}}
func GetUsers(c echo.Context) error {
	// TODO: Implement get users logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "John Doe", "email": "john@example.com"},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		},
	})
}
{{- else if eq .Framework "fiber"}}
func GetUsers(c *fiber.Ctx) error {
	// TODO: Implement get users logic
	return c.JSON(fiber.Map{
		"users": []map[string]interface{}{
			{"id": 1, "name": "John Doe", "email": "john@example.com"},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		},
	})
}
{{- end}}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [get]
{{- if eq .Framework "gin"}}
func GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement get user by ID logic
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}
{{- else if eq .Framework "echo"}}
func GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// TODO: Implement get user by ID logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}
{{- else if eq .Framework "fiber"}}
func GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// TODO: Implement get user by ID logic
	return c.JSON(fiber.Map{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}
{{- end}}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /users [post]
{{- if eq .Framework "gin"}}
func CreateUser(c *gin.Context) {
	// TODO: Implement create user logic
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    3,
			"name":  "New User",
			"email": "newuser@example.com",
		},
	})
}
{{- else if eq .Framework "echo"}}
func CreateUser(c echo.Context) error {
	// TODO: Implement create user logic
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":    3,
			"name":  "New User",
			"email": "newuser@example.com",
		},
	})
}
{{- else if eq .Framework "fiber"}}
func CreateUser(c *fiber.Ctx) error {
	// TODO: Implement create user logic
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":    3,
			"name":  "New User",
			"email": "newuser@example.com",
		},
	})
}
{{- end}}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [put]
{{- if eq .Framework "gin"}}
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement update user logic
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":    id,
			"name":  "Updated User",
			"email": "updated@example.com",
		},
	})
}
{{- else if eq .Framework "echo"}}
func UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// TODO: Implement update user logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"user": map[string]interface{}{
			"id":    id,
			"name":  "Updated User",
			"email": "updated@example.com",
		},
	})
}
{{- else if eq .Framework "fiber"}}
func UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// TODO: Implement update user logic
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"user": fiber.Map{
			"id":    id,
			"name":  "Updated User",
			"email": "updated@example.com",
		},
	})
}
{{- end}}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [delete]
{{- if eq .Framework "gin"}}
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement delete user logic
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
}
{{- else if eq .Framework "echo"}}
func DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// TODO: Implement delete user logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
		"id":      id,
	})
}
{{- else if eq .Framework "fiber"}}
func DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// TODO: Implement delete user logic
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
		"id":      id,
	})
}
{{- end}}

{{- if ne .Config.Auth "none"}}
// Login godoc
// @Summary User login
// @Description Authenticate user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
{{- if eq .Framework "gin"}}
func Login(c *gin.Context) {
	// TODO: Implement login logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		{{- if eq .Config.Auth "jwt"}}
		"token": "your-jwt-token-here",
		{{- end}}
	})
}
{{- else if eq .Framework "echo"}}
func Login(c echo.Context) error {
	// TODO: Implement login logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		{{- if eq .Config.Auth "jwt"}}
		"token": "your-jwt-token-here",
		{{- end}}
	})
}
{{- else if eq .Framework "fiber"}}
func Login(c *fiber.Ctx) error {
	// TODO: Implement login logic
	return c.JSON(fiber.Map{
		"message": "Login successful",
		{{- if eq .Config.Auth "jwt"}}
		"token": "your-jwt-token-here",
		{{- end}}
	})
}
{{- end}}

// Register godoc
// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /auth/register [post]
{{- if eq .Framework "gin"}}
func Register(c *gin.Context) {
	// TODO: Implement registration logic
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
{{- else if eq .Framework "echo"}}
func Register(c echo.Context) error {
	// TODO: Implement registration logic
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
	})
}
{{- else if eq .Framework "fiber"}}
func Register(c *fiber.Ctx) error {
	// TODO: Implement registration logic
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
{{- end}}

{{- if eq .Config.Auth "jwt"}}
// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh an expired JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/refresh [post]
{{- if eq .Framework "gin"}}
func RefreshToken(c *gin.Context) {
	// TODO: Implement token refresh logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token": "new-jwt-token-here",
	})
}
{{- else if eq .Framework "echo"}}
func RefreshToken(c echo.Context) error {
	// TODO: Implement token refresh logic
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Token refreshed successfully",
		"token": "new-jwt-token-here",
	})
}
{{- else if eq .Framework "fiber"}}
func RefreshToken(c *fiber.Ctx) error {
	// TODO: Implement token refresh logic
	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"token": "new-jwt-token-here",
	})
}
{{- end}}
{{- end}}
{{- end}}
`

	if err := g.templateEngine.RenderToFile(handlersTemplate, filepath.Join(projectPath, "internal/handlers/handlers.go"), data); err != nil {
		return err
	}

	return nil
}

// generateLogger generates logger configuration
func (g *Generator) generateLogger(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	loggerTemplate := `package logger

import (
	{{- if eq .Config.Logging "zap"}}
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	{{- else if eq .Config.Logging "logrus"}}
	"github.com/sirupsen/logrus"
	{{- else if eq .Config.Logging "charm"}}
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/lipgloss"
	"os"
	{{- else}}
	"log"
	"os"
	{{- end}}
)

{{- if eq .Config.Logging "zap"}}
var logger *zap.Logger

func Init(level, format string) {
	config := zap.NewProductionConfig()
	
	if format == "console" {
		config = zap.NewDevelopmentConfig()
	}
	
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}
	
	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

{{- else if eq .Config.Logging "logrus"}}
var logger *logrus.Logger

func Init(level, format string) {
	logger = logrus.New()
	
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}
	
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	}
}

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}

{{- else if eq .Config.Logging "charm"}}
var logger *log.Logger

func Init(level, format string) {
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      "15:04:05",
		Prefix:          "{{.ProjectName}} 🪵",
	})

	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	}

	// Customize styles for a beautiful look
	styles := log.DefaultStyles()
	
	// Beautiful gradient colors for different levels
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("63")).   // Purple
		Foreground(lipgloss.Color("255"))   // White
	
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("86")).   // Cyan
		Foreground(lipgloss.Color("0"))     // Black
	
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("214")).  // Orange
		Foreground(lipgloss.Color("0"))     // Black
	
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("196")).  // Red
		Foreground(lipgloss.Color("255"))   // White
	
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("124")).  // Dark Red
		Foreground(lipgloss.Color("255"))   // White
		Bold(true)

	// Beautiful key styles
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Keys["error"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Keys["status"] = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	styles.Keys["method"] = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	styles.Keys["path"] = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	styles.Keys["latency"] = lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
	styles.Keys["user"] = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	// Beautiful value styles
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	styles.Values["error"] = lipgloss.NewStyle().Bold(true)
	
	// Set timestamp and caller styles
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	styles.Caller = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styles.Prefix = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	
	logger.SetStyles(styles)
	
	// Set formatter based on format preference
	if format == "json" {
		logger.SetFormatter(log.JSONFormatter)
	} else if format == "logfmt" {
		logger.SetFormatter(log.LogfmtFormatter)
	}
	// Default is TextFormatter with beautiful colors
}

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}

func Print(msg string, args ...interface{}) {
	logger.Print(msg, args...)
}

// WithFields creates a new logger with additional fields
func WithFields(fields map[string]interface{}) *log.Logger {
	keyValues := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		keyValues = append(keyValues, k, v)
	}
	return logger.With(keyValues...)
}

// Helper method for structured logging
func With(key string, value interface{}) *log.Logger {
	return logger.With(key, value)
}

{{- else}}
var logger *log.Logger

func Init(level, format string) {
	logger = log.New(os.Stdout, "[{{.ProjectName}}] ", log.LstdFlags)
}

func Info(msg string, args ...interface{}) {
	logger.Printf("[INFO] "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Printf("[ERROR] "+msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Printf("[DEBUG] "+msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Printf("[WARN] "+msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatalf("[FATAL] "+msg, args...)
}
{{- end}}
`

	if err := g.templateEngine.RenderToFile(loggerTemplate, filepath.Join(projectPath, "pkg/logger/logger.go"), data); err != nil {
		return err
	}

	return nil
}

// Update the main generateFeatureFiles method to call these new generators
func (g *Generator) generateFeatureFiles(cfg *config.ProjectConfig, projectPath string) error {
	// Generate routes
	if err := g.generateRoutes(cfg, projectPath); err != nil {
		return err
	}

	// Generate handlers
	if err := g.generateHandlers(cfg, projectPath); err != nil {
		return err
	}

	// Generate logger
	if err := g.generateLogger(cfg, projectPath); err != nil {
		return err
	}

	return nil
}

// generateBeautifulStartup generates beautiful startup package with colors and ASCII art
func (g *Generator) generateBeautifulStartup(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Create main startup.go file
	startupMainTemplate := `package startup

import (
	"fmt"
	"os"
	"time"
)

// Modern gradient color palette
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	
	// Modern gradient colors
	Purple    = "\033[38;5;129m"  // Rich purple
	Pink      = "\033[38;5;205m"  // Bright pink
	Cyan      = "\033[38;5;87m"   // Bright cyan
	Blue      = "\033[38;5;75m"   // Sky blue
	Green     = "\033[38;5;120m"  // Mint green
	Yellow    = "\033[38;5;227m"  // Bright yellow
	Orange    = "\033[38;5;215m"  // Peach orange
	Red       = "\033[38;5;196m"  // Bright red
	Indigo    = "\033[38;5;63m"   // Indigo
	
	// Neutral colors
	White     = "\033[37m"        // Normal white
	Gray      = "\033[38;5;245m"  // Medium gray
	LightGray = "\033[38;5;250m"  // Light gray
	DarkGray  = "\033[38;5;240m"  // Dark gray
	
	// Special effects
	Rainbow1  = "\033[38;5;196m"  // Red
	Rainbow2  = "\033[38;5;202m"  // Orange
	Rainbow3  = "\033[38;5;226m"  // Yellow  
	Rainbow4  = "\033[38;5;46m"   // Green
	Rainbow5  = "\033[38;5;21m"   // Blue
	Rainbow6  = "\033[38;5;129m"  // Purple
)

var projectInfo = ProjectInfo{
	Framework:    "{{.Framework | title}}",
	Database:     "{{.Database | title}}",
	Architecture: "{{.Config.Architecture | title}}",
}

type ProjectInfo struct {
	Framework    string
	Database     string
	Architecture string
}

// ShowWelcome displays a beautiful welcome message
func ShowWelcome(appName, version, port string) {
	if os.Getenv("SILENT") == "true" {
		return
	}
	
	clearScreen()
	showHeader(appName)
	showInfo(version, port)
	showFooter()
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func showHeader(appName string) {
	
	// Modern logo with gradient effect
	fmt.Printf("\n%s", Bold)
	fmt.Printf("%s    ▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓    %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s  ▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓  %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s ▓▓▓%s▓▓▓%s░░░%s░░░%s░░░%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓ %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s▓▓▓%s▓▓▓%s░%s(●)%s░%s░%s(●)%s░%s░%s▓▓▓%s▓▓▓%s▓▓▓%s\n", Blue, Blue, Blue, White, Blue, White, Blue, White, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s▓▓▓%s▓▓▓%s░░%s ~~~~ %s░░%s░░░%s▓▓▓%s▓▓▓%s▓▓▓ %s\n", Blue, Blue, Blue, White, White, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s ▓▓▓%s▓▓▓%s░░░%s░░░%s░░░ %s▓▓▓%s▓▓▓%s▓▓▓%s\n", Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s  ▓▓▓%s▓▓▓%s▓▓▓%s▓▓▓ %s▓▓▓%s▓▓▓%s▓▓▓  %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s    ▓▓▓%s▓▓▓%s▓▓▓ %s▓▓▓%s▓▓▓%s▓▓▓    %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	fmt.Printf("%s    ▓▓▓%s▓▓▓%s▓▓▓ %s▓▓▓%s▓▓▓%s▓▓▓    %s\n", Blue, Blue, Blue, Blue, Blue, Blue, Reset)
	
	// Beautiful project name with sparkles
	fmt.Printf("\n%s", Bold)
	fmt.Printf("%s%s🌟 %s %s🌟%s\n", Pink, Bold, appName, Pink, Reset)
	
	// Subtitle with elegant styling
	fmt.Printf("%s%s✨ %s🔧 Gool CLI Generator%s ✨%s\n", White, Bold, Cyan, White, Reset)
	fmt.Printf("%s%s✨ %s⚡ Fast, Modern, Beautiful API%s ✨%s\n", White, Bold, Blue, White, Reset)
	fmt.Printf("%s%s✨ %s🛠️  Make Command + Ctrl + C%s ✨%s\n", White, Reset, White, White, Reset)
}

func showInfo(version, port string) {
	fmt.Printf("%s──────────────────────────────────────────────────────────%s\n", Blue, Bold)
	fmt.Printf("\n%s%sProject Stack%s\n", Blue, Bold, Reset)
	fmt.Printf("%s%s %s🚀 Framework     %-25s  %s%s\n", DarkGray, Bold, Cyan,  fmt.Sprintf("%s%s%s", White, projectInfo.Framework, Reset), DarkGray, Reset)
	fmt.Printf("%s%s %s💾 Database      %-25s   %s%s\n", DarkGray, Bold, Pink,  fmt.Sprintf("%s%s%s", White, projectInfo.Database, Reset), DarkGray, Reset)
	fmt.Printf("%s%s %s🏗️  Architecture  %-25s   %s%s\n", DarkGray, Bold, Green,  fmt.Sprintf("%s%s%s", White, projectInfo.Architecture, Reset), DarkGray, Reset)
	fmt.Printf("%s%s %s📦 Version       %-25s   %s%s\n", DarkGray, Bold, Yellow,  fmt.Sprintf("%s%s%s", White, version, Reset), DarkGray, Reset)
	
	fmt.Printf("%s──────────────────────────────────────────────────────────%s\n", Blue, Bold)
	fmt.Printf("\n%s%sServer Endpoints%s %s(Click to open)%s\n", Blue, Bold, Bold, Gray, Bold)
	
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	healthURL := fmt.Sprintf("http://localhost:%s/api/v1/health", port)
	docsURL := fmt.Sprintf("http://localhost:%s/swagger/index.html", port)
	apiURL := fmt.Sprintf("http://localhost:%s/api/v1", port)
	
	fmt.Printf("%s%s %s🏠 Server       %s%s%-35s%s      %s%s\n", 
		DarkGray, Bold, Indigo, White, Reset, serverURL, Bold, DarkGray, Bold)
	fmt.Printf("%s%s %s💖 Health       %s%s%-35s%s      %s%s\n", 
		DarkGray, Bold, Pink, White, Reset, healthURL, Bold, DarkGray, Bold)
	fmt.Printf("%s%s %s📚 API Docs     %s%s%-35s%s %s%s\n", 
		DarkGray, Bold, Green, White, Reset, docsURL, Bold, DarkGray, Bold)
	fmt.Printf("%s%s %s🔗 API Base     %s%s%-35s%s      %s%s\n", 
		DarkGray, Bold, Orange, White, Reset, apiURL, Bold, DarkGray, Bold)
	fmt.Printf("%s%s %s⏰ Started      %s%s%-25s%s                %s%s\n", 
		DarkGray, Bold, Cyan, White, Reset, 
		time.Now().Format("2006-01-02 15:04:05"), Bold, DarkGray, Bold)
}

func showFooter() {
	// Beautiful animated startup sequence
	fmt.Printf("\n%s%sInitializing Server%s\n", Green, Bold, Bold)
	fmt.Printf("%s", White)
	
	loadingChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	colors := []string{Rainbow1, Rainbow2, Rainbow3, Rainbow4, Rainbow5, Rainbow6}
	
	for i := 0; i < 20; i++ {
		fmt.Printf("\r%s%s%s%s Loading", White, colors[i%len(colors)], loadingChars[i%len(loadingChars)], Reset)
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Printf("\r%s%s✅ Server Ready!%s", White, Green, Bold)
}

func centerText(text string, width int) {
	textLen := len(stripANSI(text))
	padding := (width - textLen) / 2
	fmt.Printf("%*s%s%*s", padding, "", text, width-textLen-padding, "")
}

func stripANSI(text string) string {
	// Enhanced ANSI stripper for better length calculation
	result := ""
	inEscape := false
	inOSC := false
	
	for i, r := range text {
		if r == '\033' {
			if i+1 < len(text) && text[i+1] == ']' {
				inOSC = true
			} else {
				inEscape = true
			}
			continue
		}
		
		if inOSC {
			if r == '\a' || (r == '\\' && i > 0 && text[i-1] == '\033') {
				inOSC = false
			}
			continue
		}
		
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		
		result += string(r)
	}
	return result
}

// ShowShutdown displays a beautiful shutdown message  
func ShowShutdown(appName string) {
	fmt.Printf("\n%s%s🌙 Graceful Shutdown%s\n", Purple, Bold, Reset)
	fmt.Printf("%s╭─────────────────────────────────────────────────╮%s\n", DarkGray, Reset)
	fmt.Printf("%s│%s %s📦 Project      %s│%s %-25s %s│%s\n", DarkGray, Reset, Cyan, Reset, White, fmt.Sprintf("%s%s%s", Pink, appName, Reset), DarkGray, Reset)
	fmt.Printf("%s│%s %s🔴 Status       %s│%s %-25s %s│%s\n", DarkGray, Reset, Red, Reset, White, fmt.Sprintf("%s%s%s", Gray, "Stopped", Reset), DarkGray, Reset)
	fmt.Printf("%s│%s %s⏰ Time         %s│%s %-25s %s│%s\n", DarkGray, Reset, Yellow, Reset, White, time.Now().Format("15:04:05"), DarkGray, Reset)
	fmt.Printf("%s╰─────────────────────────────────────────────────╯%s\n", DarkGray, Reset)
	
	// Cute goodbye message
	fmt.Printf("\n%s", Pink)
	fmt.Print("                    ╭─────────────────────╮\n")
	fmt.Printf("                    │%s  %s✨ See you soon! ✨%s  %s│\n", Reset, Cyan, Reset, Pink)
	fmt.Printf("                    │%s                     %s│\n", Reset, Pink)
	fmt.Printf("                    │%s     %s(=^･ω･^=)%s      %s│\n", Reset, Yellow, Reset, Pink)
	fmt.Print("                    ╰─────────────────────╯\n")
	fmt.Printf("%s\n", Reset)
	
	fmt.Printf("%s%s💝 Thanks for using Gool! Happy coding! 💝%s\n\n", Green, Bold, Reset)
}
`

	if err := g.templateEngine.RenderToFile(startupMainTemplate, filepath.Join(projectPath, "pkg/startup/startup.go"), data); err != nil {
		return err
	}

	return nil
}
