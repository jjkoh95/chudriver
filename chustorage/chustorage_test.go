package chustorage_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/jjkoh95/chudriver/chustorage"
	"google.golang.org/api/option"
)

const testFile = "babe.jpg"
const bucketName = "chudriver"
const bucketFolder = "test"
const bucketFilename = "test/babe.jpg"
const testBucketFile = "babe_test.jpg"

var chustorageWrapper chustorage.Chustorage

func TestMain(m *testing.M) {
	login()
	code := m.Run()
	os.Exit(code)
}

func login() {
	clientOptions := option.WithCredentialsFile("credentials.json")
	chustorageWrapper.Storage, _ = storage.NewClient(context.Background(), clientOptions)
}

func TestUploadAndDeleteFile(t *testing.T) {
	r, err := os.Open(testFile)
	if err != nil {
		t.Error("Expected to open test file without error")
	}
	err = chustorageWrapper.UploadFile(context.Background(), bucketName, bucketFilename, r)
	if err != nil {
		t.Error("Expected to upload file without error")
	}

	err = chustorageWrapper.DeleteFile(context.Background(), bucketName, bucketFilename)
	if err != nil {
		t.Error("Expected to delete file without error")
	}
}

func TestDownloadFile(t *testing.T) {
	f, err := os.Create(testBucketFile)
	if err != nil {
		t.Error("Expected to create file without error")
	}
	defer f.Close()

	err = chustorageWrapper.DownloadFile(context.Background(), bucketName, testBucketFile, f)
	if err != nil {
		t.Error("Expected to download file from bucket without error")
	}
}

func TestReadFile(t *testing.T) {
	data, err := chustorageWrapper.ReadFile(context.Background(), bucketName, testBucketFile)
	if err != nil {
		t.Error("Expected to read file without error")
	}
	if len(data) == 0 {
		t.Error("Expected to have content in ReadFile")
	}
}
