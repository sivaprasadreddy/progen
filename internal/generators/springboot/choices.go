package springboot

import (
	"slices"
)

// Choice is a selectable value for a ProjectConfig field, with a human-readable label.
// Frontends (CLI prompts, desktop app) use these to render their inputs.
type Choice struct {
	Value string
	Label string
}

// ConfigChoices lists the supported values for each selectable ProjectConfig field.
// The first entry of every list is the default.
type ConfigChoices struct {
	BuildTools       []Choice
	PersistenceTypes []Choice
	DatabaseTypes    []Choice
	DbMigrationTools []Choice
	TaskRunners      []Choice
	Frontends        []Choice
	SpringAIModels   []Choice
}

// Choices returns the supported values for each selectable ProjectConfig field.
func Choices() ConfigChoices {
	aiModels := make([]Choice, 0, len(allAIModelProviders))
	for _, p := range allAIModelProviders {
		aiModels = append(aiModels, Choice{p.String(), p.Label()})
	}
	return ConfigChoices{
		BuildTools: []Choice{
			{Maven.String(), Maven.String()},
			{Gradle.String(), Gradle.String()},
		},
		PersistenceTypes: []Choice{
			{SpringDataJPA.String(), SpringDataJPA.String()},
			{SpringJdbcClient.String(), SpringJdbcClient.String()},
			{SpringJOOQ.String(), SpringJOOQ.String()},
		},
		DatabaseTypes: []Choice{
			{PostgreSQL.String(), PostgreSQL.String()},
			{MySQL.String(), MySQL.String()},
			{MariaDB.String(), MariaDB.String()},
		},
		DbMigrationTools: []Choice{
			{Flyway.String(), Flyway.String()},
			{Liquibase.String(), Liquibase.String()},
		},
		TaskRunners: []Choice{
			{TaskRunnerMise.String(), "Mise"},
			{TaskRunnerTaskfile.String(), "Taskfile (with SDKMAN)"},
		},
		Frontends: []Choice{
			{FrontendNone.String(), "None (REST API)"},
			{FrontendThymeleaf.String(), FrontendThymeleaf.String()},
			{FrontendJTE.String(), FrontendJTE.String()},
			{FrontendAngular.String(), FrontendAngular.String()},
			{FrontendReact.String(), FrontendReact.String()},
		},
		SpringAIModels: aiModels,
	}
}

// baseFeatures are the optional features available for every frontend.
var baseFeatures = []string{
	FeatureSpringCloudAWSSupport,
	FeatureEmailSupport,
	FeatureRabbitMQSupport,
	FeatureRedisCachingSupport,
	FeatureOpenTelemetrySupport,
	FeatureAxelixSupport,
	FeatureK8sSupport,
}

// AvailableFeatures returns the optional features that can be enabled for the given frontend.
// HTMX is only offered for server-rendered views.
func AvailableFeatures(frontend Frontend) []string {
	features := slices.Clone(baseFeatures)
	if (ProjectConfig{Frontend: frontend}).ServerRenderedViews() {
		features = append(features, FeatureHTMXSupport)
	}
	return features
}

// SetFeatures enables exactly the given optional features (see AvailableFeatures)
// and disables the rest. Features not applicable to the selected frontend are ignored,
// so set Frontend before calling it.
func (pc *ProjectConfig) SetFeatures(features []string) {
	pc.AxelixSupport = slices.Contains(features, FeatureAxelixSupport)
	pc.SpringCloudAWSSupport = slices.Contains(features, FeatureSpringCloudAWSSupport)
	pc.HTMXSupport = pc.ServerRenderedViews() && slices.Contains(features, FeatureHTMXSupport)
	pc.EmailSupport = slices.Contains(features, FeatureEmailSupport)
	pc.RabbitMQSupport = slices.Contains(features, FeatureRabbitMQSupport)
	pc.RedisCachingSupport = slices.Contains(features, FeatureRedisCachingSupport)
	pc.OpenTelemetrySupport = slices.Contains(features, FeatureOpenTelemetrySupport)
	pc.K8sSupport = slices.Contains(features, FeatureK8sSupport)
}

// EnabledFeatures returns the optional features currently enabled, in AvailableFeatures order.
func (pc ProjectConfig) EnabledFeatures() []string {
	enabled := map[string]bool{
		FeatureAxelixSupport:         pc.AxelixSupport,
		FeatureSpringCloudAWSSupport: pc.SpringCloudAWSSupport,
		FeatureEmailSupport:          pc.EmailSupport,
		FeatureRabbitMQSupport:       pc.RabbitMQSupport,
		FeatureRedisCachingSupport:   pc.RedisCachingSupport,
		FeatureOpenTelemetrySupport:  pc.OpenTelemetrySupport,
		FeatureK8sSupport:            pc.K8sSupport,
		FeatureHTMXSupport:           pc.HTMXSupport,
	}
	var result []string
	for _, f := range AvailableFeatures(pc.Frontend) {
		if enabled[f] {
			result = append(result, f)
		}
	}
	return result
}
