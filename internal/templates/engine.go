package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gool-cli/gool/internal/config"
)

// Engine handles template processing
type Engine struct {
	templates map[string]*template.Template
}

// NewEngine creates a new template engine
func NewEngine() *Engine {
	return &Engine{
		templates: make(map[string]*template.Template),
	}
}

// getFuncMap returns template functions
func (e *Engine) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"capitalize": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},
	}
}

// RenderToFile renders a template to a file
func (e *Engine) RenderToFile(templateContent, filePath string, data interface{}) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Parse template with custom functions
	tmpl, err := template.New("template").Funcs(e.getFuncMap()).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// WriteFile writes content directly to a file
func (e *Engine) WriteFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// TemplateData holds common template data
type TemplateData struct {
	Config      *config.ProjectConfig
	ProjectName string
	ModulePath  string
	Framework   string
	ORM         string
	Database    string
}

// NewTemplateData creates template data from config
func NewTemplateData(cfg *config.ProjectConfig) *TemplateData {
	return &TemplateData{
		Config:      cfg,
		ProjectName: cfg.ProjectName,
		ModulePath:  cfg.ModulePath,
		Framework:   cfg.Framework,
		ORM:         cfg.ORM,
		Database:    cfg.Database,
	}
} 