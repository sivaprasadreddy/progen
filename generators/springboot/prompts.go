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
	var appType = RestApi

	appTypeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[AppType]().
				Title("Select App Type:").
				Options(
					huh.NewOption(RestApi.String(), RestApi).Selected(true),
					huh.NewOption(WebApp.String(), WebApp),
					huh.NewOption(SpringBootAngularFullStack.String(), SpringBootAngularFullStack),
				).Value(&appType),
		),
	)
	if err := runForm(appTypeForm); err != nil {
		return nil, err
	}

	answers := ProjectConfig{
		AppType:         appType,
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
	}

	otherFeatureOptions := []huh.Option[string]{
		huh.NewOption(FeatureSpringCloudAWSSupport, FeatureSpringCloudAWSSupport),
		huh.NewOption(FeatureEmailSupport, FeatureEmailSupport),
		huh.NewOption(FeatureRabbitMQSupport, FeatureRabbitMQSupport),
		huh.NewOption(FeatureRedisCachingSupport, FeatureRedisCachingSupport),
		huh.NewOption(FeatureOpenTelemetrySupport, FeatureOpenTelemetrySupport),
	}

	if answers.AppType == WebApp {
		otherFeatureOptions = append(otherFeatureOptions, huh.NewOption(FeatureHTMXSupport, FeatureHTMXSupport))
	}

	otherFeaturesSelect := huh.NewMultiSelect[string]().
		Title("Select Features:").
		Options(otherFeatureOptions...).
		Height(len(otherFeatureOptions) + 3).
		Value(&features)

	inputs = append(inputs, otherFeaturesSelect)

	form := huh.NewForm(huh.NewGroup(inputs...))
	if err := runForm(form); err != nil {
		return &answers, err
	}
	updateFeatureFlags(&answers, features)
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
	pc.ThymeleafSupport = pc.AppType == WebApp
	pc.HTMXSupport = isEnabled(features, FeatureHTMXSupport)
	pc.EmailSupport = isEnabled(features, FeatureEmailSupport)
	pc.RabbitMQSupport = isEnabled(features, FeatureRabbitMQSupport)
	pc.RedisCachingSupport = isEnabled(features, FeatureRedisCachingSupport)
	pc.OpenTelemetrySupport = isEnabled(features, FeatureOpenTelemetrySupport)
}

func isEnabled(features []string, feature string) bool {
	return features != nil && slices.Contains(features, feature)
}
