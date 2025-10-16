package spring_boot

const FilePermission = 0644
const DirPermission = 0755

const srcMainJavaPath = "src/main/java/"
const srcMainResourcesPath = "src/main/resources/"
const srcTestJavaPath = "src/test/java/"
const srcTestResourcesPath = "src/test/resources/"

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

type BuildTool string

const (
	BuildToolMaven  BuildTool = "Maven"
	BuildToolGradle BuildTool = "Gradle"
)

func (s BuildTool) String() string {
	return string(s)
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

type DbMigrationTool string

const (
	Flyway    DbMigrationTool = "Flyway"
	Liquibase DbMigrationTool = "Liquibase"
)

func (s DbMigrationTool) String() string {
	return string(s)
}

const FeatureDockerComposeSupport = "Docker Compose"
const FeatureSpringModulithSupport = "Spring Modulith"
const FeatureSpringCloudAWSSupport = "Spring Cloud AWS"
const FeatureThymeleafSupport = "Thymeleaf"
const FeatureHTMXSupport = "HTMX"
const FeatureSecuritySupport = "Security"
const FeatureJwtSecuritySupport = "JWT Security"
