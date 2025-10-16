package cmd

import (
	"github.com/sivaprasadreddy/progen/generators/helpers"
	"github.com/sivaprasadreddy/progen/generators/springboot"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .progen.json",
	Long:  `Initialize .progen.json with default values`,
	Run: func(cmd *cobra.Command, args []string) {
		err := springboot.GenerateInitConfig()
		helpers.FatalIfErrOrMsg(err, "Generated .progen.json successfully")
	},
}
