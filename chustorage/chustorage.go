package chustorage

import (
	"context"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// Chustorage - Google Cloud Storage (GCS) wrapper
type Chustorage struct {
	Storage *storage.Client
}

// ReadFile - Read file from storage as raw bytes.
func (chustorage *Chustorage) ReadFile(ctx context.Context, bucket, bucketFilename string) ([]byte, error) {
	obj, err := chustorage.getObject(ctx, bucket, bucketFilename)
	if err != nil {
		return nil, err
	}

	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DownloadFile - Download file from storage.
func (chustorage *Chustorage) DownloadFile(ctx context.Context, bucket, bucketFilename string, fileWriter io.Writer) error {
	obj, err := chustorage.getObject(ctx, bucket, bucketFilename)
	if err != nil {
		return err
	}

	rc, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	// write file
	if _, err := io.Copy(fileWriter, rc); err != nil {
		return nil
	}

	return nil
}

// UploadFile - Upload file to storage.
func (chustorage *Chustorage) UploadFile(ctx context.Context, bucket, bucketFilename string, fileReader io.Reader) error {
	obj, err := chustorage.getObject(ctx, bucket, bucketFilename)
	if err != nil {
		return err
	}

	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, fileReader); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// DeleteFile - Delete file from storage.
func (chustorage *Chustorage) DeleteFile(ctx context.Context, bucket, bucketFilename string) error {
	obj, err := chustorage.getObject(ctx, bucket, bucketFilename)
	if err != nil {
		return err
	}

	if err := obj.Delete(ctx); err != nil {
		return err
	}

	return nil
}

// getObject - Get access to object/file in storage.
func (chustorage *Chustorage) getObject(ctx context.Context, bucket, bucketFilename string) (*storage.ObjectHandle, error) {
	bh := chustorage.Storage.Bucket(bucket)
	// check if bucket exists
	if _, err := bh.Attrs(ctx); err != nil {
		return nil, err
	}
	obj := bh.Object(bucketFilename)
	return obj, nil
}
