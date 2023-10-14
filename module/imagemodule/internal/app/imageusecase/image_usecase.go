package imageusecase

import (
	"context"
	"errors"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"gorm.io/gorm"
)

type ImageUsecase struct {
	db *gorm.DB
	imageRepo imageentity.ImageRepository
}

func NewImageUsecase(
	db *gorm.DB,
	imageRepo imageentity.ImageRepository,
) *ImageUsecase {
	return &ImageUsecase{
		db:        db,
		imageRepo: imageRepo,
	}
}

func (u *ImageUsecase) ConvertImages(ctx context.Context) error {
	return errors.New("not implemented")
}