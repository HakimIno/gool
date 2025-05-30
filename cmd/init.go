package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/generator"
	"github.com/gool-cli/gool/internal/prompts"
	"github.com/spf13/cobra"
)

var (
	projectName string
	interactive bool
	framework   string
	orm         string
	database    string
	arch        string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Go project",
	Long: `Initialize a new Go project with customizable structure and configurations.
	
🎯 Interactive Mode (Default):
  gool init my-app                    # Full interactive setup
  gool init                          # Interactive with default name

🚀 Quick Mode:
  gool init my-app --framework=gin --orm=gorm --database=postgresql --arch=simple --interactive=false

✨ Examples:
  gool init my-api --framework=gin --database=postgresql
  gool init my-service --arch=clean --interactive=false
  gool init my-microservice --framework=echo --orm=sqlx`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags for non-interactive mode
	initCmd.Flags().StringVarP(&framework, "framework", "f", "", "Web framework (gin, echo, fiber, revel)")
	initCmd.Flags().StringVarP(&orm, "orm", "o", "", "ORM/Database layer (gorm, sqlx, raw, none)")
	initCmd.Flags().StringVarP(&database, "database", "d", "", "Database type (postgresql, mysql, sqlite, mongodb, redis, memory)")
	initCmd.Flags().StringVarP(&arch, "arch", "a", "", "Architecture (simple, clean, hexagonal, mvc, custom)")
	initCmd.Flags().BoolVar(&interactive, "interactive", true, "Run in interactive mode (default: true)")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Get project name from args or use default
	if len(args) > 0 {
		projectName = args[0]
	}

	var cfg *config.ProjectConfig
	var err error

	// Check if user wants non-interactive mode by providing flags
	isNonInteractive := !interactive || (framework != "" || orm != "" || database != "" || arch != "")

	if !isNonInteractive {
		// Interactive mode (default)
		cfg, err = prompts.CollectProjectConfig()
		if err != nil {
			color.Red("❌ Failed to collect configuration: %v", err)
			return err
		}

		// Override project name if provided as argument
		if projectName != "" {
			cfg.ProjectName = projectName
		}
	} else {
		// Non-interactive mode with flags
		color.Cyan("🚀 Quick Setup Mode")
		fmt.Println()

		cfg = &config.ProjectConfig{
			ProjectName:  projectName,
			Framework:    framework,
			ORM:          orm,
			Database:     database,
			Architecture: arch,
		}

		// Validate and set defaults for missing values
		if err := validateAndSetDefaults(cfg); err != nil {
			color.Red("❌ Configuration error: %v", err)
			printConfigurationHelp()
			return err
		}

		// Show quick setup summary
		printQuickSetupSummary(cfg)
	}

	// Validate configuration
	if err := validateProjectConfig(cfg); err != nil {
		color.Red("❌ Configuration validation failed: %v", err)
		printConfigurationHelp()
		return err
	}

	// Check if project directory already exists
	projectPath := prompts.GetProjectPath(cfg.ProjectName)
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		color.Red("❌ Directory '%s' already exists", projectPath)
		color.Yellow("💡 Try using a different project name or remove the existing directory:")
		color.White("   rm -rf %s", projectPath)
		return fmt.Errorf("directory '%s' already exists", projectPath)
	}

	// Generate the project
	printGenerationStart(cfg)

	gen := generator.New()
	if err := gen.Generate(cfg); err != nil {
		color.Red("❌ Failed to generate project: %v", err)

		// Provide specific help based on error type
		printErrorHelp(err)

		// Clean up partial generation if it exists
		if _, statErr := os.Stat(projectPath); statErr == nil {
			color.Yellow("🧹 Cleaning up partial generation...")
			if removeErr := os.RemoveAll(projectPath); removeErr != nil {
				color.Red("❌ Failed to clean up: %v", removeErr)
			}
		}

		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Success message
	printSuccessMessage(cfg, projectPath)

	return nil
}

// validateAndSetDefaults validates configuration and sets defaults
func validateAndSetDefaults(cfg *config.ProjectConfig) error {
	// Set defaults for missing values
	if cfg.ProjectName == "" {
		cfg.ProjectName = "my-go-app"
	}

	// Validate project name
	if !isValidProjectName(cfg.ProjectName) {
		return fmt.Errorf("invalid project name '%s'. Use only letters, numbers, hyphens, and underscores", cfg.ProjectName)
	}

	if cfg.Framework == "" {
		cfg.Framework = config.FrameworkGin
	} else if !isValidFramework(cfg.Framework) {
		return fmt.Errorf("invalid framework '%s'. Valid options: gin, echo, fiber, revel", cfg.Framework)
	}

	if cfg.ORM == "" {
		cfg.ORM = config.ORMGorm
	} else if !isValidORM(cfg.ORM) {
		return fmt.Errorf("invalid ORM '%s'. Valid options: gorm, sqlx, raw, none", cfg.ORM)
	}

	if cfg.Database == "" && cfg.ORM != config.ORMNone {
		cfg.Database = config.DBPostgreSQL
	} else if cfg.Database != "" && !isValidDatabase(cfg.Database) {
		return fmt.Errorf("invalid database '%s'. Valid options: postgresql, mysql, sqlite, mongodb, redis, memory", cfg.Database)
	}

	if cfg.Architecture == "" {
		cfg.Architecture = config.ArchSimple
	} else if !isValidArchitecture(cfg.Architecture) {
		return fmt.Errorf("invalid architecture '%s'. Valid options: simple, clean, hexagonal, mvc, custom", cfg.Architecture)
	}

	// Set module path
	cfg.ModulePath = fmt.Sprintf("github.com/username/%s", cfg.ProjectName)
	cfg.Config = config.ConfigYAML
	cfg.Auth = config.AuthJWT
	cfg.Logging = config.LogZap
	cfg.Testing = true
	cfg.Docker = true
	cfg.CICD = config.CICDGitHub

	// Enable common middleware
	cfg.Middleware.CORS = true
	cfg.Middleware.Logging = true
	cfg.Middleware.ErrorHandler = true

	// Enable common features
	cfg.Features.HealthCheck = true
	cfg.Features.Swagger = true

	return nil
}

// validateProjectConfig validates the entire project configuration
func validateProjectConfig(cfg *config.ProjectConfig) error {
	if cfg.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if cfg.ModulePath == "" {
		return fmt.Errorf("module path is required")
	}

	// Check ORM and Database compatibility
	if cfg.ORM == config.ORMNone && cfg.Database != "" {
		color.Yellow("⚠️  Warning: Database '%s' specified but ORM is 'none'. Database will be ignored.", cfg.Database)
		cfg.Database = ""
	}

	return nil
}

// Validation helper functions
func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}
	return true
}

func isValidFramework(framework string) bool {
	validFrameworks := []string{config.FrameworkGin, config.FrameworkEcho, config.FrameworkFiber, config.FrameworkRevel}
	for _, valid := range validFrameworks {
		if framework == valid {
			return true
		}
	}
	return false
}

func isValidORM(orm string) bool {
	validORMs := []string{config.ORMGorm, config.ORMSqlx, config.ORMRaw, config.ORMNone}
	for _, valid := range validORMs {
		if orm == valid {
			return true
		}
	}
	return false
}

func isValidDatabase(database string) bool {
	validDatabases := []string{config.DBPostgreSQL, config.DBMySQL, config.DBSQLite, config.DBMongoDB, config.DBRedis, config.DBMemory}
	for _, valid := range validDatabases {
		if database == valid {
			return true
		}
	}
	return false
}

func isValidArchitecture(arch string) bool {
	validArchs := []string{config.ArchSimple, config.ArchClean, config.ArchHexagonal, config.ArchMVC, config.ArchCustom}
	for _, valid := range validArchs {
		if arch == valid {
			return true
		}
	}
	return false
}

// printErrorHelp provides specific help based on error type
func printErrorHelp(err error) {
	yellow := color.New(color.FgYellow, color.Bold)
	white := color.New(color.FgWhite)

	fmt.Println()
	yellow.Println("💡 Troubleshooting Tips:")

	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "template"):
		white.Println("  • Template error detected. This might be a bug in gool.")
		white.Println("  • Try using quick mode: --interactive=false")
		white.Println("  • Report this issue at: https://github.com/gool-cli/gool/issues")

	case strings.Contains(errStr, "permission"):
		white.Println("  • Permission denied. Check if you have write access to this directory.")
		white.Println("  • Try: sudo chown -R $USER:$USER .")

	case strings.Contains(errStr, "space"):
		white.Println("  • Not enough disk space. Free up some space and try again.")

	case strings.Contains(errStr, "network"):
		white.Println("  • Network error. Check your internet connection.")
		white.Println("  • Some templates may require internet access for dependencies.")

	default:
		white.Println("  • Check if you have sufficient permissions in this directory")
		white.Println("  • Ensure you have enough disk space")
		white.Println("  • Try running with different configuration options")
		white.Println("  • For help: gool init --help")
	}
	fmt.Println()
}

// printConfigurationHelp shows configuration help
func printConfigurationHelp() {
	yellow := color.New(color.FgYellow, color.Bold)
	white := color.New(color.FgWhite)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	yellow.Println("📚 Configuration Help:")
	fmt.Println()

	cyan.Println("Valid Frameworks:")
	white.Println("  gin, echo, fiber, revel")
	fmt.Println()

	cyan.Println("Valid ORMs:")
	white.Println("  gorm, sqlx, raw, none")
	fmt.Println()

	cyan.Println("Valid Databases:")
	white.Println("  postgresql, mysql, sqlite, mongodb, redis, memory")
	fmt.Println()

	cyan.Println("Valid Architectures:")
	white.Println("  simple, clean, hexagonal, mvc, custom")
	fmt.Println()

	yellow.Println("💡 Examples:")
	white.Println("  gool init my-app --framework=gin --database=postgresql")
	white.Println("  gool init my-service --arch=clean --orm=gorm")
	white.Println("  gool init my-api --interactive=true")
	fmt.Println()
}

func printQuickSetupSummary(cfg *config.ProjectConfig) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	cyan.Println("📋 Quick Setup Summary:")
	yellow.Printf("  • Project: %s\n", cfg.ProjectName)
	yellow.Printf("  • Framework: %s\n", cfg.Framework)
	yellow.Printf("  • ORM: %s\n", cfg.ORM)
	if cfg.Database != "" {
		yellow.Printf("  • Database: %s\n", cfg.Database)
	}
	yellow.Printf("  • Architecture: %s\n", cfg.Architecture)
	fmt.Println()
}

func printGenerationStart(cfg *config.ProjectConfig) {
	magenta := color.New(color.FgMagenta, color.Bold)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	magenta.Println("🔨 Generating your awesome Go project...")
	fmt.Println()
	cyan.Printf("  📦 Creating project structure...\n")
	cyan.Printf("  🏗️  Setting up %s framework...\n", cfg.Framework)
	if cfg.ORM != config.ORMNone {
		cyan.Printf("  🗄️  Configuring %s with %s...\n", cfg.ORM, cfg.Database)
	}
	cyan.Printf("  🏛️  Applying %s architecture...\n", cfg.Architecture)
	cyan.Printf("  ⚙️  Adding configuration files...\n")
	if cfg.Docker {
		cyan.Printf("  🐳 Adding Docker support...\n")
	}
	if cfg.Features.Swagger {
		cyan.Printf("  📚 Setting up API documentation...\n")
	}
	cyan.Printf("  ✨ Adding finishing touches...\n")
	fmt.Println()
}

func printSuccessMessage(cfg *config.ProjectConfig, projectPath string) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	white := color.New(color.FgWhite)
	magenta := color.New(color.FgMagenta)

	fmt.Println()
	green.Println("🎉 Project generated successfully!")
	fmt.Println()

	cyan.Printf("📁 Project location: %s\n", projectPath)
	fmt.Println()

	yellow.Println("🚀 Next steps:")
	white.Printf("  cd %s\n", cfg.ProjectName)
	white.Printf("  go mod tidy\n")
	white.Printf("  go run main.go\n")
	fmt.Println()

	if cfg.Features.Swagger {
		magenta.Println("📚 API Documentation:")
		white.Printf("  http://localhost:8080/swagger/index.html\n")
		fmt.Println()
	}

	if cfg.Features.HealthCheck {
		magenta.Println("❤️  Health Check:")
		white.Printf("  http://localhost:8080/api/v1/health\n")
		fmt.Println()
	}

	if cfg.Docker {
		yellow.Println("🐳 Docker commands:")
		white.Printf("  docker-compose up -d     # Start services\n")
		white.Printf("  make docker-build        # Build image\n")
		white.Printf("  make docker-run          # Run container\n")
		fmt.Println()
	}

	yellow.Println("🛠️  Development commands:")
	white.Printf("  make help                 # Show all commands\n")
	white.Printf("  make run                  # Run application\n")
	white.Printf("  make test                 # Run tests\n")
	white.Printf("  make dev                  # Hot reload (requires air)\n")
	fmt.Println()

	green.Println("Happy coding! 🎯")
}
