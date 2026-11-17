package app

import "github.com/sivaprasadreddy/progen/internal/sbgen/config"

// option is a choice the wizard offers: its value, its name, and a line
// describing it.
type option[T comparable] struct {
	value       T
	name        string
	description string
}

var buildToolOptions = []option[config.BuildTool]{
	{config.BuildToolMaven, "Maven", "pom.xml"},
	{config.BuildToolGradle, "Gradle", "build.gradle"},
}

var taskRunnerOptions = []option[config.TaskRunner]{
	{config.TaskRunnerMise, "Mise", "mise.toml installs the JDK and runs tasks"},
	{config.TaskRunnerTaskfile, "Taskfile", "Taskfile.yml with .sdkmanrc for SDKMAN"},
}

var frontendOptions = []option[config.FrontendType]{
	{config.FrontendNone, "None", "Backend REST API only"},
	{config.FrontendThymeleaf, "Thymeleaf", "Server-rendered web application using Thymeleaf"},
	{config.FrontendJTE, "JTE", "Server-rendered web application using Java Template Engine"},
	{config.FrontendAngular, "Angular", "Spring Boot backend with an Angular single-page application"},
	{config.FrontendReact, "React", "Spring Boot backend with a React (Vite) single-page application"},
}

var databaseOptions = []option[config.DatabaseEngine]{
	{config.DatabasePostgreSQL, "PostgreSQL", "Reliable, advanced SQL database"},
	{config.DatabaseMySQL, "MySQL", "Popular open-source SQL database"},
	{config.DatabaseMariaDB, "MariaDB", "Community-developed MySQL alternative"},
}

var persistenceOptions = []option[config.PersistenceFramework]{
	{config.PersistenceSpringDataJPA, "Spring Data JPA", "Repository-based persistence with JPA"},
	{config.PersistenceSpringJdbcClient, "Spring JdbcClient", "Fluent JDBC-based database access"},
	{config.PersistenceJOOQ, "jOOQ", "Type-safe SQL query building with jOOQ"},
}

var migrationOptions = []option[config.MigrationTool]{
	{config.MigrationFlyway, "Flyway", "Versioned SQL migrations"},
	{config.MigrationLiquibase, "Liquibase", "Flexible database change tracking"},
}

var featureOptions = []option[config.Feature]{
	{config.FeatureSpringCloudAWS, "Spring Cloud AWS", "Integrate with AWS services using Spring"},
	{config.FeatureEmail, "Email", "Send transactional and application emails"},
	{config.FeatureRabbitMQ, "RabbitMQ", "Add messaging support with Spring AMQP"},
	{config.FeatureRedisCaching, "Redis Caching", "Add caching support backed by Redis"},
	{config.FeatureOpenTelemetry, "OpenTelemetry", "Add distributed tracing and observability"},
	{config.FeatureKubernetes, "Kubernetes Manifests", "Generate deployment manifests for Kubernetes"},
	{config.FeatureAxelix, "Axelix", "Monitor and manage the service with Axelix"},
}

var aiModelOptions = []option[config.AIModel]{
	{config.AIModelOllama, "Ollama", "Run open models locally with Ollama"},
	{config.AIModelOpenAI, "OpenAI", "Use OpenAI chat models"},
	{config.AIModelAnthropic, "Anthropic", "Use Anthropic Claude chat models"},
}

// labels are the names the Review step and template summaries show.
var labels = map[string]string{
	"maven": "Maven", "gradle": "Gradle",
	"mise": "Mise", "taskfile": "Taskfile (with SDKMAN)",
	"none": "None (REST API)", "thymeleaf": "Thymeleaf", "jte": "JTE", "angular": "Angular", "react": "React",
	"postgresql": "PostgreSQL", "mysql": "MySQL", "mariadb": "MariaDB",
	"flyway": "Flyway", "liquibase": "Liquibase",
	"spring-data-jpa": "Spring Data JPA", "spring-jdbc-client": "Spring JdbcClient", "jooq": "jOOQ",
	"email": "Email", "spring-cloud-aws": "Spring Cloud AWS", "rabbitmq": "RabbitMQ",
	"redis-caching": "Redis Caching", "opentelemetry": "OpenTelemetry", "kubernetes": "Kubernetes Manifests", "axelix": "Axelix",
	"ollama": "Ollama", "openai": "OpenAI", "anthropic": "Anthropic",
}

// label is the name of an option's value.
func label[T ~string](v T) string {
	if l, ok := labels[string(v)]; ok {
		return l
	}
	return string(v)
}
