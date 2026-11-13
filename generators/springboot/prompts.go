package springboot

import (
	"errors"
	"os"
	"slices"
	"strings"

	"charm.land/huh/v2"
	"github.com/sivaprasadreddy/progen/generators/helpers"
)

func GetAnswers() (ProjectConfig, error) {
	answers, err := getProjectConfigAnswers()
	if err != nil {
		return ProjectConfig{}, err
	}
	return *answers, nil
}

func getProjectConfigAnswers() (*ProjectConfig, error) {
	answers := ProjectConfig{
		Frontend:        FrontendNone,
		AppName:         "myapp",
		GroupID:         "com.mycompany",
		ArtifactID:      "myapp",
		AppVersion:      "1.0.0",
		BuildTool:       Maven,
		PersistenceType: SpringDataJPA,
		DbType:          PostgreSQL,
		DbMigrationTool: Flyway,
	}

	identityForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Application Name (ex: boot-demo):").
				Validate(func(str string) error {
					return helpers.ValidateApplicationName(str)
				}).Value(&answers.AppName),

			huh.NewInput().
				Title("GroupId (ex: com.mycompany):").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("GroupId is required")
					}
					return nil
				}).
				Value(&answers.GroupID),

			huh.NewInput().
				Title("ArtifactId (ex: boot-demo):").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("ArtifactId is required")
					}
					return nil
				}).Value(&answers.ArtifactID),
		),
	)
	if err := runForm(identityForm); err != nil {
		return nil, err
	}

	answers.BasePackage = deriveBasePackage(answers.GroupID, answers.ArtifactID)

	var features []string

	inputs := []huh.Field{
		huh.NewInput().
			Title("Package Name (ex: com.mycompany.myapp):").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("package name is required")
				}
				return nil
			}).Value(&answers.BasePackage),

		huh.NewInput().
			Title("Application Version (ex: 1.0.0):").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("application version is required")
				}
				return nil
			}).Value(&answers.AppVersion),

		huh.NewSelect[BuildTool]().
			Title("Build Tool:").
			Options(
				huh.NewOption(Maven.String(), Maven).Selected(true),
				huh.NewOption(Gradle.String(), Gradle),
			).Value(&answers.BuildTool),

		huh.NewSelect[PersistenceType]().
			Title("Persistence:").
			Options(
				huh.NewOption(SpringDataJPA.String(), SpringDataJPA).Selected(true),
				huh.NewOption(SpringJdbcClient.String(), SpringJdbcClient),
				huh.NewOption(SpringJOOQ.String(), SpringJOOQ),
			).Value(&answers.PersistenceType),

		huh.NewSelect[DatabaseType]().
			Title("Database:").
			Options(
				huh.NewOption(PostgreSQL.String(), PostgreSQL).Selected(true),
				huh.NewOption(MySQL.String(), MySQL),
				huh.NewOption(MariaDB.String(), MariaDB),
			).Value(&answers.DbType),

		huh.NewSelect[DbMigrationTool]().
			Title("Database Migration Tool:").
			Options(
				huh.NewOption(Flyway.String(), Flyway).Selected(true),
				huh.NewOption(Liquibase.String(), Liquibase),
			).Value(&answers.DbMigrationTool),

		huh.NewSelect[Frontend]().
			Title("Select Frontend:").
			Options(
				huh.NewOption("None (REST API)", FrontendNone).Selected(true),
				huh.NewOption(FrontendThymeleaf.String(), FrontendThymeleaf),
				huh.NewOption(FrontendJTE.String(), FrontendJTE),
				huh.NewOption(FrontendAngular.String(), FrontendAngular),
			).Value(&answers.Frontend),
	}

	baseFeatureOptions := []huh.Option[string]{
		huh.NewOption(FeatureSpringCloudAWSSupport, FeatureSpringCloudAWSSupport),
		huh.NewOption(FeatureEmailSupport, FeatureEmailSupport),
		huh.NewOption(FeatureRabbitMQSupport, FeatureRabbitMQSupport),
		huh.NewOption(FeatureRedisCachingSupport, FeatureRedisCachingSupport),
		huh.NewOption(FeatureOpenTelemetrySupport, FeatureOpenTelemetrySupport),
		huh.NewOption(FeatureK8sSupport, FeatureK8sSupport),
	}

	// Options depend on the selected frontend (HTMX only applies to server-rendered views),
	// so they are recomputed whenever answers.Frontend changes.
	otherFeaturesSelect := huh.NewMultiSelect[string]().
		Title("Select Features:").
		OptionsFunc(func() []huh.Option[string] {
			opts := slices.Clone(baseFeatureOptions)
			if answers.ServerRenderedViews() {
				opts = append(opts, huh.NewOption(FeatureHTMXSupport, FeatureHTMXSupport))
			}
			return opts
		}, &answers.Frontend).
		Height(len(baseFeatureOptions) + 4).
		Value(&features)

	var aiModels []AIModelProvider
	aiModelOptions := make([]huh.Option[AIModelProvider], 0, len(allAIModelProviders))
	for _, p := range allAIModelProviders {
		aiModelOptions = append(aiModelOptions, huh.NewOption(p.Label(), p))
	}
	aiModelsSelect := huh.NewMultiSelect[AIModelProvider]().
		Title("Spring AI Models (optional):").
		Description("Select the chat model providers to integrate using Spring AI.").
		Options(aiModelOptions...).
		Height(len(aiModelOptions) + 4).
		Value(&aiModels)

	inputs = append(inputs, otherFeaturesSelect, aiModelsSelect)

	form := huh.NewForm(huh.NewGroup(inputs...))
	if err := runForm(form); err != nil {
		return &answers, err
	}
	updateFeatureFlags(&answers, features)
	answers.SpringAIModels = sanitizeAIModels(aiModels)
	return &answers, nil
}

func runForm(form *huh.Form) error {
	err := form.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	}
	return err
}

func deriveBasePackage(groupID, artifactID string) string {
	groupID = strings.TrimSpace(groupID)
	segment := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(artifactID), "-", ""))
	if groupID == "" {
		return segment
	}
	if segment == "" {
		return groupID
	}
	return groupID + "." + segment
}

func updateFeatureFlags(pc *ProjectConfig, features []string) {
	pc.SpringCloudAWSSupport = isEnabled(features, FeatureSpringCloudAWSSupport)
	pc.HTMXSupport = pc.ServerRenderedViews() && isEnabled(features, FeatureHTMXSupport)
	pc.EmailSupport = isEnabled(features, FeatureEmailSupport)
	pc.RabbitMQSupport = isEnabled(features, FeatureRabbitMQSupport)
	pc.RedisCachingSupport = isEnabled(features, FeatureRedisCachingSupport)
	pc.OpenTelemetrySupport = isEnabled(features, FeatureOpenTelemetrySupport)
	pc.K8sSupport = isEnabled(features, FeatureK8sSupport)
}

func isEnabled(features []string, feature string) bool {
	return features != nil && slices.Contains(features, feature)
}
