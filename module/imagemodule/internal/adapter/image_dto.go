package imagerepository

import (
	"fmt"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

type ImageDTO struct {
	ID                  int64  `json:"id" gorm:"primaryKey;autoIncrement:true"`
	ObjectName          string `json:"object_name" gorm:"unique"`
	ResizeWidthPercent  int    `json:"resize_width_percent" gorm:"not null"`
	ResizeHeightPercent int    `json:"resize_height_percent" gorm:"not null"`
	EncodeFormatID      int    `json:"encode_format_id" gorm:"not null"`
	StatusID            int    `json:"status_id" gorm:"not null"`
	ConvertedImageURL   string `json:"converted_image_url"`
}

func (ImageDTO) TableName() string {
	return "images"
}

func ImageFromDTO(dto *ImageDTO) (*imagemodel.Image, error) {
	encodeFormat, err := ImageEncodingFormatFromID(dto.EncodeFormatID)
	if err != nil {
		return nil, fmt.Errorf("invalid encode format id %v", dto)
	}
	status, err := ImageStatusFromStatusID(dto.StatusID)
	if err != nil {
		return nil, err
	}

	model := &imagemodel.Image{
		ID:                  dto.ID,
		ObjectName:          dto.ObjectName,
		ResizeWidthPercent:  dto.ResizeWidthPercent,
		ResizeHeightPercent: dto.ResizeHeightPercent,
		EncodeFormat:        encodeFormat,
		Status:              status,
		ConvertedImageURL:   dto.ConvertedImageURL,
	}

	return model, nil
}

func ImagesFromDTO(dtos []ImageDTO) ([]imagemodel.Image, error) {
	if dtos == nil {
		return nil, nil
	}

	models := make([]imagemodel.Image, len(dtos))
	for i := range dtos {
		m, err := ImageFromDTO(&dtos[i])
		if err != nil {
			return nil, err
		}
		models[i] = *m
	}
	return models, nil
}

func DTOFromImage(model *imagemodel.Image) (*ImageDTO, error) {
	encodeFormatID, err := ImageEncodingFormatToID(model.EncodeFormat)
	if err != nil {
		return nil, err
	}
	statusID, err := ImageStatusToStatusID(model.Status)
	if err != nil {
		return nil, err
	}

	dto := &ImageDTO{
		ID:                  model.ID,
		ObjectName:          model.ObjectName,
		ResizeWidthPercent:  model.ResizeWidthPercent,
		ResizeHeightPercent: model.ResizeHeightPercent,
		EncodeFormatID:      encodeFormatID,
		StatusID:            statusID,
		ConvertedImageURL:   model.ConvertedImageURL,
	}

	return dto, nil
}

func DTOFromImages(models []imagemodel.Image) ([]ImageDTO, error) {
	if models == nil {
		return nil, nil
	}

	dtos := make([]ImageDTO, len(models))
	for i := range models {
		dto, err := DTOFromImage(&models[i])
		if err != nil {
			return nil, err
		}
		dtos[i] = *dto
	}
	return dtos, nil
}
