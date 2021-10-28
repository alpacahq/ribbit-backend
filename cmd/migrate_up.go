package cmd

import (
	"log"
	"strconv"

	"github.com/alpacahq/ribbit-backend/migration"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up [target]",
	Short: "runs all available migrations or up to the target if provided",
	Long:  `runs all available migrations or up to the target if provided`,
	Run: func(cmd *cobra.Command, args []string) {
		var passthrough = []string{"up"}
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
	migrateCmd.AddCommand(upCmd)
}
