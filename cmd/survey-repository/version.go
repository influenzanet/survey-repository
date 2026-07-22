package main

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/influenzanet/survey-repository/pkg/version"
	"github.com/urfave/cli/v3"
)

func showVersion(ctx context.Context, cmd *cli.Command) error {
	info, ok := debug.ReadBuildInfo()

	v := version.Version()

	fmt.Printf("version %s\n", v.Tag)
	fmt.Printf("Revision %s\n", v.Revision)
	fmt.Printf("Dirty %t\n", v.Dirty)

	if !ok {
		return fmt.Errorf("Build info are not available")
	}

	fmt.Println(info)
	return nil
}
