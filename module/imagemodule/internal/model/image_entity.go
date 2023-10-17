package imageentity

import (
	"errors"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/mo"
)

type ImageEntity struct {
	image    imagemodel.Image
	original *imagemodel.Image
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

func NewImageEntityToCreate(
	objectName string, resizeWidthPercent int, resizeHeightPercent int, encodeFormat imagemodel.EncodeFormat,
) (*ImageEntity, error) {
	return &ImageEntity{
		image: imagemodel.Image{
			ObjectName:          objectName,
			ResizeWidthPercent:  resizeWidthPercent,
			ResizeHeightPercent: resizeHeightPercent,
			EncodeFormat:        encodeFormat,
			Status:              imagemodel.Waiting,
		},
	}, nil
}

func NewAsyncJobEntityToUpdate(asyncJob *imagemodel.Image) *ImageEntity {
	if asyncJob == nil {
		return nil
	}

	return &ImageEntity{
		image:    *asyncJob.DeepClone(),
		original: asyncJob,
	}
}

type ImageUpdatedValues struct {
	Status            mo.Option[imagemodel.ImageStatus]
	ConvertedImageURL mo.Option[string]
}

func (e *ImageEntity) UpdateAsyncJob(asyncJobToUpdate *imagemodel.ImageToUpdate) error {
	if v, ok := asyncJobToUpdate.Status.Get(); ok {
		e.image.Status = v
	}
	if v, ok := asyncJobToUpdate.ConvertedImageURL.Get(); ok {
		e.image.ConvertedImageURL = v
	}

	return nil
}

func (e *ImageEntity) UpdateValues() (*ImageUpdatedValues, error) {
	if e.original == nil {
		return nil, errors.New("new entity doesn't have updated values")
	}

	v := &ImageUpdatedValues{}

	if !cmp.Equal(e.image.Status, e.original.Status) {
		v.Status = mo.Some(e.image.Status)
	}
	if !cmp.Equal(e.image.ConvertedImageURL, e.original.ConvertedImageURL) {
		v.ConvertedImageURL = mo.Some(e.image.ConvertedImageURL)
	}

	return v, nil
}
