package imagerepository

import "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"

func New() *repo {
	return &repo{}
}

type repo struct{}

var _ imageentity.ImageRepository = (*repo)(nil)
