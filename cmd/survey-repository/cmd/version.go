package cmd

import (
	"fmt"
	"runtime/debug"
	"github.com/influenzanet/survey-repository/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		
		info, ok := debug.ReadBuildInfo()
		
		fmt.Printf("version %s\n", version.Version)
		fmt.Printf("Revision %s\n", version.Revision)

		if !ok {
			fmt.Println("Build info are not available")
			return
		}

		fmt.Println(info)
	},
}
