package springboot

import "strings"

const (
	javaVersion                            = 25
	springBootVersion                      = "4.2.0"
	gradleDependencyMgmtVersion            = "1.1.7"
	gradleGitPropertiesVersion             = "4.0.1"
	spotlessGradlePluginVersion            = "8.12.0"
	spotlessMavenPluginVersion             = "3.12.0"
	springModulithVersion                  = "2.2.0"
	springCloudAWSVersion                  = "4.2.2"
	springBootFlociTcVersion               = "2.16.1"
	springdocOpenapiVersion                = "3.2.1"
	tcMailpitVersion                       = "1.4.0"
	fontAwesomeVersion                     = "7.3.0"
	htmxVersion                            = "4.0.2"
	htmxSpringBootThymeleafVersion         = "6.0.1"
	htmxSpringBootVersion                  = "6.0.1"
	jteVersion                             = "3.2.5"
	springAIVersion                        = "2.1.0"
	datasourceMicrometerSpringBootVersion  = "2.3.1"
	jjwtVersion                            = "0.13.0"
	bootUIVersion                          = "2.5.0"
	axelixVersion                          = "1.3.1"
	archunitVersion                        = "1.5.1"
	taikaiVersion                          = "3.0.0"
	palantirJavaFormatVersion              = "2.103.0"
	testcontainersJooqCodegenPluginVersion = "1.0.0"
	jacocoMavenPluginVersion               = "0.8.16"
	frontendMavenPluginVersion             = "2.0.2"
	gradleNodePluginVersion                = "7.1.0"
	nodeVersion                            = "26.13.1"
	npmVersion                             = "12.2.3"
	mavenAntrunPluginVersion               = "3.2.0"
	opentelemetryLogbackAppenderVersion    = "2.34.0-alpha"
	mavenVersion                           = "3.10.2"
	mavenWrapperVersion                    = "3.3.4"
	gradleVersion                          = "9.9.2"
	checkoutActionVersion                  = "v7"
	setupJavaActionVersion                 = "v6"
	setupTaskActionVersion                 = "v2"
	miseActionVersion                      = "v5"
	dockerLoginActionVersion               = "v4"
	postgresImage                          = "postgres:19-alpine"
	mysqlImage                             = "mysql:9"
	mariadbImage                           = "mariadb:12"
	flociImage                             = "floci/floci:latest-compat"
	mailpitImage                           = "axllent/mailpit:v1.34"
	rabbitmqImage                          = "rabbitmq:4-management"
	redisImage                             = "redis:8-alpine"
	grafanaLgtmImage                       = "grafana/otel-lgtm:0.41.1"
	ollamaImage                            = "ollama/ollama:0.40.17"
)

// Versions holds all third-party dependency and plugin versions used in generated projects.
type Versions struct {
	JavaVersion                            int
	SpringBootVersion                      string
	GradleDependencyMgmtVersion            string
	GradleGitPropertiesVersion             string
	SpotlessGradlePluginVersion            string
	SpotlessMavenPluginVersion             string
	SpringModulithVersion                  string
	SpringCloudAWSVersion                  string
	SpringBootFlociTcVersion               string
	SpringdocOpenapiVersion                string
	TcMailpitVersion                       string
	FontAwesomeVersion                     string
	HtmxVersion                            string
	HtmxSpringBootThymeleafVersion         string
	HtmxSpringBootVersion                  string
	JteVersion                             string
	SpringAIVersion                        string
	DatasourceMicrometerSpringBootVersion  string
	JjwtVersion                            string
	BootUIVersion                          string
	AxelixVersion                          string
	ArchunitVersion                        string
	TaikaiVersion                          string
	PalantirJavaFormatVersion              string
	TestcontainersJooqCodegenPluginVersion string
	JacocoMavenPluginVersion               string
	FrontendMavenPluginVersion             string
	GradleNodePluginVersion                string
	NodeVersion                            string
	NpmVersion                             string
	MavenAntrunPluginVersion               string
	OpenTelemetryLogbackAppenderVersion    string
	MavenVersion                           string
	MavenWrapperVersion                    string
	GradleVersion                          string
	CheckoutActionVersion                  string
	SetupJavaActionVersion                 string
	SetupTaskActionVersion                 string
	MiseActionVersion                      string
	DockerLoginActionVersion               string
	PostgresImage                          string
	MysqlImage                             string
	MariadbImage                           string
	FlociImage                             string
	MailpitImage                           string
	RabbitMQImage                          string
	RedisImage                             string
	GrafanaLgtmImage                       string
	OllamaImage                            string
}

func defaultVersions() Versions {
	return Versions{
		JavaVersion:                            javaVersion,
		SpringBootVersion:                      springBootVersion,
		GradleDependencyMgmtVersion:            gradleDependencyMgmtVersion,
		GradleGitPropertiesVersion:             gradleGitPropertiesVersion,
		SpotlessGradlePluginVersion:            spotlessGradlePluginVersion,
		SpotlessMavenPluginVersion:             spotlessMavenPluginVersion,
		SpringModulithVersion:                  springModulithVersion,
		SpringCloudAWSVersion:                  springCloudAWSVersion,
		SpringBootFlociTcVersion:               springBootFlociTcVersion,
		SpringdocOpenapiVersion:                springdocOpenapiVersion,
		TcMailpitVersion:                       tcMailpitVersion,
		FontAwesomeVersion:                     fontAwesomeVersion,
		HtmxVersion:                            htmxVersion,
		HtmxSpringBootThymeleafVersion:         htmxSpringBootThymeleafVersion,
		HtmxSpringBootVersion:                  htmxSpringBootVersion,
		JteVersion:                             jteVersion,
		SpringAIVersion:                        springAIVersion,
		DatasourceMicrometerSpringBootVersion:  datasourceMicrometerSpringBootVersion,
		BootUIVersion:                          bootUIVersion,
		AxelixVersion:                          axelixVersion,
		JjwtVersion:                            jjwtVersion,
		ArchunitVersion:                        archunitVersion,
		TaikaiVersion:                          taikaiVersion,
		PalantirJavaFormatVersion:              palantirJavaFormatVersion,
		TestcontainersJooqCodegenPluginVersion: testcontainersJooqCodegenPluginVersion,
		JacocoMavenPluginVersion:               jacocoMavenPluginVersion,
		FrontendMavenPluginVersion:             frontendMavenPluginVersion,
		GradleNodePluginVersion:                gradleNodePluginVersion,
		NodeVersion:                            nodeVersion,
		NpmVersion:                             npmVersion,
		MavenAntrunPluginVersion:               mavenAntrunPluginVersion,
		OpenTelemetryLogbackAppenderVersion:    opentelemetryLogbackAppenderVersion,
		MavenVersion:                           mavenVersion,
		MavenWrapperVersion:                    mavenWrapperVersion,
		GradleVersion:                          gradleVersion,
		CheckoutActionVersion:                  checkoutActionVersion,
		SetupJavaActionVersion:                 setupJavaActionVersion,
		SetupTaskActionVersion:                 setupTaskActionVersion,
		MiseActionVersion:                      miseActionVersion,
		DockerLoginActionVersion:               dockerLoginActionVersion,
		PostgresImage:                          postgresImage,
		MysqlImage:                             mysqlImage,
		MariadbImage:                           mariadbImage,
		FlociImage:                             flociImage,
		MailpitImage:                           mailpitImage,
		RabbitMQImage:                          rabbitmqImage,
		RedisImage:                             redisImage,
		GrafanaLgtmImage:                       grafanaLgtmImage,
		OllamaImage:                            ollamaImage,
	}
}

// NodeMajorVersion returns the major line of NodeVersion (e.g. "26"), used where
// tools pick the latest release of that line, such as mise.toml.
func (v Versions) NodeMajorVersion() string {
	major, _, _ := strings.Cut(v.NodeVersion, ".")
	return major
}
