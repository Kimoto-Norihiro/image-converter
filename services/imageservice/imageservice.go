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
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	for _, image := range images {
		if image.Status == imagemodel.ImageStatus(imagemodel.Succeeded) {
			continue
		}

		err := s.gcs.GCSFunction.DownloadFile(ctx, image.ObjectName)
		if err != nil {
			return nil, fmt.Errorf("failed to download file: %w", err)
		}

		err = image.Convert()
		var status imagemodel.ImageStatus
		if err != nil {
			status = imagemodel.ImageStatus(imagemodel.Failed)
		} else {
			status = imagemodel.ImageStatus(imagemodel.Succeeded)
		}

		convertedImagePath := fmt.Sprintf("converted-%s", image.ObjectName)
		convertedImageURL, err := s.gcs.GCSFunction.UploadFile(ctx, convertedImagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file: %w", err)
		}

		err = s.imageModule.ImageUsecase.UpdateImage(ctx, image.ID, &status, convertedImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to update image: %w", err)
		}
	}
	return &imageservicepb.ConvertImagesResponse{}, nil
}

func (s *ImageService) ListImages(ctx context.Context, req *imageservicepb.ListImagesRequest) (*imageservicepb.ListImagesResponse, error) {
	images, err := s.imageModule.ImageUsecase.ListImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
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
	err := s.gcs.GCSFunction.UploadNonConvertedFile(ctx, reader, req.ObjectName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	err = s.imageModule.ImageUsecase.CreateImage(ctx, req.ObjectName, int(req.ResizeWidthPercent), int(req.ResizeHeightPercent), imagemodel.EncodeFormat(req.EncodeFormat))
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}
	return &imageservicepb.CreateImageResponse{}, nil
}
