package generator

import (
	"path/filepath"

	"github.com/gool-cli/gool/internal/config"
	"github.com/gool-cli/gool/internal/templates"
)

// generateModels generates example models
func (g *Generator) generateModels(cfg *config.ProjectConfig, projectPath string) error {
	data := templates.NewTemplateData(cfg)

	// Generate User model
	userModelTemplate := `package models

import (
	{{- if eq .ORM "gorm"}}
	"gorm.io/gorm"
	{{- end}}
	"time"
)

type User struct {
	{{- if eq .ORM "gorm"}}
	gorm.Model
	{{- else}}
	ID        uint      ` + "`json:\"id\" db:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
	{{- end}}
	Name      string    ` + "`json:\"name\" db:\"name\" validate:\"required,min=2,max=100\"`" + `
	Email     string    ` + "`json:\"email\" db:\"email\" validate:\"required,email\" gorm:\"unique\"`" + `
	Password  string    ` + "`json:\"-\" db:\"password\" validate:\"required,min=6\"`" + `
	Role      string    ` + "`json:\"role\" db:\"role\" validate:\"required\" gorm:\"default:user\"`" + `
	IsActive  bool      ` + "`json:\"is_active\" db:\"is_active\" gorm:\"default:true\"`" + `
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserResponse represents the user response without sensitive data
type UserResponse struct {
	ID       uint      ` + "`json:\"id\"`" + `
	Name     string    ` + "`json:\"name\"`" + `
	Email    string    ` + "`json:\"email\"`" + `
	Role     string    ` + "`json:\"role\"`" + `
	IsActive bool      ` + "`json:\"is_active\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Role:     u.Role,
		IsActive: u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CreateUserRequest represents request for creating user
type CreateUserRequest struct {
	Name     string ` + "`json:\"name\" validate:\"required,min=2,max=100\"`" + `
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\" validate:\"required,min=6\"`" + `
	Role     string ` + "`json:\"role\" validate:\"required\"`" + `
}

// UpdateUserRequest represents request for updating user
type UpdateUserRequest struct {
	Name     *string ` + "`json:\"name,omitempty\" validate:\"omitempty,min=2,max=100\"`" + `
	Email    *string ` + "`json:\"email,omitempty\" validate:\"omitempty,email\"`" + `
	Role     *string ` + "`json:\"role,omitempty\"`" + `
	IsActive *bool   ` + "`json:\"is_active,omitempty\"`" + `
}

{{- if eq .Config.Auth "jwt"}}
// LoginRequest represents login request
type LoginRequest struct {
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\" validate:\"required\"`" + `
}

// LoginResponse represents login response
type LoginResponse struct {
	User         UserResponse ` + "`json:\"user\"`" + `
	Token        string       ` + "`json:\"token\"`" + `
	RefreshToken string       ` + "`json:\"refresh_token,omitempty\"`" + `
	ExpiresIn    int64        ` + "`json:\"expires_in\"`" + `
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string ` + "`json:\"refresh_token\" validate:\"required\"`" + `
}
{{- end}}
`

	var modelPath string
	switch cfg.Architecture {
	case config.ArchSimple:
		modelPath = "internal/models"
	case config.ArchClean:
		modelPath = "internal/entity"
	case config.ArchHexagonal:
		modelPath = "internal/core/domain"
	case config.ArchMVC:
		modelPath = "internal/models"
	default:
		modelPath = "internal/models"
	}

	if err := g.templateEngine.RenderToFile(userModelTemplate, filepath.Join(projectPath, modelPath, "user.go"), data); err != nil {
		return err
	}

	// Generate common response types
	responseTemplate := `package models

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        ` + "`json:\"success\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
	Error   *APIError   ` + "`json:\"error,omitempty\"`" + `
}

// APIError represents an API error
type APIError struct {
	Code    string                 ` + "`json:\"code\"`" + `
	Message string                 ` + "`json:\"message\"`" + `
	Details map[string]interface{} ` + "`json:\"details,omitempty\"`" + `
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   ` + "`json:\"page\"`" + `
	Limit      int   ` + "`json:\"limit\"`" + `
	Total      int64 ` + "`json:\"total\"`" + `
	TotalPages int   ` + "`json:\"total_pages\"`" + `
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{}     ` + "`json:\"data\"`" + `
	Pagination *PaginationMeta ` + "`json:\"pagination\"`" + `
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status    string                 ` + "`json:\"status\"`" + `
	Message   string                 ` + "`json:\"message\"`" + `
	Timestamp string                 ` + "`json:\"timestamp\"`" + `
	Services  map[string]interface{} ` + "`json:\"services,omitempty\"`" + `
}
`

	if err := g.templateEngine.RenderToFile(responseTemplate, filepath.Join(projectPath, modelPath, "response.go"), data); err != nil {
		return err
	}

	return nil
}
