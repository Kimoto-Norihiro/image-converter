package gcs

import "context"

type GCSUsecase interface {
	DownloadFile(ctx context.Context, gcsFileName string) error
	UploadFile(ctx context.Context, localFileName string) error
}
