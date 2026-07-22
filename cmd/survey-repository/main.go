package main

import (
	"context"
	"log"
	"os"

	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/manager"
	"github.com/urfave/cli/v3"
)

const flagConfig = "config"

func main() {

	cmd := &cli.Command{
		Name: "survey-repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagConfig,
				Value:   "app.yml",
				Usage:   "Config path",
				Sources: cli.EnvVars("MONITOR_CONFIG"),
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "start",
				Description: "Start server",
				Action:      startServer,
			},
			{
				Name:        "version",
				Description: "Show version",
				Action:      showVersion,
			},
			{
				Name:        "load",
				Description: "load from file",
				Action:      loadCommand,
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "file", UsageText: "file to load"},
				},
			},
			{
				Name:        "password",
				Description: "Hash a password",
				Action:      hashPassword,
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "password", UsageText: "password"},
				},
			},
			{
				Name:        "ns",
				Description: "Namespace commands",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Description: "list namespace",
						Action:      listNamespace,
					},
					{
						Name:        "create",
						Description: "create namespace",
						Action:      createNamespace,
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "name", UsageText: "name"},
						},
					},
				},
			},
			{
				Name:        "survey",
				Description: "survey commands",
				Commands: []*cli.Command{
					{
						Name:        "show",
						Description: "show survey",
						Action:      showSurvey,
						Arguments: []cli.Argument{
							&cli.IntArg{Name: "id", UsageText: "survey id"},
						},
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
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
