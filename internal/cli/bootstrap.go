package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
)

func invokeGenerator(ctx context.Context, configFile string) error {
	var projectConfig springboot.ProjectConfig
	if strings.TrimSpace(configFile) != "" {
		pc, warnings, err := springboot.LoadConfigFile(configFile)
		if err != nil {
			return err
		}
		for _, w := range warnings {
			fmt.Println("WARNING: " + w)
		}
		projectConfig = pc
	} else {
		pc, err := getAnswers()
		if err != nil {
			return err
		}
		projectConfig = pc
	}
	// The CLI has always replaced an existing project directory; interactive mode
	// already rejects names of existing directories while prompting.
	result, err := springboot.Generate(ctx, projectConfig, springboot.GenerateOptions{
		Overwrite: true,
		Output:    os.Stdout,
	})
	if err != nil {
		return err
	}
	for _, w := range result.Warnings {
		fmt.Println("WARNING: " + w)
	}
	return nil
}
