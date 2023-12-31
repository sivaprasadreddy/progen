package minimal_java

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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
	var questions = []*survey.Question{
		{
			Name: "AppName",
			Prompt: &survey.Input{
				Message: "What is the base name of your application?",
				Help:    "Name of your application",
				Default: "myapp",
			},
			Validate:  helpers.ValidateApplicationName,
			Transform: helpers.TrimString,
		},
		{
			Name: "GroupID",
			Prompt: &survey.Input{
				Message: "What is your application groupId?",
				Default: "com.mycompany",
			},
			Validate: survey.Required,
		},
		{
			Name: "ArtifactID",
			Prompt: &survey.Input{
				Message: "What is your application artifactId?",
				Default: "myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "AppVersion",
			Prompt: &survey.Input{
				Message: "What is your application version?",
				Default: "1.0.0-SNAPSHOT",
			},
			Validate: survey.Required,
		},
		{
			Name: "BasePackage",
			Prompt: &survey.Input{
				Message: "What is your application base package?",
				Help:    "Base package",
				Default: "com.mycompany.myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "BuildTool",
			Prompt: &survey.Select{
				Message: "Which build tool would you like to use?",
				Options: []string{"Maven", "Gradle"},
				Default: "Maven",
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
