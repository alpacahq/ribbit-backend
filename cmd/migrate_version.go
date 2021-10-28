package cmd

import (
	"fmt"
	"log"

	"github.com/alpacahq/ribbit-backend/migration"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version prints current db version",
	Long:  `version prints current db version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version called")
		err := migration.Run("version")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(versionCmd)
}
