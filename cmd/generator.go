package cmd

import (
	"github.com/sivaprasadreddy/progen/generators/springboot"
)

func invokeGenerator(configFile string) error {
	return springboot.Run(configFile)
}
