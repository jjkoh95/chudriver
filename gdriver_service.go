package gdriver

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

// Service - wrapper for all services
type Service struct {
	Drive  *drive.Service
	Sheets *sheets.Service
	Docs   *docs.Service
}

// getClient - this is some sort of middleware to connect to google
func getClient(credentialFile, scope string) *http.Client {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.JWTConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config.Client(context.Background())
}

// LoginDrive - function wrappr to login drive
func (service *Service) LoginDrive(credentialFile, scope string) {
	service.Drive, _ = drive.New(getClient(credentialFile, scope))
}

// LoginSheets - function wrapper to login sheet
func (service *Service) LoginSheets(credentialFile, scope string) {
	service.Sheets, _ = sheets.New(getClient(credentialFile, scope))
}

// LoginDocs - function wrapper to login docs
func (service *Service) LoginDocs(credentialFile, scope string) {
	service.Docs, _ = docs.New(getClient(credentialFile, scope))
}

// ExportSheetToJSON - export sheet to json file
func (service *Service) ExportSheetToJSON(sheetID, baseFilePath string) {
	sheets := service.GetAllSheetTitles(sheetID)
	for _, sheet := range sheets {
		// making this fancy
		go func(sheet string) {
			resp, err := service.Sheets.Spreadsheets.Values.Get(sheetID, sheet).Do()
			if err == nil {
				WriteJSON(resp, fmt.Sprintf("%s%s", baseFilePath, sheet))
			}
		}(sheet)
	}
}

// ReadSheet - read sheet as []map[string]interface{}
func (service *Service) ReadSheet(spreadsheetID, readRange string) []map[string]interface{} {
	resp, err := service.Sheets.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil
	}
	return convertRowsToJSON(parseInterface2String(resp.Values))
}

// ReadDoc - read doc as raw string
func (service *Service) ReadDoc(documentID string) (string, error) {
	doc, err := service.Docs.Documents.Get(documentID).Do()
	if err != nil {
		return "", err
	}
	return readStructuralElements(doc.Body.Content), nil
}

// ExportDoc - export doc to raw text file
func (service *Service) ExportDoc(documentID, baseFilePath string) {
	doc, err := service.Docs.Documents.Get(documentID).Do()
	if err != nil {
		log.Println(err.Error())
	}
	WriteText(readStructuralElements(doc.Body.Content), fmt.Sprintf("%s%s", baseFilePath, doc.Title))
}

// GetAllSheetTitles - return all sheets titles
func (service *Service) GetAllSheetTitles(spreadsheetID string) []string {
	resp, err := service.Sheets.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	titles := make([]string, len(resp.Sheets))
	for i, sheet := range resp.Sheets {
		titles[i] = sheet.Properties.Title
	}
	return titles
}

// DownloadFile - Download file from google drive
func (service *Service) DownloadFile(fileID string) error {
	return nil
}

// UploadFile - upload file to google drive
func (service *Service) UploadFile(filename, parentID string) (*drive.File, error) {
	mimeType := mime.TypeByExtension(filename)

	f := &drive.File{
		MimeType: mimeType,
		Name:     filename,
		Parents:  []string{parentID},
	}

	content, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	file, err := service.Drive.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CreateFolder - create a folder in google drive
func (service *Service) CreateFolder(folderName, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := service.Drive.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}
