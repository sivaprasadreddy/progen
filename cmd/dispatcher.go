package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/sivaprasadreddy/progen/generators/minimal-java"
	"os"
	"strings"
)

const appTypeMinimalJava = "Minimal Java"
const appTypeSpringBoot = "Spring Boot"

func invokeGenerator(genType string) error {
	if genType == "" {
		answers, err := getAppTypeAnswers()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		genType = answers.AppType
	}

	if strings.EqualFold(genType, appTypeMinimalJava) {
		minimal_java.Run()
	} else if strings.EqualFold(genType, appTypeSpringBoot) {
		fmt.Println("Spring Boot - Work In Progress")
	} else {
		fmt.Println("Unknown generator type")
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
				Options: []string{appTypeMinimalJava, appTypeSpringBoot},
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
