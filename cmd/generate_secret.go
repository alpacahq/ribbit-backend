package cmd

import (
	"fmt"
	"log"

	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/spf13/cobra"
)

// createCmd represents the migrate command
var generateSecretCmd = &cobra.Command{
	Use:   "generate_secret",
	Short: "generate_secret",
	Long:  `generate_secret`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := secret.GenerateRandomString(256)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\nJWT_SECRET=%s\n\n", s)
	},
}

func init() {
	rootCmd.AddCommand(generateSecretCmd)
}
