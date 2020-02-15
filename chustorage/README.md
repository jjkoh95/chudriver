# chustorage
chustorage provides friendlier APIs to perform download and upload in GCS

## How to use
```go
// upload
r, err := os.Open(filename)
err = chustorageWrapper.UploadFile(ctx, bucket, object, r)

// download
f, err := os.Create(filename)
err = chustorageWrapper.DownloadFile(ctx, bucket, object, f)
```