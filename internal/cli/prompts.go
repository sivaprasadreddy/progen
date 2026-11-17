package cli

import (
	"errors"
	"fmt"
	"os"

	"charm.land/huh/v2"
	sb "github.com/sivaprasadreddy/progen/internal/generators/springboot"
)

// getAnswers collects the project configuration through interactive terminal prompts.
// It returns huh.ErrUserAborted when the user cancels.
func getAnswers() (sb.ProjectConfig, error) {
	answers := sb.ProjectConfig{
		Frontend:        sb.FrontendNone,
		AppName:         "myapp",
		GroupID:         "com.mycompany",
		ArtifactID:      "myapp",
		AppVersion:      "1.0.0",
		BuildTool:       sb.Maven,
		PersistenceType: sb.SpringDataJPA,
		DbType:          sb.PostgreSQL,
		DbMigrationTool: sb.Flyway,
		TaskRunner:      sb.TaskRunnerMise,
	}
	choices := sb.Choices()

	identityForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Application Name (ex: boot-demo):").
				Validate(validateNewAppDir).
				Value(&answers.AppName),

			huh.NewInput().
				Title("GroupId (ex: com.mycompany):").
				Validate(required("GroupId is required")).
				Value(&answers.GroupID),

			huh.NewInput().
				Title("ArtifactId (ex: boot-demo):").
				Validate(required("ArtifactId is required")).
				Value(&answers.ArtifactID),
		),
	)
	if err := identityForm.Run(); err != nil {
		return sb.ProjectConfig{}, err
	}

	answers.BasePackage = sb.DeriveBasePackage(answers.GroupID, answers.ArtifactID)

	var features []string

	inputs := []huh.Field{
		huh.NewInput().
			Title("Package Name (ex: com.mycompany.myapp):").
			Validate(required("package name is required")).
			Value(&answers.BasePackage),

		huh.NewInput().
			Title("Application Version (ex: 1.0.0):").
			Validate(required("application version is required")).
			Value(&answers.AppVersion),

		huh.NewSelect[sb.BuildTool]().
			Title("Build Tool:").
			Options(toOptions(choices.BuildTools, answers.BuildTool)...).
			Value(&answers.BuildTool),

		huh.NewSelect[sb.PersistenceType]().
			Title("Persistence:").
			Options(toOptions(choices.PersistenceTypes, answers.PersistenceType)...).
			Value(&answers.PersistenceType),

		huh.NewSelect[sb.DatabaseType]().
			Title("Database:").
			Options(toOptions(choices.DatabaseTypes, answers.DbType)...).
			Value(&answers.DbType),

		huh.NewSelect[sb.DbMigrationTool]().
			Title("Database Migration Tool:").
			Options(toOptions(choices.DbMigrationTools, answers.DbMigrationTool)...).
			Value(&answers.DbMigrationTool),

		huh.NewSelect[sb.TaskRunner]().
			Title("Task Runner:").
			Options(toOptions(choices.TaskRunners, answers.TaskRunner)...).
			Value(&answers.TaskRunner),

		huh.NewSelect[sb.Frontend]().
			Title("Select Frontend:").
			Options(toOptions(choices.Frontends, answers.Frontend)...).
			Value(&answers.Frontend),
	}

	// Options depend on the selected frontend (HTMX only applies to server-rendered views),
	// so they are recomputed whenever answers.Frontend changes.
	otherFeaturesSelect := huh.NewMultiSelect[string]().
		Title("Select Features:").
		OptionsFunc(func() []huh.Option[string] {
			return huh.NewOptions(sb.AvailableFeatures(answers.Frontend)...)
		}, &answers.Frontend).
		Height(len(sb.AvailableFeatures(sb.FrontendNone)) + 1).
		Value(&features)

	var aiModels []sb.AIModelProvider
	aiModelsSelect := huh.NewMultiSelect[sb.AIModelProvider]().
		Title("Spring AI Models (optional):").
		Description("Select the chat model providers to integrate using Spring AI.").
		Options(toOptions[sb.AIModelProvider](choices.SpringAIModels, "")...).
		Height(len(choices.SpringAIModels) + 1).
		Value(&aiModels)

	inputs = append(inputs, otherFeaturesSelect, aiModelsSelect)

	form := huh.NewForm(huh.NewGroup(inputs...))
	if err := form.Run(); err != nil {
		return answers, err
	}
	answers.SetFeatures(features)
	answers.SpringAIModels = sb.NormalizeAIModels(aiModels)
	return answers, nil
}

// toOptions converts library choices into huh options, pre-selecting the current value.
func toOptions[T ~string](choices []sb.Choice, current T) []huh.Option[T] {
	opts := make([]huh.Option[T], 0, len(choices))
	for _, c := range choices {
		opts = append(opts, huh.NewOption(c.Label, T(c.Value)).Selected(T(c.Value) == current))
	}
	return opts
}

func required(msg string) func(string) error {
	return func(str string) error {
		if str == "" {
			return errors.New(msg)
		}
		return nil
	}
}

// validateNewAppDir checks the application name and that no directory with that
// name exists in the working directory, where the CLI generates the project.
func validateNewAppDir(name string) error {
	if err := sb.ValidateAppName(name); err != nil {
		return err
	}
	if stat, err := os.Stat(name); err == nil && stat.IsDir() {
		return fmt.Errorf("a directory with name '%s' already exists", name)
	}
	return nil
}
