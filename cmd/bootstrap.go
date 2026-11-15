package cmd

import (
	"strings"

	"github.com/sivaprasadreddy/progen/generators/springboot"
)

func invokeGenerator(configFile string) error {
	if strings.TrimSpace(configFile) != "" {
		return springboot.GenerateProjectFromConfigFile(configFile)
	}
	projectConfig, err := springboot.GetAnswers()
	if err != nil {
		return err
	}
	return springboot.GenerateProject(projectConfig)
}
