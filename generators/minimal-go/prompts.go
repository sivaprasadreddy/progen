package minimal_go

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
		AppName:        "myapp",
		ModulePath:     "github.com/username/myapp",
		RoutingLibrary: "Gin",
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
		), huh.NewGroup(
			huh.NewInput().
				Title("What is your application module path?").
				Description("Module path ex: github.com/username/modulename").
				Value(&answers.ModulePath),
		), huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose Routing Library").
				Options(
					huh.NewOption("GIN", "Gin"),
					huh.NewOption("CHI", "Chi"),
				).Value(&answers.RoutingLibrary),
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
