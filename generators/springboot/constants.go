package springboot

const srcMainJavaPath = "src/main/java/"
const srcMainResourcesPath = "src/main/resources/"
const srcTestJavaPath = "src/test/java/"
const srcTestResourcesPath = "src/test/resources/"
const ProjectConfigFile = ".progen.json"

type AppType string

const (
	RestApi AppType = "REST API"
	WebApp  AppType = "Web App"
)

func (a AppType) String() string {
	switch a {
	case RestApi:
		return "REST API"
	case WebApp:
		return "Web App"
	default:
		return "Unknown"
	}
}

func (a AppType) IsValid() bool {
	switch a {
	case RestApi, WebApp:
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
const FeatureThymeleafSupport = "Thymeleaf"
const FeatureHTMXSupport = "HTMX"
