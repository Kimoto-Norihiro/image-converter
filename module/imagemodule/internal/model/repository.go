package imageentity

import (
	"context"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=${GOPACKAGE}mock/$GOFILE -package=${GOPACKAGE}mock

type ImageRepository interface {
	ListImages(ctx context.Context, db *gorm.DB) ([]imagemodel.Image, error)
	Update(ctx context.Context, tx *gorm.DB, entity *ImageEntity) error
	Create(ctx context.Context, tx *gorm.DB, entity *ImageEntity) error
	FindForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*ImageEntity, error)
}
