package chudrive_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/jjkoh95/chudriver/chuauth"
	"github.com/jjkoh95/chudriver/chudrive"
	"google.golang.org/api/drive/v3"
)

const (
	testPrefix           = "chudrive-test-"
	testDriveID          = "1F8CvTSdLFqPuvYya6AeCYkX1u7wISLrT"
	testDriveName        = "chudrive"
	testFilePathID       = "1aZgRU8MJFFFttScVR1oVvvokjP9SZ7FZ"
	testFileName         = "babe.jpg"
	testDownloadFileID   = "1rueVP_EGaFtu_amkNOeCktiqnllKXvmF"
	testDownloadFileName = "test-babe.jpg"
)

var chudriveWrapper chudrive.Chudrive

func TestMain(m *testing.M) {
	login()
	code := m.Run()
	os.Exit(code)
}

func login() {
	chudriveWrapper.Drive, _ = drive.New(chuauth.GetClientFromJSON("credentials.json", drive.DriveScope))
}

func TestCreateFolder(t *testing.T) {
	folderName := testPrefix + uuid.New().String()
	_, err := chudriveWrapper.CreateFolder(folderName, testDriveID)
	if err != nil {
		t.Error("Should create directory without error")
	}
}

func TestList(t *testing.T) {
	res, err := chudriveWrapper.ListByQuery("")
	if err != nil {
		t.Error("Expected to call ListByQuery without error")
	}
	if _, ok := res[testDriveName]; !ok {
		t.Error("Expected to get test folder")
	}

	folderQuery := fmt.Sprintf("\"%s\" in parents", testDriveID)
	res, err = chudriveWrapper.ListFolder(folderQuery)
	if err != nil {
		t.Error("Expected to call ListFolder without error")
	}

	fileQuery := fmt.Sprintf("\"%s\" in parents", testDriveID)
	res, err = chudriveWrapper.ListFile(fileQuery)
	if err != nil {
		t.Error("Expected to call ListFile without error")
	}
	if len(res) != 0 {
		t.Error("Expected no file in ListFile")
	}
}

func TestUploadAndDeleteFile(t *testing.T) {
	driveFile, err := chudriveWrapper.UploadFile(testFileName, testFilePathID)
	if err != nil {
		t.Error("Expected to upload file without error")
	}

	res, err := chudriveWrapper.ListByQuery("")
	if _, ok := res[testFileName]; !ok {
		t.Error("Expected to find the newly uploaded file")
	}

	err = chudriveWrapper.DeleteFile(driveFile.Id)
	if err != nil {
		t.Error("Expected to trigger delete file action without error")
	}
}

func TestDownloadFile(t *testing.T) {
	err := chudriveWrapper.DownloadFile(testDownloadFileName, testDownloadFileID)
	if err != nil {
		fmt.Println(err)
		t.Error("Expected to download file without error")
	}
}
