package imageservice

import (
	"context"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/imagemoduleclient"
	imageservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
	"gorm.io/gorm"
)

type ImageService struct {
	imageservicepb.UnimplementedImageServiceServer
	imageModule *imagemodule.ImageModule
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		imageModule: imagemoduleclient.New(db),
	}
}

func (s *ImageService) ConvertImages(ctx context.Context, req *imageservicepb.ConvertImagesRequest) (*imageservicepb.ConvertImagesResponse, error) {
	return nil, nil
}

// func (s *ImageService) ListImages(ctx context.Context, req *imageservicepb.ListImagesRequest) (*imageservicepb.ListImagesResponse, error) {
// 	return nil, nil
// }

// func (s *ImageService) CreateImage(ctx context.Context, req *imageservicepb.CreateImageRequest) (*imageservicepb.CreateImageResponse, error) {
// 	return nil, nil
// }

