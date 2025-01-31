package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	minimalgo "github.com/sivaprasadreddy/progen/generators/minimal-go"

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

	generatorType := GeneratorType{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose application type:").
				Options(
					huh.NewOption(appTypeSpringBoot, appTypeSpringBoot).Selected(true),
					huh.NewOption(appTypeMinimalJava, appTypeMinimalJava),
					huh.NewOption(appTypeMinimalGo, appTypeMinimalGo),
				).Value(&generatorType.AppType),
		),
	)

	err := form.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	} else if err != nil {
		return generatorType, err
	}

	return generatorType, nil
}
