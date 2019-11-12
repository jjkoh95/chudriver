package chudrive

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"

	"google.golang.org/api/drive/v3"
)

const (
	folderMimeType  = "application/vnd.google-apps.folder"
	defaultPageSize = 1000
)

// Chudrive - google drive wrapper
type Chudrive struct {
	Drive *drive.Service
}

// CreateFolder - create a folder
func (chudrive *Chudrive) CreateFolder(folderName, parentID string) (*drive.File, error) {

	d := &drive.File{
		Name:     folderName,
		MimeType: folderMimeType,
		Parents:  []string{parentID},
	}

	file, err := chudrive.Drive.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

// ListByQuery - list everything by query
func (chudrive *Chudrive) ListByQuery(query string) (map[string]string, error) {
	fileList, err := chudrive.Drive.Files.List().
		PageSize(defaultPageSize).Q(query).
		Do()
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for _, f := range fileList.Files {
		res[f.Name] = f.Id
	}
	return res, nil
}

// ListFolder - list folder
func (chudrive *Chudrive) ListFolder(query string) (map[string]string, error) {
	newQuery := fmt.Sprintf("mimeType = \"%s\"", folderMimeType)
	if query != "" {
		newQuery = fmt.Sprintf("%s and %s", query, newQuery)
	}
	return chudrive.ListByQuery(newQuery)
}

// ListFile - list file, it is super weird, google treats directory as files too
func (chudrive *Chudrive) ListFile(query string) (map[string]string, error) {
	newQuery := fmt.Sprintf("mimeType != \"%s\"", folderMimeType)
	if query != "" {
		newQuery = fmt.Sprintf("%s and %s", query, newQuery)
	}
	return chudrive.ListByQuery(newQuery)
}

// UploadFile - upload file to path
func (chudrive *Chudrive) UploadFile(filename, parentID string) (*drive.File, error) {
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

	file, err := chudrive.Drive.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

// DownloadFile - download file to local
func (chudrive *Chudrive) DownloadFile(filename, fileID string) error {
	res, err := chudrive.Drive.Files.Get(fileID).Download()
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, bodyBytes, 0644)
	return nil
}

// DeleteFile - delete file
func (chudrive *Chudrive) DeleteFile(fileID string) error {
	return chudrive.Drive.Files.Delete(fileID).Do()
}
