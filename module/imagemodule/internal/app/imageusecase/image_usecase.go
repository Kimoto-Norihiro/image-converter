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

func (u *ImageUsecase) ListImages(ctx context.Context) ([]imagemodel.Image, error) {
	images, err := u.imageRepo.ListImages(ctx, u.db)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (u *ImageUsecase) CreateImage(ctx context.Context, objectName string, resizeWidthPercent int, resizeHeightPercent int, encodeFormat imagemodel.EncodeFormat) error {	
	u.db.Transaction(func(tx *gorm.DB) error {
		entity, err := imageentity.NewImageEntityToCreate(objectName, resizeWidthPercent, resizeHeightPercent, encodeFormat)
		if err != nil {
			return err
		}

		err = u.imageRepo.Create(ctx, tx, entity)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (u *ImageUsecase) UpdateImage(ctx context.Context, id int64, statusID *imagemodel.ImageStatus, convertedImageURL *string) error {
	var entity *imageentity.ImageEntity

	u.db.Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = u.imageRepo.FindForUpdate(ctx, tx, id)
		if err != nil {
			return err
		}
		if entity == nil {
			return errors.New("image not found")
		}
		if statusID != nil {
			entity.SetStatus(*statusID)
		}
		if convertedImageURL != nil {
			entity.SetConvertedImageURL(*convertedImageURL)
		}
		err = u.imageRepo.Update(ctx, tx, entity)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}