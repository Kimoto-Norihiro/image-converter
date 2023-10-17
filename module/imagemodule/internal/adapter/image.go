package imagerepository

import (
	"context"
	"errors"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"gorm.io/gorm"
)

func (r *repo) ListImages(ctx context.Context, db *gorm.DB) ([]imagemodel.Image, error) {
	return nil, errors.New("not implemented")
}

func (r *repo) Update(ctx context.Context, db *gorm.DB, entity *imageentity.ImageEntity) error {
	return errors.New("not implemented")
}

func (r *repo) Create(ctx context.Context, db *gorm.DB, entity *imageentity.ImageEntity) error {
	return errors.New("not implemented")
}

func (r *repo) FindForUpdate(ctx context.Context, tx *gorm.DB, id int) (*imageentity.ImageEntity, error) {
	return nil, errors.New("not implemented")
}