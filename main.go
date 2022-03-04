package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rstorr/wham-platform/db"
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

	dbApp, err := db.Init(ctx)
	if err != nil {
		util.Logger.Fatal(errors.ErrorStack(err))
	}
	defer dbApp.CloseDB()

	if err = server.Run(dbApp); err != nil {
		util.Logger.Fatal(errors.ErrorStack(err))
	}
}
