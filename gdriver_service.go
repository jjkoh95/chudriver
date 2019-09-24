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

// GDriver - Wrapper for all services
type GDriver struct {
	Drive  *drive.Service
	Sheets *sheets.Service
	Docs   *docs.Service
}

// getClient - This is some sort of middleware to connect to google
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

// LoginDrive - Function wrapper to login drive
func (gdriver *GDriver) LoginDrive(credentialFile, scope string) {
	gdriver.Drive, _ = drive.New(getClient(credentialFile, scope))
}

// LoginSheets - Function wrapper to login sheet
func (gdriver *GDriver) LoginSheets(credentialFile, scope string) {
	gdriver.Sheets, _ = sheets.New(getClient(credentialFile, scope))
}

// LoginDocs - Function wrapper to login docs
func (gdriver *GDriver) LoginDocs(credentialFile, scope string) {
	gdriver.Docs, _ = docs.New(getClient(credentialFile, scope))
}

// ExportSheetToJSON - Export sheet to json file
func (gdriver *GDriver) ExportSheetToJSON(sheetID, baseFilePath string) {
	sheets := gdriver.GetAllSheetTitles(sheetID)
	for _, sheet := range sheets {
		// making this fancy
		go func(sheet string) {
			resp, err := gdriver.Sheets.Spreadsheets.Values.Get(sheetID, sheet).Do()
			if err == nil {
				WriteJSON(resp, fmt.Sprintf("%s%s", baseFilePath, sheet))
			}
		}(sheet)
	}
}

// ReadSheet - Read sheet as []map[string]interface{}
func (gdriver *GDriver) ReadSheet(spreadsheetID, readRange string) []map[string]interface{} {
	resp, err := gdriver.Sheets.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil
	}
	return convertRowsToJSON(parseInterface2String(resp.Values))
}

// CreateSheet - Create a new Spreadsheet
func (gdriver *GDriver) CreateSheet() {

}

// AddSheet - add a blank sheet
func (gdriver *GDriver) AddSheet() {

}

// InsertRow - insert row to sheet
func (gdriver *GDriver) InsertRow() {

}

// ReadDoc - Read doc as raw string
func (gdriver *GDriver) ReadDoc(documentID string) (string, error) {
	doc, err := gdriver.Docs.Documents.Get(documentID).Do()
	if err != nil {
		return "", err
	}
	return readStructuralElements(doc.Body.Content), nil
}

// ExportDoc - Export doc to raw text file
func (gdriver *GDriver) ExportDoc(documentID, baseFilePath string) {
	doc, err := gdriver.Docs.Documents.Get(documentID).Do()
	if err != nil {
		log.Println(err.Error())
	}
	WriteText(readStructuralElements(doc.Body.Content), fmt.Sprintf("%s%s", baseFilePath, doc.Title))
}

// GetAllSheetTitles - Return all sheets titles
func (gdriver *GDriver) GetAllSheetTitles(spreadsheetID string) []string {
	resp, err := gdriver.Sheets.Spreadsheets.Get(spreadsheetID).Do()
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
// TODO - needs a smart method to get mimetype
func (gdriver *GDriver) DownloadFile(fileID string) error {
	return nil
}

// UploadFile - Upload file to google drive
func (gdriver *GDriver) UploadFile(filename, parentID string) (*drive.File, error) {
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

	file, err := gdriver.Drive.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CreateFolder - Create a folder in google drive
func (gdriver *GDriver) CreateFolder(folderName, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := gdriver.Drive.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}
