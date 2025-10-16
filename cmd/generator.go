package cmd

import (
	"github.com/sivaprasadreddy/progen/generators/springboot"
)

func invokeGenerator() error {
	return springboot.Run()
}
