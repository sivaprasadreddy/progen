package springboot

const srcMainJavaPath = "src/main/java/"
const srcMainResourcesPath = "src/main/resources/"
const srcTestJavaPath = "src/test/java/"
const srcTestResourcesPath = "src/test/resources/"
const ProjectConfigFile = ".progen.json"

type AppType string

const (
	RestApi                    AppType = "REST API"
	WebApp                     AppType = "Web App"
	SpringBootAngularFullStack AppType = "Spring Boot + Angular Full Stack"
)

func (a AppType) String() string {
	switch a {
	case RestApi:
		return "REST API"
	case WebApp:
		return "Web App"
	case SpringBootAngularFullStack:
		return "Spring Boot + Angular Full Stack"
	default:
		return "Unknown"
	}
}

func (a AppType) IsValid() bool {
	switch a {
	case RestApi, WebApp, SpringBootAngularFullStack:
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
