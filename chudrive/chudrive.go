package chudrive

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"

	"google.golang.org/api/drive/v3"
)

const (
	folderMimeType  = "application/vnd.google-apps.folder"
	defaultPageSize = 1000
)

// Chudrive - Google drive wrapper.
type Chudrive struct {
	Drive *drive.Service
}

// CreateFolder - Create a folder.
func (chudrive *Chudrive) CreateFolder(folderName, parentID string) (*drive.File, error) {

	d := &drive.File{
		Name:     folderName,
		MimeType: folderMimeType,
		Parents:  []string{parentID},
	}

	file, err := chudrive.Drive.Files.Create(d).Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

// ListByQuery - List everything by query.
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

// ListFolder - List folder.
func (chudrive *Chudrive) ListFolder(query string) (map[string]string, error) {
	newQuery := fmt.Sprintf("mimeType = \"%s\"", folderMimeType)
	if query != "" {
		newQuery = fmt.Sprintf("%s and %s", query, newQuery)
	}
	return chudrive.ListByQuery(newQuery)
}

// ListFile - List file, it is super weird, google treats directory as files too.
func (chudrive *Chudrive) ListFile(query string) (map[string]string, error) {
	newQuery := fmt.Sprintf("mimeType != \"%s\"", folderMimeType)
	if query != "" {
		newQuery = fmt.Sprintf("%s and %s", query, newQuery)
	}
	return chudrive.ListByQuery(newQuery)
}

// UploadFileLocal - Upload file to path.
func (chudrive *Chudrive) UploadFileLocal(filename, parentID string) (*drive.File, error) {
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

// DownloadFileLocal - Download file to local.
func (chudrive *Chudrive) DownloadFileLocal(filename, fileID string) error {
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

// DeleteFile - Delete file.
func (chudrive *Chudrive) DeleteFile(fileID string) error {
	return chudrive.Drive.Files.Delete(fileID).Do()
}

// TransferOwnership - Transfer ownership.
// This is particularly useful when you are using service account to create.
// Please make sure you own whatever file instead of your service account!
func (chudrive *Chudrive) TransferOwnership(fileID, email string) error {
	perm := drive.Permission{
		EmailAddress: email,
		Role:         "owner",
		Type:         "user",
	}
	_, err := chudrive.Drive.Permissions.Create(fileID, &perm).TransferOwnership(true).Do()
	return err
}
