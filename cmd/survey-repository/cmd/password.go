package cmd

import (
	"fmt"

	"github.com/influenzanet/survey-repository/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(passwordCmd)
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Hash a password",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		hash, err := utils.HashPassword(password)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hash : '%s'\n", hash)
	},
}
