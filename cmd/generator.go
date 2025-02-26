package cmd

import (
	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
)

func invokeGenerator() error {
	springboot.Run()
	return nil
}
