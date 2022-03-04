package util

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	zlog, err := zap.NewProduction()

	if err != nil {
		// yes panic, this is a developer error.
		panic(errors.ErrorStack(err))
	}

	Logger = zlog.Sugar()
}

func SetDebugLogger() error {
	zlog, err := zap.NewDevelopment()
	if err != nil {
		return errors.Trace(err)
	}
	Logger = zlog.Sugar()

	// Needed?
	jLog := loggo.GetLogger("juju.worker.dependency")
	jLog.SetLogLevel(loggo.TRACE)

	return nil
}
