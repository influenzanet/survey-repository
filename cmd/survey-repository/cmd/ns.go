package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nsCmd)
	nsCmd.AddCommand(NsListCmd)
	nsCmd.AddCommand(NsCreateCmd)
}

var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "Namespace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var NsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show Namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		mg := loadManager()
		nn := mg.GetNamespaces()
		for id, name := range nn {
			fmt.Printf("- %d '%s'\n", id, name)
		}
	},
}

var NsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create Namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mg := loadManager()
		name := args[0]
		id, err := mg.CreateNamespace(name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("namespace created %d '%s'\n", id, name)
	},
}
