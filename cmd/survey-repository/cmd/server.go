package cmd

import (
	"log"

	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/manager"
	"github.com/influenzanet/survey-repository/pkg/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ServerCmd)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch server",
	Long:  `Load HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Unable to load config : %s", err)
		}

		mg := manager.NewManager(cfg)
		err = mg.Start()

		if err != nil {
			log.Fatalf("Unable to start manager : %s", err)
		}

		server := server.NewHttpServer(cfg, mg)
		server.Start()

		if err != nil {
			log.Fatalf("Unable to start manager : %s", err)
		}

	},
}
