package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func listNamespace(ctx context.Context, cmd *cli.Command) error {
	mg := loadManager()
	nn := mg.GetNamespaces()
	for id, name := range nn {
		fmt.Printf("- %d '%s'\n", id, name)
	}
	return nil
}

var ErrEmptyNamespace = errors.New("Namespace name must not be empty")

func createNamespace(ctx context.Context, cmd *cli.Command) error {

	name := cmd.StringArg("name")

	if name == "" {
		return ErrEmptyNamespace
	}

	mg := loadManager()
	id, err := mg.CreateNamespace(name)
	if err != nil {
		return err
	}
	fmt.Printf("namespace created %d '%s'\n", id, name)
	return nil
}
