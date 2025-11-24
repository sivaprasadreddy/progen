package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "v1.0.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Progen",
	Long:  `Print the version number of Progen CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Progen : Spring Boot Application Generator %s\n", VERSION)
	},
}
