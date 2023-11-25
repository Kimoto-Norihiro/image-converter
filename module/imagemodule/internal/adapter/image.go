package imagerepository

import (
	"context"

	imageentity "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"github.com/Kimoto-Norihiro/image-converter/utils/myerror"
	"gorm.io/gorm"
)

func (r *repo) ListImages(ctx context.Context, db *gorm.DB) (_ []imagemodel.Image, reterr error) {
	defer myerror.Wrap(&reterr, "ListImages in ImageRepository")

	var images []ImageDTO
	err := db.Find(&images).Error
	if err != nil {
		return nil, err
	}

	return ImagesFromDTO(images)
}

func (r *repo) Update(ctx context.Context, tx *gorm.DB, entity *imageentity.ImageEntity) (reterr error) {
	defer myerror.Wrap(&reterr, "Update in ImageRepository")
	
	upValues, err := entity.UpdateValues()
	if err != nil {
		return err
	}

	m := dtoImageUpdatedValues(upValues)
	if len(m) == 0 {
		return nil
	}

	return tx.Model(&ImageDTO{}).
		Where("id = ?", entity.Image().ID).
		Updates(m).
		Error
}

func (r *repo) Create(ctx context.Context, tx *gorm.DB, entity *imageentity.ImageEntity) (reterr error) {
	defer myerror.Wrap(&reterr, "Create in ImageRepository")
	
	ImageDTO, err := DTOFromImage(entity.Image())
	if err != nil {
		return err
	}
	return tx.Create(ImageDTO).Error
}

func (r *repo) FindForUpdate(ctx context.Context, tx *gorm.DB, id int64) (_ *imageentity.ImageEntity, reterr error) {
	defer myerror.Wrap(&reterr, "FindForUpdate in ImageRepository")
	
	var image ImageDTO
	err := tx.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	imageDTO, err := ImageFromDTO(&image)
	if err != nil {
		return nil, err
	}
	return imageentity.NewImageEntityToUpdate(imageDTO), nil
}

func dtoImageUpdatedValues(v *imageentity.ImageUpdatedValues) map[string]any {
	if v == nil {
		return nil
	}

	m := make(map[string]any)
	if status, ok := v.Status.Get(); ok {
		m["status_id"] = status
	}
	if convertedImageURL, ok := v.ConvertedImageURL.Get(); ok {
		m["converted_image_url"] = convertedImageURL
	}

	return m
}
