package cmd

import (
	"errors"
	"os"
	"strings"

	minimalgo "github.com/sivaprasadreddy/progen/generators/minimal-go"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	minimaljava "github.com/sivaprasadreddy/progen/generators/minimal-java"
	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
)

const appTypeMinimalJava = "Minimal Java"
const appTypeSpringBoot = "Spring Boot"
const appTypeMinimalGo = "Minimal Go"

func invokeGenerator() error {
	answers, err := getAppTypeAnswers()
	if err != nil {
		return err
	}
	if strings.EqualFold(answers.AppType, appTypeMinimalJava) {
		minimaljava.Run()
	} else if strings.EqualFold(answers.AppType, appTypeSpringBoot) {
		springboot.Run()
	} else if strings.EqualFold(answers.AppType, appTypeMinimalGo) {
		minimalgo.Run()
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
				Options: []string{appTypeSpringBoot, appTypeMinimalJava, appTypeMinimalGo},
				Default: appTypeSpringBoot,
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
