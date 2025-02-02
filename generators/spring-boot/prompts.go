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
	return answers, nil
}

func getProjectConfigAnswers() (ProjectConfig, error) {
	answers := ProjectConfig{
		AppName:         "myapp",
		GroupID:         "com.mycompany",
		ArtifactID:      "myapp",
		AppVersion:      "1.0.0-SNAPSHOT",
		BasePackage:     "com.mycompany.myapp",
		BuildTool:       "Maven",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Flyway",
		Features:        []string{"Spring Modulith", "Thymeleaf", "Security"},
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the base name of your application?").
				Description("Name of your application").
				Validate(func(str string) error {
					return helpers.ValidateApplicationName(str)
				}).
				Value(&answers.AppName),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What is your application groupId?").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("groupID is required")
					}
					return nil
				}).
				Value(&answers.GroupID),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What is your applciation artifactID").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("artifactID is required")
					}
					return nil
				}).Value(&answers.ArtifactID),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What is you application version?").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("app version is required")
					}
					return nil
				}).Value(&answers.AppVersion),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What is your application base package?").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("base package is required")
					}
					return nil
				}).Value(&answers.BasePackage),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which build tool would you like to use?").
				Options(
					huh.NewOption("Maven", "Maven").Selected(true),
					huh.NewOption("Gradle", "Gradle"),
				).Value(&answers.BuildTool),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which database would you like to use?").
				Options(
					huh.NewOption("PostgreSQL", "PostgreSQL").Selected(true),
					huh.NewOption("MySQL", "MySQL"),
					huh.NewOption("MariaDB", "MariaDB"),
				).Value(&answers.DbType),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which database migration tool would you like to use?").
				Options(
					huh.NewOption("Flyway", "Flyway"),
					huh.NewOption("Liquibase", "Liquibase"),
				).Value(&answers.DbMigrationTool),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("What features you would like to add?").
				Options(
					huh.NewOption("Spring Modulith", "Spring Modulith").Selected(true),
					huh.NewOption("Spring Cloud AWS", "Spring Cloud AWS"),
					huh.NewOption("Thymeleaf", "Thymeleaf").Selected(true),
					huh.NewOption("HTMX", "HTMX"),
					huh.NewOption("Security", "Security").Selected(true),
					huh.NewOption("JWT Security", "JWT Security"),
				).Value(&answers.Features),
		),
	)
	err := form.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	} else if err != nil {
		return answers, err
	}
	return answers, nil
}
