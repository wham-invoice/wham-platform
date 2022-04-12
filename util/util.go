package util

import (
	"os"
	"time"

	"github.com/juju/errors"
)

const (
	GCP_CLIENT_ID     = "GCP_CLIENT_ID"
	GCP_CLIENT_SECRET = "GCP_CLIENT_SECRET"
)

func ToFormattedDate(t time.Time) string {
	return t.Format("02-Jan-2006")
}

func GetEnvVar(key string) (string, error) {
	env := os.Getenv(key)

	if env == "" {
		return "", errors.Errorf("expected an explicit %s", key)
	}
	return env, nil
}
