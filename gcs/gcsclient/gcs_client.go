package gcsclient

import (
	"context"

	"github.com/Kimoto-Norihiro/image-converter/gcs"
	"github.com/Kimoto-Norihiro/image-converter/gcs/internal/app/gcsfunction"
)

func New() *gcs.GCS {
	ctx := context.Background()
	credentialsFile := "../../credentials.json"
	return &gcs.GCS{
		GCSFunction: gcsfunction.NewGCSFunction(ctx, credentialsFile),
	}
}
