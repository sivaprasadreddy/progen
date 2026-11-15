package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sivaprasadreddy/progen/generators/springboot"
	"github.com/spf13/cobra"
)

// Execute runs the CLI. It is the only place that terminates the process;
// library packages under generators/ always return errors to the caller.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, springboot.ErrUserAborted) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "progen",
	Short:         "Spring Boot Project Generator CLI",
	Long:          `A CLI tool for generating Spring Boot applications`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		if err := invokeGenerator(configFile); err != nil {
			return err
		}
		fmt.Println("Project generated successfully")
		return nil
	},
}

func init() {
	rootCmd.Version = VERSION
	rootCmd.Flags().StringP("config", "c", "", "--config "+springboot.ProjectConfigFile)
}
