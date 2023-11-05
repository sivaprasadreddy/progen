package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	minimal_java "github.com/sivaprasadreddy/progen/generators/minimal-java"
	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
)

const minimalJava = "Minimal Java"
const springBoot = "Spring Boot"

func invokeGenerator() error {
	answers, err := getAppTypeAnswers()
	if err != nil {
		return err
	}
	if strings.EqualFold(answers.AppType, minimalJava) {
		minimal_java.Run()
	} else if strings.EqualFold(answers.AppType, springBoot) {
		springboot.Run()
	} else {
		return errors.New("unknown generator type")
	}
	return nil
}

type GeneratorType struct {
	AppType string
}

func getAppTypeAnswers() (GeneratorType, error) {
	var answers = []*survey.Question{
		{
			Name: "AppType",
			Prompt: &survey.Select{
				Message: "Choose application type:",
				Options: []string{minimalJava, springBoot},
				Default: "Minimal Java",
			},
		},
	}
	generatorType := GeneratorType{}
	err := survey.Ask(answers, &generatorType)
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(0)
	} else if err != nil {
		return generatorType, err
	}
	return generatorType, nil
}
