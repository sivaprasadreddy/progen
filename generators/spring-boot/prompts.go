package spring_boot

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
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
			Name: "ApplicationName",
			Prompt: &survey.Input{
				Message: "What is ApplicationName?",
				Help:    "Name of your application",
				Default: "myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "GroupID",
			Prompt: &survey.Input{
				Message: "What is GroupID?",
				Default: "com.mycompany",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) < 1 {
					return errors.New("invalid groupId")
				}
				return nil
			},
		},
		{
			Name: "ArtifactID",
			Prompt: &survey.Input{
				Message: "What is ArtifactID?",
				Default: "myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "ApplicationVersion",
			Prompt: &survey.Input{
				Message: "What is Application Version?",
				Default: "1.0.0-SNAPSHOT",
			},
			Validate: survey.Required,
		},
		{
			Name: "BasePackage",
			Prompt: &survey.Input{
				Message: "What is base package?",
				Help:    "Base package name",
				Default: "com.mycompany.myapp",
			},
			Validate: survey.Required,
		},
		{
			Name: "BuildTool",
			Prompt: &survey.Select{
				Message: "Choose Build Tool:",
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
