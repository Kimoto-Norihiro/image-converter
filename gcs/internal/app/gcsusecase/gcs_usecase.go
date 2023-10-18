package gcsusecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSClientUsecase struct {
	client *storage.Client
}

func NewGCSClientUsecase(ctx context.Context, credentialsFile string) *GCSClientUsecase {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatal(err)
	}

	return &GCSClientUsecase{
		client: client,
	}
}

func (u *GCSClientUsecase) DownloadFile(ctx context.Context, gcsFileName string) error {
	bucketName := "image-nonconverted"
	localFileName := gcsFileName
	bucket := u.client.Bucket(bucketName)

	file, err := os.Create(localFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// バケットからオブジェクトをダウンロード
	rc, err := bucket.Object(gcsFileName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	// ファイルに出力
	if _, err := file.ReadFrom(rc); err != nil {
		return err
	}

	fmt.Println("file download success")
	return nil
}

func (u *GCSClientUsecase) UploadFile(ctx context.Context, fileName string) error {
	bucketName := "image-converted"
	bucket := u.client.Bucket(bucketName)

	filePath := fmt.Sprintf("/tmp/%s", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// バケット内のアップロード先のオブジェクトを作成
	obj := bucket.Object(fileName)

	// ファイルをアップロード
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		log.Fatal(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}

	// fileを削除
	err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (u *GCSClientUsecase) UploadNonConvertedFile(ctx context.Context, reader *bytes.Reader, fileName string) error {
	bucketName := "image-nonconverted"
	bucket := u.client.Bucket(bucketName)

	obj := bucket.Object(fileName)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, reader); err != nil {
		log.Fatal(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
