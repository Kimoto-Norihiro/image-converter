package imageusecase

import (
	"context"
	"errors"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
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

func (u *ImageUsecase) ConvertImage(ctx context.Context) error {
	return errors.New("not implemented")
}

func (u *ImageUsecase) ListImages(ctx context.Context) ([]imagemodel.Image, error) {
	return nil, errors.New("not implemented")
}

func (u *ImageUsecase) CreateImage(ctx context.Context) error {
	return errors.New("not implemented")
}