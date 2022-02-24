package main

import (
	"context"
	"log"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server"

	"github.com/juju/errors"
)

func main() {
	ctx := context.Background()

	dbApp, err := db.Init(ctx)
	if err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
	defer dbApp.CloseDB()

	if err = server.Run(dbApp); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
}
