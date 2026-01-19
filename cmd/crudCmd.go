package cmd

import (
	"fmt"

	"github.com/sivaprasadreddy/progen/generators/crud"
	"github.com/spf13/cobra"
)

var crudCmd = &cobra.Command{
	Use:   "crud [entity]",
	Short: "Generate CRUD code for an entity",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an entity name. Example: progen crud User")
			return
		}

		entity := args[0]
		crud.Generate(entity)

		fmt.Println("CRUD file generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(crudCmd)
}
