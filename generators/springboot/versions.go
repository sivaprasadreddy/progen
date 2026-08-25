package springboot

const (
	javaVersion                           = 25
	springBootVersion                     = "4.1.1"
	gradleDependencyMgmtVersion           = "1.1.7"
	gradleGitPropertiesVersion            = "4.0.1"
	spotlessGradlePluginVersion           = "8.10.0"
	spotlessMavenPluginVersion            = "3.10.0"
	springModulithVersion                 = "2.1.0"
	springCloudAWSVersion                 = "4.1.0"
	springBootFlociTcVersion              = "2.15.0"
	springdocOpenapiVersion               = "3.1.0"
	tcMailpitVersion                      = "1.3.1"
	fontAwesomeVersion                    = "7.3.0"
	htmxVersion                           = "2.0.10"
	htmxSpringBootThymeleafVersion        = "5.1.0"
	datasourceMicrometerSpringBootVersion = "2.2.1"
	jjwtVersion                           = "0.13.0"
	bootUIVersion                         = "1.14.1"
	archunitVersion                       = "1.5.0"
	taikaiVersion                         = "1.66.0"
	palantirJavaFormatVersion             = "2.97.0"
	postgresImage                         = "postgres:18-alpine"
	mysqlImage                            = "mysql:9"
	mariadbImage                          = "mariadb:12"
	flociImage                            = "floci/floci:latest-compat"
	mailpitImage                          = "axllent/mailpit:v1.31"
	rabbitmqImage                         = "rabbitmq:4-management"
	redisImage                            = "redis:8-alpine"
	grafanaLgtmImage                      = "grafana/otel-lgtm:0.30.0"
)

// Versions holds all third-party dependency and plugin versions used in generated projects.
type Versions struct {
	JavaVersion                           int
	SpringBootVersion                     string
	GradleDependencyMgmtVersion           string
	GradleGitPropertiesVersion            string
	SpotlessGradlePluginVersion           string
	SpotlessMavenPluginVersion            string
	SpringModulithVersion                 string
	SpringCloudAWSVersion                 string
	SpringBootFlociTcVersion              string
	SpringdocOpenapiVersion               string
	TcMailpitVersion                      string
	FontAwesomeVersion                    string
	HtmxVersion                           string
	HtmxSpringBootThymeleafVersion        string
	DatasourceMicrometerSpringBootVersion string
	JjwtVersion                           string
	BootUIVersion                         string
	ArchunitVersion                       string
	TaikaiVersion                         string
	PalantirJavaFormatVersion             string
	PostgresImage                         string
	MysqlImage                            string
	MariadbImage                          string
	FlociImage                            string
	MailpitImage                          string
	RabbitMQImage                         string
	RedisImage                            string
	GrafanaLgtmImage                      string
}

func defaultVersions() Versions {
	return Versions{
		JavaVersion:                           javaVersion,
		SpringBootVersion:                     springBootVersion,
		GradleDependencyMgmtVersion:           gradleDependencyMgmtVersion,
		GradleGitPropertiesVersion:            gradleGitPropertiesVersion,
		SpotlessGradlePluginVersion:           spotlessGradlePluginVersion,
		SpotlessMavenPluginVersion:            spotlessMavenPluginVersion,
		SpringModulithVersion:                 springModulithVersion,
		SpringCloudAWSVersion:                 springCloudAWSVersion,
		SpringBootFlociTcVersion:              springBootFlociTcVersion,
		SpringdocOpenapiVersion:               springdocOpenapiVersion,
		TcMailpitVersion:                      tcMailpitVersion,
		FontAwesomeVersion:                    fontAwesomeVersion,
		HtmxVersion:                           htmxVersion,
		HtmxSpringBootThymeleafVersion:        htmxSpringBootThymeleafVersion,
		DatasourceMicrometerSpringBootVersion: datasourceMicrometerSpringBootVersion,
		BootUIVersion:                         bootUIVersion,
		JjwtVersion:                           jjwtVersion,
		ArchunitVersion:                       archunitVersion,
		TaikaiVersion:                         taikaiVersion,
		PalantirJavaFormatVersion:             palantirJavaFormatVersion,
		PostgresImage:                         postgresImage,
		MysqlImage:                            mysqlImage,
		MariadbImage:                          mariadbImage,
		FlociImage:                            flociImage,
		MailpitImage:                          mailpitImage,
		RabbitMQImage:                         rabbitmqImage,
		RedisImage:                            redisImage,
		GrafanaLgtmImage:                      grafanaLgtmImage,
	}
}
