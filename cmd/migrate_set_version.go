package cmd

import (
	"log"
	"strconv"

	"github.com/alpacahq/ribbit-backend/migration"

	"github.com/spf13/cobra"
)

// setVersionCmd represents the version command
var setVersionCmd = &cobra.Command{
	Use:   "set_version [version]",
	Short: "sets db version without running migrations",
	Long:  `sets db version without running migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		var passthrough = []string{"set_version"}
		if len(args) > 0 {
			_, err := strconv.Atoi(args[0])
			if err != nil {
				passthrough = append(passthrough, args[0])
			}
		}

		err := migration.Run(passthrough...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(setVersionCmd)
}
