package gcsfunction

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSFunction struct {
	client *storage.Client
}

func NewGCSFunction(ctx context.Context, credentialsFile string) *GCSFunction {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		fmt.Errorf("failed to create client: %v", err)
	}

	return &GCSFunction{
		client: client,
	}
}

func (u *GCSFunction) DownloadFile(ctx context.Context, gcsFileName string) error {
	bucket := u.client.Bucket(os.Getenv("NON_CONVERTED_BUCKET_NAME"))

	localFilePath := fmt.Sprintf("./img/%s", gcsFileName)
	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file in DownloadFile: %v", err)
	}
	defer file.Close()

	rc, err := bucket.Object(gcsFileName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reader in DownloadFile: %v", err)
	}
	defer rc.Close()

	if _, err := file.ReadFrom(rc); err != nil {
		return fmt.Errorf("failed to read file in DownloadFile: %v", err)
	}

	fmt.Printf("file download success: %s\n", gcsFileName)
	return nil
}

func (u *GCSFunction) UploadFile(ctx context.Context, fileName string) (*string, error) {
	bucket := u.client.Bucket(os.Getenv("CONVERTED_BUCKET_NAME"))

	filePath := fmt.Sprintf("./img/%s", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file in UploadFile: %v", err)
	}
	defer file.Close()

	obj := bucket.Object(fileName)

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("failed to copy file in UploadFile: %v", err)
	}
	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file in UploadFile: %v", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get attrs in UploadFile: %v", err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove file in UploadFile: %v", err)
	}

	fmt.Printf("file upload success: %s\n", fileName)
	return &attrs.MediaLink, nil
}

func (u *GCSFunction) UploadNonConvertedFile(ctx context.Context, reader *bytes.Reader, fileName string) error {
	bucket := u.client.Bucket(os.Getenv("NON_CONVERTED_BUCKET_NAME"))

	obj := bucket.Object(fileName)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, reader); err != nil {
		return fmt.Errorf("failed to copy file in UploadNonConvertedFile: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close file in UploadNonConvertedFile: %v", err)
	}

	fmt.Printf("file upload success: %s\n", fileName)
	return nil
}
