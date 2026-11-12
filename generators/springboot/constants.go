package springboot

const srcMainJavaPath = "src/main/java/"
const srcMainResourcesPath = "src/main/resources/"
const srcTestJavaPath = "src/test/java/"
const srcTestResourcesPath = "src/test/resources/"
const ProjectConfigFile = ".progen.json"

// Frontend determines the kind of UI generated alongside the Spring Boot backend.
// FrontendNone generates a REST API only.
type Frontend string

const (
	FrontendNone      Frontend = "None"
	FrontendThymeleaf Frontend = "Thymeleaf"
	FrontendJTE       Frontend = "JTE"
	FrontendAngular   Frontend = "Angular"
)

func (f Frontend) String() string {
	return string(f)
}

func (f Frontend) IsValid() bool {
	switch f {
	case FrontendNone, FrontendThymeleaf, FrontendJTE, FrontendAngular:
		return true
	default:
		return false
	}
}

type BuildTool string

const (
	Maven  BuildTool = "Maven"
	Gradle BuildTool = "Gradle"
)

func (s BuildTool) String() string {
	return string(s)
}

func (s BuildTool) IsValid() bool {
	switch s {
	case Maven, Gradle:
		return true
	default:
		return false
	}
}

type DatabaseType string

type PersistenceType string

const (
	SpringDataJPA    PersistenceType = "Spring Data JPA"
	SpringJdbcClient PersistenceType = "Spring JdbcClient"
	SpringJOOQ       PersistenceType = "jOOQ"
)

func (p PersistenceType) String() string { return string(p) }

func (p PersistenceType) IsValid() bool {
	return p == SpringDataJPA || p == SpringJdbcClient || p == SpringJOOQ
}

const (
	PostgreSQL DatabaseType = "PostgreSQL"
	MySQL      DatabaseType = "MySQL"
	MariaDB    DatabaseType = "MariaDB"
)

func (s DatabaseType) String() string {
	return string(s)
}

func (s DatabaseType) IsValid() bool {
	switch s {
	case PostgreSQL, MySQL, MariaDB:
		return true
	default:
		return false
	}
}

type DbMigrationTool string

const (
	Flyway    DbMigrationTool = "Flyway"
	Liquibase DbMigrationTool = "Liquibase"
)

func (s DbMigrationTool) String() string {
	return string(s)
}

func (s DbMigrationTool) IsValid() bool {
	switch s {
	case Flyway, Liquibase:
		return true
	default:
		return false
	}
}

const FeatureSpringCloudAWSSupport = "Spring Cloud AWS"
const FeatureHTMXSupport = "HTMX"
const FeatureEmailSupport = "Email"
const FeatureRabbitMQSupport = "RabbitMQ"
const FeatureRedisCachingSupport = "Redis Caching"
const FeatureOpenTelemetrySupport = "OpenTelemetry"
const FeatureK8sSupport = "Kubernetes Manifests"
