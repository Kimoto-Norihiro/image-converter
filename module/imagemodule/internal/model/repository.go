package imageentity

import (
	"context"

	"gorm.io/gorm"
)

type ImageRepository interface {
	ListImages(ctx context.Context, db *gorm.DB) ([]ImageEntity, error)
	Update(ctx context.Context, db *gorm.DB, entity *ImageEntity) error
	Convert(ctx context.Context, db *gorm.DB, entity *ImageEntity) error
}
