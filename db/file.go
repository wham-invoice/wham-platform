package db

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/juju/errors"
)

// UploadFile uploads an object.
func (a *App) UploadFile(ctx context.Context, object, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return errors.Annotatef(err, "failed to open file: %s", filePath)
	}
	defer f.Close()

	// NOTE what this do?
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	bucket, err := a.storageClient.DefaultBucket()
	if err != nil {
		return errors.Trace(err)
	}

	wc := bucket.Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return errors.Trace(err)
	}
	if err := wc.Close(); err != nil {
		return errors.Trace(err)
	}

	log.Printf("file %s uploaded as %s", filePath, object)
	return nil
}
