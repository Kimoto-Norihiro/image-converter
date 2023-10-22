package gcs

import (
	"bytes"
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=${GOPACKAGE}mock/$GOFILE -package=${GOPACKAGE}mock

type GCSFunction interface {
	DownloadFile(ctx context.Context, gcsFileName string) error
	UploadFile(ctx context.Context, fileName string) (*string, error)
	UploadNonConvertedFile(ctx context.Context, reader *bytes.Reader, fileName string) error
}
