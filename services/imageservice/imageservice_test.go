package imageservice

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Kimoto-Norihiro/image-converter/gcs"
	"github.com/Kimoto-Norihiro/image-converter/gcs/gcsmock"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/imagemodulemock"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"gorm.io/gorm"

	imageservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
)

type imageServiceMockSet struct {
	imageUsecase *imagemodulemock.MockImageUsecase
	gcsFunction  *gcsmock.MockGCSFunction
	db           *gorm.DB
}

func makeImageServiceMockSet(ctrl *gomock.Controller) imageServiceMockSet {
	return imageServiceMockSet{
		imageUsecase: imagemodulemock.NewMockImageUsecase(ctrl),
		gcsFunction:  gcsmock.NewMockGCSFunction(ctrl),
	}
}

func mustImageService(m imageServiceMockSet) *ImageService {
	return &ImageService{
		imageModule: &imagemodule.ImageModule{
			ImageUsecase: m.imageUsecase,
		},
		gcs: &gcs.GCS{
			GCSFunction: m.gcsFunction,
		},
	}
}

func dummyCreateImageRequest() *grpc.CreateImageRequest {
	return &grpc.CreateImageRequest{
		ObjectName:          "test.jpg",
		ResizeWidthPercent:  50,
		ResizeHeightPercent: 50,
		EncodeFormat:        grpc.EncodeFormat_JPEG,
	}
}

func TestCreateImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageServiceMockSet(ctrl)

	m.gcsFunction.EXPECT().UploadNonConvertedFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	m.imageUsecase.EXPECT().CreateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	ctx := context.Background()
	req := dummyCreateImageRequest()

	want := &grpc.CreateImageResponse{}
	got, err := mustImageService(m).CreateImage(ctx, req)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" || err != nil {
		t.Errorf(`CreateImage(ctx, req) = %+v, %+v; want %+v, %+v`, got, err, want, nil)
	}
}

func TestCreateImage_errors(t *testing.T) {
	ErrFailedToUploadFile := errors.New("failed to upload file")
	ErrFailedToCreateImage := errors.New("failed to create image")

	ctx := context.Background()
	req := dummyCreateImageRequest()

	testcases := []struct {
		name  string
		setup func(m *imageServiceMockSet)
		err   error
	}{
		{
			name: "failed to upload file",
			setup: func(m *imageServiceMockSet) {
				m.gcsFunction.EXPECT().UploadNonConvertedFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrFailedToUploadFile)
			},
			err: ErrFailedToUploadFile,
		},
		{
			name: "failed to create image",
			setup: func(m *imageServiceMockSet) {
				m.gcsFunction.EXPECT().UploadNonConvertedFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.imageUsecase.EXPECT().CreateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrFailedToCreateImage)
			},
			err: ErrFailedToCreateImage,
		},
	}

	for _, tc := range testcases {
		ctrl := gomock.NewController(t)
		m := makeImageServiceMockSet(ctrl)
		tc.setup(&m)

		_, err := mustImageService(m).CreateImage(ctx, req)
		if !errors.Is(tc.err, err) {
			t.Errorf(`CreateImage(ctx, req) = %+v, %+v; want %+v, %+v`, nil, err, nil, tc.err)
		}
	}
}

func TestListImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageServiceMockSet(ctrl)

	images := []imagemodel.Image{}
	for i := 0; i < 3; i++ {
		images = append(images, imagemodel.Image{
			ID:                  int64(i),
			ObjectName:          fmt.Sprintf("object_name_%d", i),
			ResizeWidthPercent:  100,
			ResizeHeightPercent: 100,
			EncodeFormat:        imagemodel.EncodeFormat(imagemodel.JPEG),
			Status:              imagemodel.ImageStatus(imagemodel.Succeeded),
			ConvertedImageURL:   fmt.Sprintf("https://example.com/converted_%d.jpg", i),
		})
	}

	imagepbs := imagemodel.ImageModelsToImageServicePbs(images)
	m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(images, nil)

	ctx := context.Background()
	req := &grpc.ListImagesRequest{}

	want := &imageservicepb.ListImagesResponse{
		Images: imagepbs,
	}
	got, err := mustImageService(m).ListImages(ctx, req)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" || err != nil {
		t.Errorf(`ListImages(ctx, req) = %+v, %+v; want %+v, %+v`, got, err, want, nil)
	}
}

func TestListImages_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageServiceMockSet(ctrl)

	ErrFailedToListImages := errors.New("failed to list images")
	m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(nil, ErrFailedToListImages)

	ctx := context.Background()
	req := &grpc.ListImagesRequest{}

	_, err := mustImageService(m).ListImages(ctx, req)
	if !errors.Is(ErrFailedToListImages, err) {
		t.Errorf(`ListImages(ctx, req) = %+v, %+v; want %+v, %+v`, nil, err, nil, ErrFailedToListImages)
	}
}

func TestConvertImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageServiceMockSet(ctrl)

	images := []imagemodel.Image{
		{
			ID:                  1,
			ObjectName:          "test.jpg",
			ResizeWidthPercent:  50,
			ResizeHeightPercent: 50,
			EncodeFormat:        imagemodel.EncodeFormat(imagemodel.JPEG),
			Status:              imagemodel.ImageStatus(imagemodel.Waiting),
			ConvertedImageURL:   "",
		},
	}
	convertedImagePath := "converted-test.jpg"
	m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(images, nil)
	m.gcsFunction.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(nil)
	m.gcsFunction.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(&convertedImagePath, nil)
	m.imageUsecase.EXPECT().UpdateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	ctx := context.Background()
	req := &grpc.ConvertImagesRequest{}

	want := &imageservicepb.ConvertImagesResponse{}
	got, err := mustImageService(m).ConvertImages(ctx, req)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" || err != nil {
		t.Errorf(`ConvertImages(ctx, req) = %+v, %+v; want %+v, %+v`, got, err, want, nil)
	}
}

func TestConvertImages_error(t *testing.T) {
	ErrFailedToListImages := errors.New("failed to list images")
	ErrFailedToDownloadFile := errors.New("failed to download file")
	ErrFailedToUploadFile := errors.New("failed to upload file")
	ErrFailedToUpdateImage := errors.New("failed to download file")

	images := []imagemodel.Image{
		{
			ID:                  1,
			ObjectName:          "test.jpg",
			ResizeWidthPercent:  50,
			ResizeHeightPercent: 50,
			EncodeFormat:        imagemodel.EncodeFormat(imagemodel.JPEG),
			Status:              imagemodel.ImageStatus(imagemodel.Waiting),
			ConvertedImageURL:   "",
		},
	}
	convertedImagePath := "converted-test.jpg"

	testcases := []struct {
		name  string
		setup func(m *imageServiceMockSet)
		err   error
	}{
		{
			name: "failed to list images",
			setup: func(m *imageServiceMockSet) {
				m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(nil, ErrFailedToListImages)
			},
			err: ErrFailedToListImages,
		},
		{
			name: "failed to download file",
			setup: func(m *imageServiceMockSet) {
				m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(images, nil)
				m.gcsFunction.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(ErrFailedToDownloadFile)
			},
			err: ErrFailedToDownloadFile,
		},
		{
			name: "failed to upload file",
			setup: func(m *imageServiceMockSet) {
				m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(images, nil)
				m.gcsFunction.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(nil)
				m.gcsFunction.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil, ErrFailedToUploadFile)
			},
			err: ErrFailedToUploadFile,
		},
		{
			name: "failed to update image",
			setup: func(m *imageServiceMockSet) {
				m.imageUsecase.EXPECT().ListImages(gomock.Any()).Return(images, nil)
				m.gcsFunction.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(nil)
				m.gcsFunction.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(&convertedImagePath, nil)
				m.imageUsecase.EXPECT().UpdateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrFailedToUpdateImage)
			},
			err: ErrFailedToUpdateImage,
		},
	}

	for _, tc := range testcases {
		ctrl := gomock.NewController(t)
		m := makeImageServiceMockSet(ctrl)
		tc.setup(&m)

		ctx := context.Background()
		req := &grpc.ConvertImagesRequest{}

		_, err := mustImageService(m).ConvertImages(ctx, req)
		if !errors.Is(tc.err, err) {
			t.Errorf(`ConvertImages(ctx, req) = %+v, %+v; want %+v, %+v`, nil, err, nil, tc.err)
		}
	}
}