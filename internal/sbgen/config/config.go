// Package config describes the Spring Boot project a user chooses in the
// wizard and validates it.
package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
)

var javaPackagePattern = regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*(\.[a-zA-Z_$][a-zA-Z0-9_$]*)*$`)

// BuildTool is the build system of the generated project.
type BuildTool string

const (
	BuildToolMaven  BuildTool = "maven"
	BuildToolGradle BuildTool = "gradle"
)

// TaskRunner runs the generated project's common tasks.
type TaskRunner string

const (
	TaskRunnerMise     TaskRunner = "mise"
	TaskRunnerTaskfile TaskRunner = "taskfile"
)

// FrontendType is how the generated project's user interface is built.
type FrontendType string

const (
	FrontendNone      FrontendType = "none"
	FrontendThymeleaf FrontendType = "thymeleaf"
	FrontendJTE       FrontendType = "jte"
	FrontendAngular   FrontendType = "angular"
	FrontendReact     FrontendType = "react"
)

// DatabaseEngine is the database the generated project connects to.
type DatabaseEngine string

const (
	DatabasePostgreSQL DatabaseEngine = "postgresql"
	DatabaseMySQL      DatabaseEngine = "mysql"
	DatabaseMariaDB    DatabaseEngine = "mariadb"
)

// MigrationTool manages the generated project's database schema.
type MigrationTool string

const (
	MigrationFlyway    MigrationTool = "flyway"
	MigrationLiquibase MigrationTool = "liquibase"
)

// PersistenceFramework is how the generated project accesses the database.
type PersistenceFramework string

const (
	PersistenceSpringDataJPA    PersistenceFramework = "spring-data-jpa"
	PersistenceSpringJdbcClient PersistenceFramework = "spring-jdbc-client"
	PersistenceJOOQ             PersistenceFramework = "jooq"
)

// Feature is an optional capability of the generated project.
type Feature string

const (
	FeatureEmail          Feature = "email"
	FeatureSpringCloudAWS Feature = "spring-cloud-aws"
	FeatureRabbitMQ       Feature = "rabbitmq"
	FeatureRedisCaching   Feature = "redis-caching"
	FeatureOpenTelemetry  Feature = "opentelemetry"
	FeatureKubernetes     Feature = "kubernetes"
	FeatureAxelix         Feature = "axelix"
)

// AIModel is a Spring AI chat model provider.
type AIModel string

const (
	AIModelOllama    AIModel = "ollama"
	AIModelOpenAI    AIModel = "openai"
	AIModelAnthropic AIModel = "anthropic"
)

// ProjectSettings holds the project coordinates and build setup.
type ProjectSettings struct {
	AppName     string     `json:"appName"`
	AppVersion  string     `json:"appVersion"`
	GroupID     string     `json:"groupId"`
	ArtifactID  string     `json:"artifactId"`
	PackageName string     `json:"packageName"`
	BuildTool   BuildTool  `json:"buildTool"`
	TaskRunner  TaskRunner `json:"taskRunner"`
}

// FrontendSettings selects the generated project's user interface.
type FrontendSettings struct {
	Type FrontendType `json:"type"`
	HTMX bool         `json:"htmxSupport"`
}

// DatabaseSettings selects the generated project's database stack.
type DatabaseSettings struct {
	Engine               DatabaseEngine       `json:"engine"`
	MigrationTool        MigrationTool        `json:"migrationTool"`
	PersistenceFramework PersistenceFramework `json:"persistenceFramework"`
}

// GeneratorConfiguration is everything the user chose in the wizard.
type GeneratorConfiguration struct {
	Project  ProjectSettings  `json:"project"`
	Frontend FrontendSettings `json:"frontend"`
	Database DatabaseSettings `json:"database"`
	Features []Feature        `json:"features"`
	AIModels []AIModel        `json:"aiModels"`
}

// IsServerRendered reports whether the frontend is rendered by the server,
// the frontends HTMX goes with.
func IsServerRendered(frontendType FrontendType) bool {
	return frontendType == FrontendThymeleaf || frontendType == FrontendJTE
}

// Validate reports the first problem with configuration, or nil when the
// generator can generate it.
func Validate(configuration GeneratorConfiguration) error {
	project := configuration.Project
	appName := strings.TrimSpace(project.AppName)
	if appName == "" {
		return errors.New("application name is required")
	}
	if appName == "." || appName == ".." || filepath.Base(appName) != appName || strings.ContainsAny(appName, `/\`) {
		return fmt.Errorf("invalid application name %q", project.AppName)
	}
	if strings.TrimSpace(project.AppVersion) == "" {
		return errors.New("application version is required")
	}
	if strings.TrimSpace(project.GroupID) == "" {
		return errors.New("group ID is required")
	}
	if strings.TrimSpace(project.ArtifactID) == "" {
		return errors.New("artifact ID is required")
	}
	if !javaPackagePattern.MatchString(strings.TrimSpace(project.PackageName)) {
		return fmt.Errorf("invalid Java package name %q", project.PackageName)
	}
	if !slices.Contains([]BuildTool{BuildToolMaven, BuildToolGradle}, project.BuildTool) {
		return fmt.Errorf("unsupported build tool %q", project.BuildTool)
	}
	if !slices.Contains([]TaskRunner{TaskRunnerMise, TaskRunnerTaskfile}, project.TaskRunner) {
		return fmt.Errorf("unsupported task runner %q", project.TaskRunner)
	}

	frontend := configuration.Frontend
	if !slices.Contains([]FrontendType{FrontendNone, FrontendThymeleaf, FrontendJTE, FrontendAngular, FrontendReact}, frontend.Type) {
		return fmt.Errorf("unsupported frontend type %q", frontend.Type)
	}
	if frontend.HTMX && !IsServerRendered(frontend.Type) {
		return errors.New("HTMX support requires the Thymeleaf or JTE frontend")
	}

	database := configuration.Database
	if !slices.Contains([]DatabaseEngine{DatabasePostgreSQL, DatabaseMySQL, DatabaseMariaDB}, database.Engine) {
		return fmt.Errorf("unsupported database engine %q", database.Engine)
	}
	if !slices.Contains([]MigrationTool{MigrationFlyway, MigrationLiquibase}, database.MigrationTool) {
		return fmt.Errorf("unsupported migration tool %q", database.MigrationTool)
	}
	if !slices.Contains([]PersistenceFramework{PersistenceSpringDataJPA, PersistenceSpringJdbcClient, PersistenceJOOQ}, database.PersistenceFramework) {
		return fmt.Errorf("unsupported persistence framework %q", database.PersistenceFramework)
	}

	supportedFeatures := []Feature{FeatureEmail, FeatureSpringCloudAWS, FeatureRabbitMQ, FeatureRedisCaching, FeatureOpenTelemetry, FeatureKubernetes, FeatureAxelix}
	for _, feature := range configuration.Features {
		if !slices.Contains(supportedFeatures, feature) {
			return fmt.Errorf("unsupported feature %q", feature)
		}
	}
	for _, model := range configuration.AIModels {
		if !springboot.AIModelProvider(model).IsValid() {
			return fmt.Errorf("unsupported Spring AI model provider %q", model)
		}
	}
	return nil
}
