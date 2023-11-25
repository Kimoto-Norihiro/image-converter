package imagerepository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kimoto-Norihiro/image-converter/database"
	imageentity "github.com/Kimoto-Norihiro/image-converter/module/imagemodule/internal/model"
	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/mo"
)

func makeImageRows(images ...ImageDTO) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "object_name", "resize_width_percent", "resize_height_percent", "encode_format_id", "status_id", "converted_image_url",
	})
	for _, i := range images {
		rows.AddRow(i.ID, i.ObjectName, i.ResizeWidthPercent, i.ResizeHeightPercent, i.EncodeFormatID, i.StatusID, i.ConvertedImageURL)
	}

	return rows
}

func TestListImages(t *testing.T) {
	db, mock := database.NewDBMock()

	var imagedtos []ImageDTO
	for i := 0; i < 3; i++ {
		imagedtos = append(imagedtos, ImageDTO{
			ID:                  int64(i),
			ObjectName:          fmt.Sprintf("object_name_%d", i),
			ResizeWidthPercent:  100,
			ResizeHeightPercent: 100,
			EncodeFormatID:      1,
			StatusID:            1,
			ConvertedImageURL:   "",
		})
	}
	images, _ := ImagesFromDTO(imagedtos)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `images`")).WillReturnRows(makeImageRows(imagedtos...))

	tx := db.Begin()

	got, err := New().ListImages(context.Background(), tx)
	if !cmp.Equal(got, images) || err != nil {
		t.Errorf(`ListImages(ctx, tx) = %+v, %+v; want %+v, %+v`, got, err, images, nil)
	}
}

func TestUpdate(t *testing.T) {
	db, mock := database.NewDBMock()

	image := imagemodel.Image{
		ID:                  1,
		ObjectName:          "object_name",
		ResizeWidthPercent:  100,
		ResizeHeightPercent: 100,
		EncodeFormat:        1,
		Status:              1,
		ConvertedImageURL:   "",
	}

	entity := imageentity.NewImageEntityToUpdate(&image)

	err := entity.UpdateImage(&imagemodel.ImageToUpdate{
		Status:            mo.Some(imagemodel.ImageStatus(imagemodel.Succeeded)),
		ConvertedImageURL: mo.Some("https://example.com/converted.jpg"),
	})
	if err != nil {
		t.Fatal("failed to set test")
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `images` SET `converted_image_url`=?,`status_id`=? WHERE id = ?")).
		WithArgs(entity.Image().ConvertedImageURL, entity.Image().Status, entity.Image().ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err = New().Update(context.Background(), tx, entity)
	if err != nil {
		t.Errorf(`Update(ctx, tx, image) = %+v; want %+v`, err, nil)
	}
}

func TestCreate(t *testing.T) {
	db, mock := database.NewDBMock()

	image := imagemodel.Image{
		ObjectName:          "object_name",
		ResizeWidthPercent:  100,
		ResizeHeightPercent: 100,
		EncodeFormat:        imagemodel.JPEG,
		Status:              imagemodel.Waiting,
		ConvertedImageURL:   "",
	}

	entity := imageentity.NewImageEntityToCreate(
		image.ObjectName,
		image.ResizeWidthPercent,
		image.ResizeHeightPercent,
		image.EncodeFormat,
	)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `images` (`object_name`,`resize_width_percent`,`resize_height_percent`,`encode_format_id`,`status_id`,`converted_image_url`) VALUES (?,?,?,?,?,?)")).
		WithArgs(image.ObjectName, image.ResizeWidthPercent, image.ResizeHeightPercent, image.EncodeFormat, image.Status, image.ConvertedImageURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx := db.Begin()
	err := New().Create(context.Background(), tx, entity)
	if err != nil {
		t.Errorf(`Create(ctx, tx, image) = %+v; want %+v`, err, nil)
	}
}

func TestFindForUpdate(t *testing.T) {
	db, mock := database.NewDBMock()

	imagedto := ImageDTO{
		ID:                  1,
		ObjectName:          "object_name",
		ResizeWidthPercent:  100,
		ResizeHeightPercent: 100,
		EncodeFormatID:      1,
		StatusID:            1,
		ConvertedImageURL:   "",
	}
	image, _ := ImageFromDTO(&imagedto)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `images` WHERE id = ?")).
		WithArgs(imagedto.ID).
		WillReturnRows(makeImageRows(imagedto))

	tx := db.Begin()
	want := imageentity.NewImageEntityToUpdate(image)
	got, err := New().FindForUpdate(context.Background(), tx, imagedto.ID)
	if !cmp.Equal(got, want, cmp.AllowUnexported(imageentity.ImageEntity{})) || err != nil {
		t.Errorf(`FindForUpdate(ctx, tx, %+v) = %+v, %+v; want %+v, %+v`, image.ID, got, err, &image, nil)
	}
}
