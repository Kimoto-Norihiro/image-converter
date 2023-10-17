package imageservice

import (
	"context"
	"fmt"
	"log"

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
	for i, image := range images {
		log.Printf("image[%d]: %v", i, image)
	}

	return &imageservicepb.ListImagesResponse{}, nil
}

// func (s *ImageService) CreateImage(ctx context.Context, req *imageservicepb.CreateImageRequest) (*imageservicepb.CreateImageResponse, error) {
// 	return nil, nil
// }
