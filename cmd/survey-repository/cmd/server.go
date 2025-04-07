package cmd

import (
	"log"
	"context"
	"os"
	"os/signal"
	"syscall"
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

		cfg.Show()

		mg := manager.NewManager(cfg)
		err = mg.Start()

		if err != nil {
			log.Fatalf("Unable to start manager : %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // cancel when we are finished consuming integers
	

		err = mg.StartRoutines(ctx)
		if err != nil {
			log.Fatalf("Unable to start routines : %s", err)
		}

		server := server.NewHttpServer(cfg, mg)
		
		go func() {
			if err := server.Start(); err != nil {
				log.Fatalf("Unable to start manager : %s", err)
			}
		}()
		
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
		// Block until a signal is received.
		sig := <-c
		log.Println("receiving signal, stopping services ", sig)
	
		server.Shutdown()
		
	},
}
