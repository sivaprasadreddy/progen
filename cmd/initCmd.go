package cmd

import (
	"fmt"

	"github.com/sivaprasadreddy/progen/generators/springboot"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize " + springboot.ProjectConfigFile,
	Long:  "Initialize " + springboot.ProjectConfigFile + " with default values",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := springboot.GenerateInitConfig(); err != nil {
			return err
		}
		fmt.Println("Generated " + springboot.ProjectConfigFile + " successfully")
		return nil
	},
}
