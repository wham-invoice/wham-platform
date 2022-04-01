package db

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/juju/errors"
	"github.com/wham-invoice/wham-platform/util"
)

func (app *App) StorePDF(ctx context.Context, fileName, filePath string) error {
	// NOTE when cancel is called all resources using ctx are released.
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	bucket, err := app.storageClient.Bucket("wham-ad61b.appspot.com")
	if err != nil {
		return errors.Trace(err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	writer := bucket.Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(writer, f); err != nil {
		return errors.Trace(err)
	}
	if err := writer.Close(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

// PDF returns the PDF file from the storage bucket.
func (app *App) PDF(ctx context.Context, fileName string) ([]byte, error) {

	bucket, err := app.storageClient.Bucket("wham-ad61b.appspot.com")
	if err != nil {
		return nil, errors.Trace(err)
	}

	util.Logger.Infof("attempting to read file %s", fileName)
	rc, err := bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer rc.Close()
	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return body, nil
}
