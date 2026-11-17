package springboot

import "strings"

const (
	javaVersion                            = 25
	springBootVersion                      = "4.1.2"
	gradleDependencyMgmtVersion            = "1.1.7"
	gradleGitPropertiesVersion             = "4.0.1"
	spotlessGradlePluginVersion            = "8.11.0"
	spotlessMavenPluginVersion             = "3.11.0"
	springModulithVersion                  = "2.1.2"
	springCloudAWSVersion                  = "4.2.2"
	springBootFlociTcVersion               = "2.16.1"
	springdocOpenapiVersion                = "3.2.0"
	tcMailpitVersion                       = "1.3.2"
	fontAwesomeVersion                     = "7.3.0"
	htmxVersion                            = "4.0.1"
	htmxSpringBootThymeleafVersion         = "6.0.1"
	htmxSpringBootVersion                  = "6.0.1"
	jteVersion                             = "3.2.5"
	springAIVersion                        = "2.0.2"
	datasourceMicrometerSpringBootVersion  = "2.3.0"
	jjwtVersion                            = "0.13.0"
	bootUIVersion                          = "2.4.0"
	axelixVersion                          = "1.3.0"
	archunitVersion                        = "1.5.1"
	taikaiVersion                          = "2.9.0"
	palantirJavaFormatVersion              = "2.103.0"
	testcontainersJooqCodegenPluginVersion = "1.0.0"
	jacocoMavenPluginVersion               = "0.8.16"
	frontendMavenPluginVersion             = "2.0.2"
	gradleNodePluginVersion                = "7.1.0"
	nodeVersion                            = "26.0.0"
	npmVersion                             = "12.2.2"
	mavenAntrunPluginVersion               = "3.2.0"
	opentelemetryLogbackAppenderVersion    = "2.33.1-alpha"
	opentelemetryAPIIncubatorVersion       = "1.67.0-alpha"
	mavenVersion                           = "3.10.1"
	mavenWrapperVersion                    = "3.3.4"
	gradleVersion                          = "9.9.1"
	checkoutActionVersion                  = "v7"
	setupJavaActionVersion                 = "v6"
	setupTaskActionVersion                 = "v2"
	miseActionVersion                      = "v5"
	dockerLoginActionVersion               = "v4"
	postgresImage                          = "postgres:18-alpine"
	mysqlImage                             = "mysql:9"
	mariadbImage                           = "mariadb:12"
	flociImage                             = "floci/floci:latest-compat"
	mailpitImage                           = "axllent/mailpit:v1.32"
	rabbitmqImage                          = "rabbitmq:4-management"
	redisImage                             = "redis:8-alpine"
	grafanaLgtmImage                       = "grafana/otel-lgtm:0.38.0"
	ollamaImage                            = "ollama/ollama:0.40.16"
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
	OpenTelemetryAPIIncubatorVersion       string
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
		OpenTelemetryAPIIncubatorVersion:       opentelemetryAPIIncubatorVersion,
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
