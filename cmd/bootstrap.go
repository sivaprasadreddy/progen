package cmd

import (
	"strings"

	"github.com/sivaprasadreddy/progen/generators/helpers"
	"github.com/sivaprasadreddy/progen/generators/springboot"
)

func invokeGenerator(configFile string) error {
	var projectConfig springboot.ProjectConfig
	var err error
	if strings.TrimSpace(configFile) == "" {
		projectConfig, err = springboot.GetAnswers()
		helpers.FatalIfErr(err)
		err = springboot.GenerateProject(projectConfig)
	} else {
		err = springboot.GenerateProjectFromConfigFile(configFile)
	}
	helpers.FatalIfErr(err)
	return err
}
