package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
// It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "progen",
	Short: "Project generator CLI",
	Long:  `A project generator CLI for Java, Spring Boot, Go applications`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := invokeGenerator()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Version = VERSION
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.Flags().StringVarP(&generatorType, "type", "t", "", "Application Type")
}
