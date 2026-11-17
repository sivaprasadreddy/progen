package cli

import (
	"fmt"

	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
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
		if err := springboot.WriteConfigFile(springboot.DefaultProjectConfig(), springboot.ProjectConfigFile); err != nil {
			return err
		}
		fmt.Println("Generated " + springboot.ProjectConfigFile + " successfully")
		return nil
	},
}
