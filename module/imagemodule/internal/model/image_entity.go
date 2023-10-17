package imageentity

import (
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

type ImageEntity struct {
	image    imagemodel.Image
	original *imagemodel.Image
}

func NewImageEntityToCreate(
	objectName string, resizeWidthPercent int, resizeHeightPercent int, encodeFormat imagemodel.EncodeFormat,
) (*ImageEntity, error) {
	return &ImageEntity{
		image: imagemodel.Image{
			ObjectName: objectName,
			ResizeWidthPercent: resizeWidthPercent,
			ResizeHeightPercent: resizeHeightPercent,
			EncodeFormat: encodeFormat,
			Status: imagemodel.Waiting,
		},
	}, nil
}


func (e *ImageEntity) Image() *imagemodel.Image {
	return &e.image
}

func (ie *ImageEntity) SetConvertedImageURL(convertedImageURL string) {
	ie.image.ConvertedImageURL = convertedImageURL
}

func (ie *ImageEntity) SetStatus(status imagemodel.ImageStatus) {
	ie.image.Status = status
}