package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// Generator handles project generation
type Generator struct {
	templateEngine *templates.Engine
}

// New creates a new generator instance
func New() *Generator {
	return &Generator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate creates a new project based on the provided configuration
func (g *Generator) Generate(cfg *config.ProjectConfig) error {
	projectPath := filepath.Join(".", cfg.ProjectName)

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate folder structure based on architecture
	if err := g.generateDirectoryStructure(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate directory structure: %w", err)
	}

	// Generate core files
	if err := g.generateCoreFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate core files: %w", err)
	}

	// Generate framework-specific files
	if err := g.generateFrameworkFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate framework files: %w", err)
	}

	// Generate database files
	if err := g.generateDatabaseFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate database files: %w", err)
	}

	// Generate middleware files
	if err := g.generateMiddlewareFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate middleware files: %w", err)
	}

	// Generate feature files
	if err := g.generateFeatureFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate feature files: %w", err)
	}

	// Generate beautiful startup
	if err := g.generateBeautifulStartup(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate beautiful startup: %w", err)
	}

	// Generate docs package (only if Swagger is enabled)
	if err := g.generateDocsPackage(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate docs package: %w", err)
	}

	// Generate models
	if err := g.generateModels(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	// Generate configuration files
	if err := g.generateConfigFiles(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate config files: %w", err)
	}

	// Generate testing files
	if cfg.Testing {
		if err := g.generateTestFiles(cfg, projectPath); err != nil {
			return fmt.Errorf("failed to generate test files: %w", err)
		}
	}

	// Generate Docker files
	if cfg.Docker {
		if err := g.generateDockerFiles(cfg, projectPath); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
		}
	}

	// Generate CI/CD files
	if cfg.CICD != config.CICDNone {
		if err := g.generateCICDFiles(cfg, projectPath); err != nil {
			return fmt.Errorf("failed to generate CI/CD files: %w", err)
		}
	}

	// Generate README
	if err := g.generateREADME(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	return nil
}

// generateDirectoryStructure creates the folder structure based on architecture
func (g *Generator) generateDirectoryStructure(cfg *config.ProjectConfig, projectPath string) error {
	var dirs []string

	switch cfg.Architecture {
	case config.ArchSimple:
		dirs = []string{
			"cmd",
			"internal/handlers",
			"internal/models",
			"internal/services",
			"internal/middleware",
			"pkg/config",
			"pkg/database",
			"pkg/logger",
			"api/routes",
			"scripts",
			"deployments",
		}
		// Add docs directory only if Swagger is enabled
		if cfg.Features.Swagger {
			dirs = append(dirs, "docs")
		}
		if cfg.Testing {
			dirs = append(dirs, "test")
		}

	case config.ArchClean:
		dirs = []string{
			"cmd",
			"internal/controller",
			"internal/usecase",
			"internal/repository",
			"internal/entity",
			"internal/delivery/http",
			"pkg/config",
			"pkg/database",
			"pkg/logger",
			"api/routes",
			"scripts",
			"deployments",
		}
		// Add docs directory only if Swagger is enabled
		if cfg.Features.Swagger {
			dirs = append(dirs, "docs")
		}
		if cfg.Testing {
			dirs = append(dirs, "test")
		}

	case config.ArchHexagonal:
		dirs = []string{
			"cmd",
			"internal/adapters/primary/http",
			"internal/adapters/secondary/database",
			"internal/domain/entities",
			"internal/domain/services",
			"internal/ports",
			"pkg/config",
			"pkg/database",
			"pkg/logger",
			"api/routes",
			"scripts",
			"deployments",
		}
		// Add docs directory only if Swagger is enabled
		if cfg.Features.Swagger {
			dirs = append(dirs, "docs")
		}
		if cfg.Testing {
			dirs = append(dirs, "test")
		}

	case config.ArchMVC:
		dirs = []string{
			"cmd",
			"internal/controllers",
			"internal/models",
			"internal/views",
			"internal/middleware",
			"pkg/config",
			"pkg/database",
			"pkg/logger",
			"api/routes",
			"scripts",
			"deployments",
		}
		// Add docs directory only if Swagger is enabled
		if cfg.Features.Swagger {
			dirs = append(dirs, "docs")
		}
		if cfg.Testing {
			dirs = append(dirs, "test")
		}

	case config.ArchCustom:
		dirs = []string{
			"cmd",
			"internal",
			"pkg/config",
			"pkg/database",
			"pkg/logger",
			"api",
			"scripts",
			"deployments",
		}
		// Add docs directory only if Swagger is enabled
		if cfg.Features.Swagger {
			dirs = append(dirs, "docs")
		}
		if cfg.Testing {
			dirs = append(dirs, "test")
		}
	}

	// Add feature-specific directories
	if cfg.Features.StaticFiles {
		dirs = append(dirs, "static/css", "static/js", "static/images")
	}
	if cfg.Features.I18n {
		dirs = append(dirs, "locales")
	}

	// Create all directories
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
