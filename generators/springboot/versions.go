package springboot

const (
	javaVersion                    = 25
	springBootVersion              = "4.1.0"
	gradleDependencyMgmtVersion    = "1.1.7"
	gradleGitPropertiesVersion     = "4.0.1"
	spotlessGradlePluginVersion    = "8.6.0"
	spotlessMavenPluginVersion     = "3.6.0"
	springModulithVersion          = "2.1.0"
	springCloudAWSVersion          = "4.0.2"
	springdocOpenapiVersion        = "3.0.3"
	bootstrapVersion               = "5.3.8"
	jqueryVersion                  = "3.7.1"
	fontAwesomeVersion             = "7.2.0"
	htmxVersion                    = "2.0.10"
	htmxSpringBootThymeleafVersion = "5.1.0"
	jjwtVersion                    = "0.13.0"
	archunitVersion                = "1.4.2"
	palantirJavaFormatVersion      = "2.93.0"
)

// Versions holds all third-party dependency and plugin versions used in generated projects.
type Versions struct {
	JavaVersion                    int
	SpringBootVersion              string
	GradleDependencyMgmtVersion    string
	GradleGitPropertiesVersion     string
	SpotlessGradlePluginVersion    string
	SpotlessMavenPluginVersion     string
	SpringModulithVersion          string
	SpringCloudAWSVersion          string
	SpringdocOpenapiVersion        string
	BootstrapVersion               string
	JqueryVersion                  string
	FontAwesomeVersion             string
	HtmxVersion                    string
	HtmxSpringBootThymeleafVersion string
	JjwtVersion                    string
	ArchunitVersion                string
	PalantirJavaFormatVersion      string
}

func defaultVersions() Versions {
	return Versions{
		JavaVersion:                    javaVersion,
		SpringBootVersion:              springBootVersion,
		GradleDependencyMgmtVersion:    gradleDependencyMgmtVersion,
		GradleGitPropertiesVersion:     gradleGitPropertiesVersion,
		SpotlessGradlePluginVersion:    spotlessGradlePluginVersion,
		SpotlessMavenPluginVersion:     spotlessMavenPluginVersion,
		SpringModulithVersion:          springModulithVersion,
		SpringCloudAWSVersion:          springCloudAWSVersion,
		SpringdocOpenapiVersion:        springdocOpenapiVersion,
		BootstrapVersion:               bootstrapVersion,
		JqueryVersion:                  jqueryVersion,
		FontAwesomeVersion:             fontAwesomeVersion,
		HtmxVersion:                    htmxVersion,
		HtmxSpringBootThymeleafVersion: htmxSpringBootThymeleafVersion,
		JjwtVersion:                    jjwtVersion,
		ArchunitVersion:                archunitVersion,
		PalantirJavaFormatVersion:      palantirJavaFormatVersion,
	}
}
