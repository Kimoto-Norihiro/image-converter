package imagemoduleclient

import (
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/adapter"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/app/imageusecase"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
) *imagemodule.ImageModule {
	imageRepo := imagerepository.New()
	imageUsecase := imageusecase.NewImageUsecase(
		db, imageRepo,
	)

	return &imagemodule.ImageModule{
		ImageUsecase: imageUsecase,
	}
}
