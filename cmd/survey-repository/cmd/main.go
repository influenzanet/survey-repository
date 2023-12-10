package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/manager"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "survey-repository",
	Short: "Survey repository",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func loadManager() *manager.Manager {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to load config : %s", err)
	}

	mg := manager.NewManager(cfg)
	err = mg.Start()

	if err != nil {
		log.Fatalf("Unable to start manager : %s", err)
	}
	return mg
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
