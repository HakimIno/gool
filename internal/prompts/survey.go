package prompts

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gool-cli/gool/internal/config"
)

// CollectProjectConfig prompts the user for project configuration
func CollectProjectConfig() (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{}

	// Create stylish header
	printHeader()

	// Project name
	projectNamePrompt := &survey.Input{
		Message: "🚀 What is your project name?",
		Default: "my-go-app",
		Help:    "Choose a descriptive name for your Go project",
	}
	if err := survey.AskOne(projectNamePrompt, &cfg.ProjectName); err != nil {
		return nil, err
	}

	// Module path
	modulePathPrompt := &survey.Input{
		Message: "📦 What is your module path?",
		Default: fmt.Sprintf("github.com/username/%s", cfg.ProjectName),
		Help:    "This will be used as the Go module name (e.g., github.com/username/project)",
	}
	if err := survey.AskOne(modulePathPrompt, &cfg.ModulePath); err != nil {
		return nil, err
	}

	// Framework selection
	frameworkPrompt := &survey.Select{
		Message: "🌐 Choose your web framework:",
		Options: []string{
			fmt.Sprintf("🔥 %s - Fast and minimalist", config.FrameworkGin),
			fmt.Sprintf("⚡ %s - High performance and extensible", config.FrameworkEcho),
			fmt.Sprintf("🚀 %s - Express inspired framework", config.FrameworkFiber),
			fmt.Sprintf("🎯 %s - Full-stack web framework", config.FrameworkRevel),
		},
		Default: fmt.Sprintf("🔥 %s - Fast and minimalist", config.FrameworkGin),
		Help:    "Select the web framework that best fits your project needs",
	}
	var selectedFramework string
	if err := survey.AskOne(frameworkPrompt, &selectedFramework); err != nil {
		return nil, err
	}
	cfg.Framework = extractFrameworkName(selectedFramework)

	// ORM selection
	ormPrompt := &survey.Select{
		Message: "🗄️  Choose your database access layer:",
		Options: []string{
			fmt.Sprintf("🏗️  %s - Feature-rich ORM", config.ORMGorm),
			fmt.Sprintf("⚙️  %s - Extensions on database/sql", config.ORMSqlx),
			fmt.Sprintf("🔧 %s - Direct SQL queries", config.ORMRaw),
			fmt.Sprintf("🚫 %s - No database", config.ORMNone),
		},
		Default: fmt.Sprintf("🏗️  %s - Feature-rich ORM", config.ORMGorm),
		Help:    "Choose how you want to interact with your database",
	}
	var selectedORM string
	if err := survey.AskOne(ormPrompt, &selectedORM); err != nil {
		return nil, err
	}
	cfg.ORM = extractORMName(selectedORM)

	// Database selection (only if ORM is not none)
	if cfg.ORM != config.ORMNone {
		dbPrompt := &survey.Select{
			Message: "💾 Choose your database:",
			Options: []string{
				fmt.Sprintf("🐘 %s - Advanced open source database", config.DBPostgreSQL),
				fmt.Sprintf("🐬 %s - Popular relational database", config.DBMySQL),
				fmt.Sprintf("📁 %s - Lightweight file-based database", config.DBSQLite),
				fmt.Sprintf("🍃 %s - NoSQL document database", config.DBMongoDB),
				fmt.Sprintf("⚡ %s - In-memory data store", config.DBRedis),
				fmt.Sprintf("💾 %s - In-memory storage", config.DBMemory),
			},
			Default: fmt.Sprintf("🐘 %s - Advanced open source database", config.DBPostgreSQL),
			Help:    "Select the database that matches your project requirements",
		}
		var selectedDB string
		if err := survey.AskOne(dbPrompt, &selectedDB); err != nil {
			return nil, err
		}
		cfg.Database = extractDBName(selectedDB)
	}

	// Architecture selection
	archPrompt := &survey.Select{
		Message: "🏛️  Choose your project architecture:",
		Options: []string{
			fmt.Sprintf("🎯 %s - Straightforward and easy to understand", config.ArchSimple),
			fmt.Sprintf("🧹 %s - Clean Architecture with clear separation", config.ArchClean),
			fmt.Sprintf("⬡ %s - Ports and Adapters pattern", config.ArchHexagonal),
			fmt.Sprintf("🎨 %s - Model-View-Controller pattern", config.ArchMVC),
			fmt.Sprintf("🔧 %s - Define your own structure", config.ArchCustom),
		},
		Default: fmt.Sprintf("🎯 %s - Straightforward and easy to understand", config.ArchSimple),
		Help:    "Choose the architectural pattern that best fits your team and project",
	}
	var selectedArch string
	if err := survey.AskOne(archPrompt, &selectedArch); err != nil {
		return nil, err
	}
	cfg.Architecture = extractArchName(selectedArch)

	// Configuration format
	configPrompt := &survey.Select{
		Message: "⚙️  Choose your configuration format:",
		Options: []string{
			fmt.Sprintf("📄 %s - Human-readable data serialization", config.ConfigYAML),
			fmt.Sprintf("📋 %s - JavaScript Object Notation", config.ConfigJSON),
			fmt.Sprintf("📜 %s - Tom's Obvious, Minimal Language", config.ConfigTOML),
		},
		Default: fmt.Sprintf("📄 %s - Human-readable data serialization", config.ConfigYAML),
		Help:    "Select the configuration file format you prefer",
	}
	var selectedConfig string
	if err := survey.AskOne(configPrompt, &selectedConfig); err != nil {
		return nil, err
	}
	cfg.Config = extractConfigName(selectedConfig)

	// Authentication
	authPrompt := &survey.Select{
		Message: "🔐 Choose your authentication method:",
		Options: []string{
			fmt.Sprintf("🎫 %s - JSON Web Tokens", config.AuthJWT),
			fmt.Sprintf("🌐 %s - OAuth2 protocol", config.AuthOAuth2),
			fmt.Sprintf("🔑 %s - Username and password", config.AuthBasic),
			fmt.Sprintf("🚫 %s - No authentication", config.AuthNone),
		},
		Default: fmt.Sprintf("🎫 %s - JSON Web Tokens", config.AuthJWT),
		Help:    "Choose how users will authenticate with your application",
	}
	var selectedAuth string
	if err := survey.AskOne(authPrompt, &selectedAuth); err != nil {
		return nil, err
	}
	cfg.Auth = extractAuthName(selectedAuth)

	// Logging
	loggingPrompt := &survey.Select{
		Message: "📝 Choose your logging library:",
		Options: []string{
			fmt.Sprintf("✨ %s - Beautiful, colorful logging with Lip Gloss", config.LogCharm),
			fmt.Sprintf("⚡ %s - Blazing fast, structured logging", config.LogZap),
			fmt.Sprintf("📋 %s - Structured logging with levels", config.LogLogrus),
			fmt.Sprintf("📄 %s - Standard Go logging", config.LogStandard),
		},
		Default: fmt.Sprintf("✨ %s - Beautiful, colorful logging with Lip Gloss", config.LogCharm),
		Help:    "Select the logging solution that best fits your needs",
	}
	var selectedLogging string
	if err := survey.AskOne(loggingPrompt, &selectedLogging); err != nil {
		return nil, err
	}
	cfg.Logging = extractLoggingName(selectedLogging)

	// Features selection
	color.Cyan("\n✨ Additional Features")
	featuresPrompt := &survey.MultiSelect{
		Message: "🎁 Select additional features you want to include:",
		Options: []string{
			"🧪 Testing templates and examples",
			"🐳 Docker support (Dockerfile + docker-compose)",
			"🔌 WebSocket support for real-time features",
			"⚡ Caching support (Redis/in-memory)",
			"❤️  Health check endpoint",
			"📚 API documentation (Swagger/OpenAPI)",
			"🖼️  Static files serving",
			"🌍 Internationalization (i18n) support",
			"📊 Prometheus metrics",
			"☁️  Cloud deployment configuration",
		},
		Help: "Select all the features you want to include in your project",
	}
	var selectedFeatures []string
	if err := survey.AskOne(featuresPrompt, &selectedFeatures); err != nil {
		return nil, err
	}

	// Set features based on selection
	for _, feature := range selectedFeatures {
		switch {
		case strings.Contains(feature, "Testing"):
			cfg.Testing = true
		case strings.Contains(feature, "Docker"):
			cfg.Docker = true
		case strings.Contains(feature, "WebSocket"):
			cfg.Features.WebSocket = true
		case strings.Contains(feature, "Caching"):
			cfg.Features.Caching = true
		case strings.Contains(feature, "Health"):
			cfg.Features.HealthCheck = true
		case strings.Contains(feature, "API documentation"):
			cfg.Features.Swagger = true
		case strings.Contains(feature, "Static"):
			cfg.Features.StaticFiles = true
		case strings.Contains(feature, "Internationalization"):
			cfg.Features.I18n = true
		case strings.Contains(feature, "Prometheus"):
			cfg.Features.Metrics = true
		case strings.Contains(feature, "Cloud"):
			cfg.Features.CloudConfig = true
		}
	}

	// Middleware selection
	color.Cyan("\n🔧 Middleware Components")
	middlewarePrompt := &survey.MultiSelect{
		Message: "🛡️  Select middleware components:",
		Options: []string{
			"🌐 CORS (Cross-Origin Resource Sharing)",
			"⏱️  Rate limiting protection",
			"📝 Request logging middleware",
			"🔐 Authentication middleware",
			"🚨 Error handling middleware",
		},
		Help: "Choose the middleware components you need",
	}
	var selectedMiddleware []string
	if err := survey.AskOne(middlewarePrompt, &selectedMiddleware); err != nil {
		return nil, err
	}

	// Set middleware based on selection
	for _, middleware := range selectedMiddleware {
		switch {
		case strings.Contains(middleware, "CORS"):
			cfg.Middleware.CORS = true
		case strings.Contains(middleware, "Rate"):
			cfg.Middleware.RateLimit = true
		case strings.Contains(middleware, "Request logging"):
			cfg.Middleware.Logging = true
		case strings.Contains(middleware, "Authentication"):
			cfg.Middleware.Auth = true
		case strings.Contains(middleware, "Error"):
			cfg.Middleware.ErrorHandler = true
		}
	}

	// CI/CD selection
	cicdPrompt := &survey.Select{
		Message: "🚀 Choose your CI/CD platform:",
		Options: []string{
			fmt.Sprintf("🐙 %s - GitHub Actions workflow", config.CICDGitHub),
			fmt.Sprintf("🦊 %s - GitLab CI pipeline", config.CICDGitLab),
			fmt.Sprintf("🚫 %s - No CI/CD setup", config.CICDNone),
		},
		Default: fmt.Sprintf("🐙 %s - GitHub Actions workflow", config.CICDGitHub),
		Help:    "Select your preferred CI/CD platform for automated builds and deployments",
	}
	var selectedCICD string
	if err := survey.AskOne(cicdPrompt, &selectedCICD); err != nil {
		return nil, err
	}
	cfg.CICD = extractCICDName(selectedCICD)

	return cfg, nil
}

// Helper functions to extract names from formatted options
func extractFrameworkName(option string) string {
	switch {
	case strings.Contains(option, config.FrameworkGin):
		return config.FrameworkGin
	case strings.Contains(option, config.FrameworkEcho):
		return config.FrameworkEcho
	case strings.Contains(option, config.FrameworkFiber):
		return config.FrameworkFiber
	case strings.Contains(option, config.FrameworkRevel):
		return config.FrameworkRevel
	default:
		return config.FrameworkGin
	}
}

func extractORMName(option string) string {
	switch {
	case strings.Contains(option, config.ORMGorm):
		return config.ORMGorm
	case strings.Contains(option, config.ORMSqlx):
		return config.ORMSqlx
	case strings.Contains(option, config.ORMRaw):
		return config.ORMRaw
	case strings.Contains(option, config.ORMNone):
		return config.ORMNone
	default:
		return config.ORMGorm
	}
}

func extractDBName(option string) string {
	switch {
	case strings.Contains(option, config.DBPostgreSQL):
		return config.DBPostgreSQL
	case strings.Contains(option, config.DBMySQL):
		return config.DBMySQL
	case strings.Contains(option, config.DBSQLite):
		return config.DBSQLite
	case strings.Contains(option, config.DBMongoDB):
		return config.DBMongoDB
	case strings.Contains(option, config.DBRedis):
		return config.DBRedis
	case strings.Contains(option, config.DBMemory):
		return config.DBMemory
	default:
		return config.DBPostgreSQL
	}
}

func extractArchName(option string) string {
	switch {
	case strings.Contains(option, config.ArchSimple):
		return config.ArchSimple
	case strings.Contains(option, config.ArchClean):
		return config.ArchClean
	case strings.Contains(option, config.ArchHexagonal):
		return config.ArchHexagonal
	case strings.Contains(option, config.ArchMVC):
		return config.ArchMVC
	case strings.Contains(option, config.ArchCustom):
		return config.ArchCustom
	default:
		return config.ArchSimple
	}
}

func extractConfigName(option string) string {
	switch {
	case strings.Contains(option, config.ConfigYAML):
		return config.ConfigYAML
	case strings.Contains(option, config.ConfigJSON):
		return config.ConfigJSON
	case strings.Contains(option, config.ConfigTOML):
		return config.ConfigTOML
	default:
		return config.ConfigYAML
	}
}

func extractAuthName(option string) string {
	switch {
	case strings.Contains(option, config.AuthJWT):
		return config.AuthJWT
	case strings.Contains(option, config.AuthOAuth2):
		return config.AuthOAuth2
	case strings.Contains(option, config.AuthBasic):
		return config.AuthBasic
	case strings.Contains(option, config.AuthNone):
		return config.AuthNone
	default:
		return config.AuthJWT
	}
}

func extractLoggingName(option string) string {
	switch {
	case strings.Contains(option, config.LogCharm):
		return config.LogCharm
	case strings.Contains(option, config.LogZap):
		return config.LogZap
	case strings.Contains(option, config.LogLogrus):
		return config.LogLogrus
	case strings.Contains(option, config.LogStandard):
		return config.LogStandard
	default:
		return config.LogCharm
	}
}

func extractCICDName(option string) string {
	switch {
	case strings.Contains(option, config.CICDGitHub):
		return config.CICDGitHub
	case strings.Contains(option, config.CICDGitLab):
		return config.CICDGitLab
	case strings.Contains(option, config.CICDNone):
		return config.CICDNone
	default:
		return config.CICDGitHub
	}
}

// Print beautiful header
func printHeader() {
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	blue := color.New(color.FgBlue)

	fmt.Println()
	cyan.Println("  ╔═══════════════════════════════════════════════════════════╗")
	cyan.Println("  ║                                                           ║")
	cyan.Print("  ║           ")
	magenta.Print("🚀 Welcome to Gool")
	cyan.Print(" - Go Project Generator")
	cyan.Println("           ║")
	cyan.Println("  ║                                                           ║")
	cyan.Print("  ║       ")
	blue.Print("✨ Create modern Go applications in seconds ✨")
	cyan.Println("       ║")
	cyan.Println("  ║                                                           ║")
	cyan.Println("  ╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	color.Yellow("Let's configure your awesome Go project! 🎯\n")
}

// GetProjectPath returns the full project path
func GetProjectPath(projectName string) string {
	return filepath.Join(".", projectName)
}
