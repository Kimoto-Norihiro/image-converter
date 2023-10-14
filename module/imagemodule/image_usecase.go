package imagemodule

import (
	"context"
)

type ImageUsecase interface {
	ConvertImages(c context.Context) error
}