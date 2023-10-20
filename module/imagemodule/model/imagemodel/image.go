package imagemodel

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/samber/mo"
	"golang.org/x/image/draw"
)

type ImageStatus int

const (
	Waiting ImageStatus = iota + 1
	Succeeded
	Failed
)

type EncodeFormat int

const (
	JPEG EncodeFormat = iota + 1
	PNG
)

type Image struct {
	ID                  int64
	ObjectName          string
	ResizeWidthPercent  int
	ResizeHeightPercent int
	EncodeFormat        EncodeFormat
	Status              ImageStatus
	ConvertedImageURL   string
}

func (i *Image) DeepClone() *Image {
	if i == nil {
		return nil
	}

	cloned := *i
	return &cloned
}

func (i *Image) Convert() error {
	filepath := fmt.Sprintf("./img/%s", i.ObjectName)
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	srcRct := srcImg.Bounds()
	dstImg := image.NewRGBA(image.Rect(0, 0, srcRct.Dx()*i.ResizeWidthPercent/100, srcRct.Dy()*i.ResizeHeightPercent/100)) // 669KB -> 142KB
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcRct, draw.Over, nil)

	resizeFileName := fmt.Sprintf("converted-%s", i.ObjectName)
	resizeFilePath := fmt.Sprintf("./img/%s", resizeFileName)

	dst, err := os.Create(resizeFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// エンコード
	switch i.EncodeFormat {
	case JPEG:
		err = jpeg.Encode(dst, dstImg, &jpeg.Options{Quality: 100})
		if err != nil {
			return fmt.Errorf("failed to encode image: %w", err)
		}
	case PNG:
		err = png.Encode(dst, dstImg)
		if err != nil {
			return fmt.Errorf("failed to encode image: %w", err)
		}
	}

	// 変換前のファイルを削除
	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

type ImageToUpdate struct {
	Status            mo.Option[ImageStatus]
	ConvertedImageURL mo.Option[string]
}
