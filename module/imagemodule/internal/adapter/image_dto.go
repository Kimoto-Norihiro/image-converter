package imagerepository

import "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"

type ImageDTO struct {
	ID                  int64  `json:"id"`
	ObjectName          string `json:"object_name"`
	ResizeWidthPercent  int    `json:"resize_width_percent"`
	ResizeHeightPercent int    `json:"resize_height_percent"`
	EncodeFormatID      int    `json:"encode_format_id"`
	StatusID            int    `json:"status_id"`
}

func (ImageDTO) TableName() string {
	return "images"
}

func ImageFromDTO(dto *ImageDTO) (*imagemodel.Image, error) {
	encodeFormat, err := ImageEncodingFormatFromID(dto.EncodeFormatID)
	if err != nil {
		return nil, err
	}
	status, err := ImageStatusFromStatusID(dto.StatusID)
	if err != nil {
		return nil, err
	}

	model := &imagemodel.Image{
		ID:        dto.ID,
		ObjectName: dto.ObjectName,
		ResizeWidthPercent:  dto.ResizeWidthPercent,
		ResizeHeightPercent: dto.ResizeHeightPercent,
		EncodeFormat:        encodeFormat,
		Status:              status,
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