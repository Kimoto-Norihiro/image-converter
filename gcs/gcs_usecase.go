package gcs

import (
	"bytes"
	"context"
)

type GCSUsecase interface {
	DownloadFile(ctx context.Context, gcsFileName string) error
	UploadFile(ctx context.Context, fileName string) error
	UploadNonConvertedFile(ctx context.Context, reader *bytes.Reader ,fileName string) error
}
