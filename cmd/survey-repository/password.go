package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/influenzanet/survey-repository/pkg/utils"
	"github.com/urfave/cli/v3"
)

func hashPassword(ctx context.Context, cmd *cli.Command) error {
	password := cmd.StringArg("password")

	if password == "" {
		return errors.New("Empty password")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	fmt.Printf("Hash : '%s'\n", hash)
	return nil
}
