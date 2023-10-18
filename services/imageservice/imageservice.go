package imageservice

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Kimoto-Norihiro/image-converter/gcs"
	gcsclient "github.com/Kimoto-Norihiro/image-converter/gcs/gcs_client"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/imagemoduleclient"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	imageservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"

	"gorm.io/gorm"
)

type ImageService struct {
	imageservicepb.UnimplementedImageServiceServer
	imageModule *imagemodule.ImageModule
	gcs         *gcs.GCS
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		imageModule: imagemoduleclient.New(db),
		gcs:         gcsclient.New(),
	}
}

func (s *ImageService) ConvertImages(ctx context.Context, req *imageservicepb.ConvertImagesRequest) (*imageservicepb.ConvertImagesResponse, error) {
	images, err := s.imageModule.ImageUsecase.ListImages(ctx)
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		err := s.gcs.GCSUsecase.DownloadFile(ctx, image.ObjectName)
		if err != nil {
			return nil, err
		}
		filePath := fmt.Sprintf("/tmp/%s", image.ObjectName)
		err = image.Converter(filePath)
		if err != nil {
			return nil, err
		}

		convertedFilePath := fmt.Sprintf("/tmp/converted-%s", image.ObjectName)
		err = s.gcs.GCSUsecase.UploadFile(ctx, convertedFilePath)
		if err != nil {
			return nil, err
		}

		status := imagemodel.ImageStatus(imagemodel.Succeeded)
		err = s.imageModule.ImageUsecase.UpdateImage(ctx, image.ID, &status, &convertedFilePath)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *ImageService) ListImages(ctx context.Context, req *imageservicepb.ListImagesRequest) (*imageservicepb.ListImagesResponse, error) {
	images, err := s.imageModule.ImageUsecase.ListImages(ctx)
	if err != nil {
		return nil, err
	}

	var imagespb []*imageservicepb.Image
	for _, image := range images {
		imagespb = append(imagespb, &imageservicepb.Image{
			Id:                  image.ID,
			ObjectName:          image.ObjectName,
			ResizeWidthPercent:  int32(image.ResizeWidthPercent),
			ResizeHeightPercent: int32(image.ResizeHeightPercent),
			EncodeFormat:        imageservicepb.EncodeFormat(image.EncodeFormat),
			Status:              imageservicepb.ImageStatus(image.Status),
			ConvertedImageUrl:   image.ConvertedImageURL,
		})
	}

	return &imageservicepb.ListImagesResponse{
		Images: imagespb,
	}, nil
}

func (s *ImageService) CreateImage(ctx context.Context, req *imageservicepb.CreateImageRequest) (*imageservicepb.CreateImageResponse, error) {
	reader := bytes.NewReader(req.ImageFile)
	err := s.gcs.GCSUsecase.UploadNonConvertedFile(ctx, reader, req.ObjectName)
	if err != nil {
		return nil, err
	}
	err = s.imageModule.ImageUsecase.CreateImage(ctx, req.ObjectName, int(req.ResizeWidthPercent), int(req.ResizeHeightPercent), imagemodel.EncodeFormat(req.EncodeFormat))
	if err != nil {
		return nil, err
	}

	return &imageservicepb.CreateImageResponse{}, nil
}
