package minimal_go

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func getAnswers() (ProjectConfig, error) {
	answers, err := getProjectConfigAnswers()
	if err != nil {
		return ProjectConfig{}, err
	}
	return answers, nil
}

func getProjectConfigAnswers() (ProjectConfig, error) {
	var questions = []*survey.Question{
		{
			Name: "AppName",
			Prompt: &survey.Input{
				Message: "What is ApplicationName?",
				Help:    "Name of your application",
				Default: "myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "ModulePath",
			Prompt: &survey.Input{
				Message: "What is Module Path?",
				Help:    "Module path ex: github.com/username/modulename",
				Default: "github.com/username/myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "RoutingLibrary",
			Prompt: &survey.Select{
				Message: "Choose Routing Library:",
				Options: []string{"Gin", "Chi"},
				Default: "Gin",
			},
		},
	}
	answers := ProjectConfig{}
	err := survey.Ask(questions, &answers)
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(0)
	} else if err != nil {
		return answers, err
	}
	return answers, nil
}
