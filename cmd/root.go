package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "progen",
	Short: "Spring Boot Project Generator CLI",
	Long:  `A CLI tool for generating Spring Boot applications`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}
		err = invokeGenerator(configFile)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Version = VERSION
	rootCmd.Flags().StringP("config", "c", "", "--config .progen.json")
}
