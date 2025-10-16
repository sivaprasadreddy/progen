package spring_boot

import (
	"errors"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/sivaprasadreddy/progen/generators/helpers"
)

func getAnswers() (ProjectConfig, error) {
	answers, err := getProjectConfigAnswers()
	if err != nil {
		return ProjectConfig{}, err
	}
	return *answers, nil
}

func getProjectConfigAnswers() (*ProjectConfig, error) {
	var appType = AppTypeRestApi

	appTypeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select App Type:").
				Options(
					huh.NewOption(AppTypeRestApi, AppTypeRestApi).Selected(true),
					huh.NewOption(AppTypeWebApp, AppTypeWebApp),
				).Value(&appType),
		),
	)
	err := appTypeForm.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	} else if err != nil {
		return nil, err
	}

	answers := ProjectConfig{
		AppType:         appType,
		AppName:         "myapp",
		GroupID:         "com.mycompany",
		ArtifactID:      "myapp",
		AppVersion:      "1.0.0",
		BasePackage:     "com.mycompany.myapp",
		BuildTool:       BuildToolMaven,
		DbType:          DbPostgreSQL,
		DbMigrationTool: DbMigrationToolFlyway,
	}

	var features []string

	inputs := []huh.Field{
		huh.NewInput().
			Title("Enter Application Name:").
			Description("Ex: spring-boot-demo").
			Validate(func(str string) error {
				return helpers.ValidateApplicationName(str)
			}).Value(&answers.AppName),

		huh.NewInput().
			Title("Enter GroupId:").
			Description("Ex: com.mycompany").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("GroupId is required")
				}
				return nil
			}).
			Value(&answers.GroupID),

		huh.NewInput().
			Title("Enter ArtifactId:").
			Description("Ex: spring-boot-demo").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("ArtifactId is required")
				}
				return nil
			}).Value(&answers.ArtifactID),

		huh.NewInput().
			Title("Enter Application Version:").
			Description("Ex: 1.0.0").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("Application version is required")
				}
				return nil
			}).Value(&answers.AppVersion),

		huh.NewInput().
			Title("Enter Package Name:").
			Description("Ex: com.mycompany.myapp").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("Package name is required")
				}
				return nil
			}).Value(&answers.BasePackage),

		huh.NewSelect[string]().
			Title("Select Build Tool:").
			Options(
				huh.NewOption(BuildToolMaven, BuildToolMaven).Selected(true),
				huh.NewOption(BuildToolGradle, BuildToolGradle),
			).Value(&answers.BuildTool),

		huh.NewSelect[string]().
			Title("Select Database:").
			Options(
				huh.NewOption(DbPostgreSQL, DbPostgreSQL).Selected(true),
				huh.NewOption(DbMySQL, DbMySQL),
				huh.NewOption(DbMariaDB, DbMariaDB),
			).Value(&answers.DbType),

		huh.NewSelect[string]().
			Title("Select Database Migration Tool:").
			Options(
				huh.NewOption(DbMigrationToolFlyway, DbMigrationToolFlyway).Selected(true),
				huh.NewOption(DbMigrationToolLiquibase, DbMigrationToolLiquibase),
			).Value(&answers.DbMigrationTool),
	}

	otherFeatureOptions := []huh.Option[string]{
		huh.NewOption(FeatureDockerComposeSupport, FeatureDockerComposeSupport).Selected(true),
		huh.NewOption(FeatureSpringModulithSupport, FeatureSpringModulithSupport),
		huh.NewOption(FeatureSpringCloudAWSSupport, FeatureSpringCloudAWSSupport),
	}

	if answers.AppType == AppTypeWebApp {
		otherFeatureOptions = append(otherFeatureOptions,
			huh.NewOption(FeatureSecuritySupport, FeatureSecuritySupport),
			huh.NewOption(FeatureThymeleafSupport, FeatureThymeleafSupport).Selected(true),
			huh.NewOption(FeatureHTMXSupport, FeatureHTMXSupport))
	}

	if answers.AppType == AppTypeRestApi {
		otherFeatureOptions = append(otherFeatureOptions,
			huh.NewOption(FeatureJwtSecuritySupport, FeatureJwtSecuritySupport),
		)
	}

	otherFeaturesSelect := huh.NewMultiSelect[string]().
		Title("Select Other Features:").
		Options(otherFeatureOptions...).
		Value(&features)

	inputs = append(inputs, otherFeaturesSelect)

	form := huh.NewForm(huh.NewGroup(inputs...))

	err = form.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	} else if err != nil {
		return &answers, err
	}
	updateFeatureFlags(&answers, features)
	return &answers, nil
}

func updateFeatureFlags(pc *ProjectConfig, features []string) {
	pc.DockerComposeSupport = isEnabled(features, FeatureDockerComposeSupport)
	pc.SpringModulithSupport = isEnabled(features, FeatureSpringModulithSupport)
	pc.SpringCloudAWSSupport = isEnabled(features, FeatureSpringCloudAWSSupport)
	pc.ThymeleafSupport = isEnabled(features, FeatureThymeleafSupport)
	pc.HTMXSupport = isEnabled(features, FeatureHTMXSupport)
	pc.SecuritySupport = isEnabled(features, FeatureSecuritySupport)
	pc.JwtSecuritySupport = isEnabled(features, FeatureJwtSecuritySupport)
}

func isEnabled(features []string, feature string) bool {
	return features != nil && slices.Contains(features, feature)
}
