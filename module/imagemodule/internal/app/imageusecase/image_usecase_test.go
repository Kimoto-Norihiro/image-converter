package imageusecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kimoto-Norihiro/image-converter/database"
	imageentity "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model/imageentitymock"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type imageUsecaseMockSet struct {
	db        *gorm.DB
	mock      sqlmock.Sqlmock
	imageRepo *imageentitymock.MockImageRepository
}

func makeImageUsecaseMockSet(ctrl *gomock.Controller) imageUsecaseMockSet {
	db, mock, err := database.NewDBMock()
	if err != nil {
		panic(err)
	}

	return imageUsecaseMockSet{
		db:        db,
		mock:      mock,
		imageRepo: imageentitymock.NewMockImageRepository(ctrl),
	}
}

func mustImageUsecase(m imageUsecaseMockSet) *ImageUsecase {
	return NewImageUsecase(
		m.db,
		m.imageRepo,
	)
}

func TestListImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageUsecaseMockSet(ctrl)

	var images []imagemodel.Image

	m.imageRepo.EXPECT().ListImages(gomock.Any(), gomock.Any()).Return(images, nil)

	ctx := context.Background()

	got, err := mustImageUsecase(m).ListImages(ctx)
	want := images
	var wantErr error = nil

	if !cmp.Equal(got, want) || err != wantErr {
		t.Errorf(`CreateAsyncJob(ctx) = %+v, %+v; want %+v, %+v`, got, err, want, wantErr)
	}
}

func TestListImages_errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageUsecaseMockSet(ctrl)

	wantErr := errors.New("failed to list images")
	m.imageRepo.EXPECT().ListImages(gomock.Any(), gomock.Any()).Return(nil, wantErr)

	ctx := context.Background()

	got, err := mustImageUsecase(m).ListImages(ctx)

	if !errors.Is(wantErr, err) {
		t.Errorf(`ListImage(ctx) = %+v, %+v; want %+v, %+v`, got, err, nil, wantErr)
	}
}

func TestCreateImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageUsecaseMockSet(ctrl)

	m.mock.ExpectBegin()
	m.imageRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	m.mock.ExpectCommit()

	ctx := context.Background()
	ObjectName := "test.jpg"
	ResizeWidthPercent := 50
	ResizeHeightPercent := 50
	EncodeFormat := imagemodel.EncodeFormat(imagemodel.JPEG)

	err := mustImageUsecase(m).CreateImage(ctx, ObjectName, ResizeWidthPercent, ResizeHeightPercent, EncodeFormat)
	var wantErr error = nil

	t.Logf("err: %+v", err)

	if err != wantErr {
		t.Errorf(`CreateImage(ctx, %+v, %+v, %+v, %+v) = %+v; want %+v`, ObjectName, ResizeWidthPercent, ResizeHeightPercent, EncodeFormat, err, nil)
	}
}

func TestCreateImage_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageUsecaseMockSet(ctrl)

	wantErr := errors.New("failed to create image")

	m.mock.ExpectBegin()
	m.imageRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(wantErr)

	ctx := context.Background()
	ObjectName := "test.jpg"
	ResizeWidthPercent := 50
	ResizeHeightPercent := 50
	EncodeFormat := imagemodel.EncodeFormat(imagemodel.JPEG)

	err := mustImageUsecase(m).CreateImage(ctx, ObjectName, ResizeWidthPercent, ResizeHeightPercent, EncodeFormat)

	if !errors.Is(wantErr, err) {
		t.Errorf(`CreateImage(ctx, %+v, %+v, %+v, %+v) = %+v; want %+v`, ObjectName, ResizeWidthPercent, ResizeHeightPercent, EncodeFormat, err, wantErr)
	}
}

func TestUpdateImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := makeImageUsecaseMockSet(ctrl)

	entity := imageentity.NewAsyncJobEntityToUpdate(&imagemodel.Image{
		ID:                  1,
		ObjectName:          "test.jpg",
		ResizeWidthPercent:  50,
		ResizeHeightPercent: 50,
		EncodeFormat:        imagemodel.EncodeFormat(imagemodel.JPEG),
		Status:              imagemodel.ImageStatus(imagemodel.Waiting),
		ConvertedImageURL:   "",
	})

	m.mock.ExpectBegin()
	m.imageRepo.EXPECT().FindForUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity, nil)
	m.imageRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	m.mock.ExpectCommit()

	ctx := context.Background()
	ID := int64(1)
	StatusID := imagemodel.ImageStatus(imagemodel.Succeeded)
	ConvertedImageURL := "https://example.com/converted.jpg"

	err := mustImageUsecase(m).UpdateImage(ctx, ID, &StatusID, &ConvertedImageURL)
	var wantErr error = nil

	if err != wantErr {
		t.Errorf(`UpdateImage(ctx, %+v, %+v, %+v) = %+v; want %+v`, ID, StatusID, ConvertedImageURL, err, nil)
	}
}

func TestUpdateImage_error(t *testing.T) {
	ctx := context.Background()
	ID := int64(1)
	StatusID := imagemodel.ImageStatus(imagemodel.Succeeded)
	ConvertedImageURL := "https://example.com/converted.jpg"

	FailFindForUpdateError := errors.New("failed to find image for update")
	FailUpdateError := errors.New("failed to update image")

	entity := imageentity.NewAsyncJobEntityToUpdate(&imagemodel.Image{
		ID:                  1,
		ObjectName:          "test.jpg",
		ResizeWidthPercent:  50,
		ResizeHeightPercent: 50,
		EncodeFormat:        imagemodel.EncodeFormat(imagemodel.JPEG),
		Status:              imagemodel.ImageStatus(imagemodel.Waiting),
		ConvertedImageURL:   "",
	})

	testcases := []struct {
		name    string
		setup   func(*imageUsecaseMockSet)
		wantErr error
	}{
		{
			name: "failed to find image",
			setup: func(m *imageUsecaseMockSet) {
				m.mock.ExpectBegin()
				m.imageRepo.EXPECT().FindForUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, FailFindForUpdateError)
			},
			wantErr: FailFindForUpdateError,
		},
		{
			name: "failed to update image",
			setup: func(m *imageUsecaseMockSet) {
				m.mock.ExpectBegin()
				m.imageRepo.EXPECT().FindForUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity, nil)
				m.imageRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(FailUpdateError)
			},
			wantErr: FailUpdateError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := makeImageUsecaseMockSet(ctrl)
			tc.setup(&m)
			err := mustImageUsecase(m).UpdateImage(ctx, ID, &StatusID, &ConvertedImageURL)
			if !errors.Is(tc.wantErr, err) {
				t.Errorf(`UpdateImage(ctx, %+v, %+v, %+v) = %+v; want %+v`, ID, StatusID, ConvertedImageURL, err, tc.wantErr)
			}
		})
	}
}
