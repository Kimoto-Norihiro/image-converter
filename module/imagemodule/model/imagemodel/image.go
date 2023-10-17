package imagemodel

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
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

func (i *Image) Converter(fileName string) {
	filepath := fmt.Sprintf("/img/%s", fileName)
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("os.Open error")
		log.Fatal(err)
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		log.Println("image.Decode error")
		fmt.Fprintln(os.Stderr, err)
	}

	srcRct := srcImg.Bounds()
	dstImg := image.NewRGBA(image.Rect(0, 0, srcRct.Dx()*i.ResizeWidthPercent/100, srcRct.Dy()*i.ResizeHeightPercent/100)) // 669KB -> 142KB
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcRct, draw.Over, nil)

	resizeFileName := fmt.Sprintf("resize-%s", fileName)
	resizeFilePath := fmt.Sprintf("/img/%s", resizeFileName)

	dst, err := os.Create(resizeFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	// エンコード
	switch i.EncodeFormat {
	case JPEG:
		err = jpeg.Encode(dst, dstImg, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatal(err)
		}
	case PNG:
		err = png.Encode(dst, dstImg)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 変換前のファイルを削除
	err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}

type ImageToUpdate struct {
	Status            mo.Option[ImageStatus]
	ConvertedImageURL mo.Option[string]
}
