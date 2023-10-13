package gcsclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSBucket struct {
	ctx    context.Context
	bucket *storage.BucketHandle
}

func NewGCSBucket(credentialsFile string, bucketName string) GCSBucket {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatal(err)
	}

	return GCSBucket{
		ctx:    ctx,
		bucket: client.Bucket(bucketName),
	}
}

func (gcs GCSBucket) UploadFile(file *os.File, gcsFileName string) {
	// バケット内のアップロード先のオブジェクトを作成
	obj := gcs.bucket.Object(gcsFileName)

	// ファイルをアップロード
	wc := obj.NewWriter(gcs.ctx)
	if _, err := io.Copy(wc, file); err != nil {
		log.Fatal(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("file upload success")
}

func (gcs GCSBucket) DownloadFile(localFileName string, gcsFileName string) {
	// ファイルの作成
	file, err := os.Create(localFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// バケットからオブジェクトをダウンロード
	rc, err := gcs.bucket.Object(gcsFileName).NewReader(gcs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	// ファイルに出力
	if _, err := file.ReadFrom(rc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("file download success")
}
