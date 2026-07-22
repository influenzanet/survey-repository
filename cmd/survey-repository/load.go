package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/influenzanet/survey-repository/pkg/surveys"
	"github.com/urfave/cli/v3"
)

func loadCommand(ctx context.Context, cmd *cli.Command) error {

	file := cmd.StringArg("file")

	if file == "" {
		log.Fatal("File argument is empty")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to load file %s", err)
	}
	survey, err := surveys.ExtractSurveyMetadata(data)
	if err != nil {
		log.Fatalf("Unable to load file %s", err)
	}
	b, err := json.Marshal(survey)
	if err != nil {
		log.Fatalf("Unable to serialize to json %s", err)
	}
	fmt.Println(string(b))
	return nil
}
