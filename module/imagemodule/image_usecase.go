package imagemodule

import (
	"context"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

type ImageUsecase interface {
	CreateImage(ctx context.Context, objectName string, resizeWidthPercent int, resizeHeightPercent int, encodeFormat imagemodel.EncodeFormat) error
	ListImages(ctx context.Context) ([]imagemodel.Image, error)
	UpdateImage(ctx context.Context, id int64, statusID *imagemodel.ImageStatus, convertedImageURL *string) error
}