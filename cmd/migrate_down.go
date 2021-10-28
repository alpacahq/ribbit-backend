package cmd

import (
	"fmt"
	"log"

	"github.com/alpacahq/ribbit-backend/migration"

	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "reverts the last migration",
	Long:  `reverts the last migration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("down called")
		err := migration.Run("down")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(downCmd)
}
