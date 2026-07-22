package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func showSurvey(ctx context.Context, cmd *cli.Command) error {
	mg := loadManager()
	id := cmd.IntArg("id")
	data, err := mg.GetSurveyData(uint(id), true)
	if err != nil {
		return fmt.Errorf("Unable to show survey : %s", err)
	}
	fmt.Println(string(data))
	return nil
}
