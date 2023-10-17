package imagemodule

import (
	"context"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

type ImageUsecase interface {
	CreateImage(ctx context.Context) error
	ListImages(ctx context.Context) ([]imagemodel.Image, error)
	ConvertImage(ctx context.Context) error
}