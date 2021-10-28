package cmd

import (
	"fmt"
	"log"

	"github.com/alpacahq/ribbit-backend/migration"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init creates version info table in the database",
	Long:  `init creates version info table in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		err := migration.Run("init")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(initCmd)
}
