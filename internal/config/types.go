package config

// ProjectConfig represents the configuration for generating a new project
type ProjectConfig struct {
	ProjectName  string           `yaml:"project_name"`
	ModulePath   string           `yaml:"module_path"`
	Framework    string           `yaml:"framework"`
	ORM          string           `yaml:"orm"`
	Database     string           `yaml:"database"`
	Architecture string           `yaml:"architecture"`
	Config       string           `yaml:"config"`
	Auth         string           `yaml:"auth"`
	Logging      string           `yaml:"logging"`
	Testing      bool             `yaml:"testing"`
	Docker       bool             `yaml:"docker"`
	CICD         string           `yaml:"cicd"`
	Middleware   MiddlewareConfig `yaml:"middleware"`
	Features     FeaturesConfig   `yaml:"features"`
}

// MiddlewareConfig represents middleware options
type MiddlewareConfig struct {
	CORS         bool `yaml:"cors"`
	RateLimit    bool `yaml:"rate_limit"`
	Logging      bool `yaml:"logging"`
	Auth         bool `yaml:"auth"`
	ErrorHandler bool `yaml:"error_handler"`
}

// FeaturesConfig represents additional features
type FeaturesConfig struct {
	WebSocket    bool `yaml:"websocket"`
	Caching      bool `yaml:"caching"`
	MessageQueue bool `yaml:"message_queue"`
	HealthCheck  bool `yaml:"health_check"`
	Swagger      bool `yaml:"swagger"`
	StaticFiles  bool `yaml:"static_files"`
	I18n         bool `yaml:"i18n"`
	Metrics      bool `yaml:"metrics"`
	CloudConfig  bool `yaml:"cloud_config"`
}

// Framework options
const (
	FrameworkGin   = "gin"
	FrameworkEcho  = "echo"
	FrameworkFiber = "fiber"
	FrameworkRevel = "revel"
)

// ORM options
const (
	ORMGorm = "gorm"
	ORMSqlx = "sqlx"
	ORMRaw  = "raw"
	ORMNone = "none"
)

// Database options
const (
	DBPostgreSQL = "postgresql"
	DBMySQL      = "mysql"
	DBSQLite     = "sqlite"
	DBMongoDB    = "mongodb"
	DBRedis      = "redis"
	DBMemory     = "memory"
)

// Architecture options
const (
	ArchSimple    = "simple"
	ArchClean     = "clean"
	ArchHexagonal = "hexagonal"
	ArchMVC       = "mvc"
	ArchCustom    = "custom"
)

// Config format options
const (
	ConfigYAML = "yaml"
	ConfigJSON = "json"
	ConfigTOML = "toml"
)

// Auth options
const (
	AuthJWT    = "jwt"
	AuthOAuth2 = "oauth2"
	AuthBasic  = "basic"
	AuthNone   = "none"
)

// Logging options
const (
	LogStandard = "standard"
	LogLogrus   = "logrus"
	LogZap      = "zap"
	LogCharm    = "charm"
)

// CI/CD options
const (
	CICDGitHub = "github"
	CICDGitLab = "gitlab"
	CICDNone   = "none"
)
