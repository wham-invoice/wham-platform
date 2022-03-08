package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rstorr/wham-platform/server"
	"github.com/rstorr/wham-platform/util"

	"github.com/juju/errors"
)

func main() {
	ctx := context.Background()

	if err := util.SetDebugLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	if err := server.Run(ctx); err != nil {
		util.Logger.Fatal(errors.ErrorStack(err))
	}
}
