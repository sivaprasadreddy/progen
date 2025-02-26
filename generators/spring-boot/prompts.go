package spring_boot

import (
	"errors"
	"os"

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

const appTypeRestApi = "REST API"
const appTypeWebApp = "Web App"

func getProjectConfigAnswers() (*ProjectConfig, error) {
	var appType = appTypeRestApi

	appTypeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select App Type:").
				Options(
					huh.NewOption(appTypeRestApi, appTypeRestApi).Selected(true),
					huh.NewOption(appTypeWebApp, appTypeWebApp),
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
		BuildTool:       "Maven",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Flyway",
		Features:        []string{},
	}

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
				huh.NewOption("Maven", "Maven").Selected(true),
				huh.NewOption("Gradle", "Gradle"),
			).Value(&answers.BuildTool),

		huh.NewSelect[string]().
			Title("Select Database:").
			Options(
				huh.NewOption("PostgreSQL", "PostgreSQL").Selected(true),
				huh.NewOption("MySQL", "MySQL"),
				huh.NewOption("MariaDB", "MariaDB"),
			).Value(&answers.DbType),

		huh.NewSelect[string]().
			Title("Select Database Migration Tool:").
			Options(
				huh.NewOption("Flyway", "Flyway").Selected(true),
				huh.NewOption("Liquibase", "Liquibase"),
			).Value(&answers.DbMigrationTool),
	}

	otherFeatureOptions := []huh.Option[string]{
		huh.NewOption("Docker Compose", "Docker Compose").Selected(true),
		huh.NewOption("Spring Modulith", "Spring Modulith"),
		huh.NewOption("Spring Cloud AWS", "Spring Cloud AWS"),
	}

	if answers.AppType == appTypeWebApp {
		otherFeatureOptions = append(otherFeatureOptions,
			huh.NewOption("Security", "Security"),
			huh.NewOption("Thymeleaf", "Thymeleaf").Selected(true),
			huh.NewOption("HTMX", "HTMX"))
	}

	if answers.AppType == appTypeRestApi {
		otherFeatureOptions = append(otherFeatureOptions,
			huh.NewOption("JWT Security", "JWT Security"),
		)
	}

	otherFeaturesSelect := huh.NewMultiSelect[string]().
		Title("Select Other Features:").
		Options(otherFeatureOptions...).
		Value(&answers.Features)

	inputs = append(inputs, otherFeaturesSelect)

	form := huh.NewForm(huh.NewGroup(inputs...))

	err = form.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	} else if err != nil {
		return &answers, err
	}
	return &answers, nil
}
