# chudrive
chudrive provides common functions to access google drive

## How to use
```go
// upload
_, err := chudriveWrapper.UploadFileLocal(filename, filePathID)

// download
err := chudriveWrapper.DownloadFileLocal(filename, fileID)
```