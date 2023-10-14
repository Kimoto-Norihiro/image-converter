package imageentity

import "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"

type ImageEntity struct {
	image    imagemodel.Image
	original *imagemodel.Image
}

func (e *ImageEntity) Image() *imagemodel.Image {
	return &e.image
}