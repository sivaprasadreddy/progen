package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"charm.land/huh/v2"
	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
	"github.com/sivaprasadreddy/progen/internal/version"
	"github.com/spf13/cobra"
)

// Execute runs the CLI. It is the only place that terminates the process;
// packages under internal/generators always return errors to the caller.
func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		stop()
		if errors.Is(err, huh.ErrUserAborted) {
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
		if err := invokeGenerator(cmd.Context(), configFile); err != nil {
			return err
		}
		fmt.Println("Project generated successfully")
		return nil
	},
}

func init() {
	rootCmd.Version = version.Version
	rootCmd.Flags().StringP("config", "c", "", "--config "+springboot.ProjectConfigFile)
}
