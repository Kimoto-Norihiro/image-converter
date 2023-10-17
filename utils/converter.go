package utils

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

func Converter(fileName string, resizeFileName string, resizeFileWidth int, resizeFileHeight int, encodeFormat string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("os.Open error")
		log.Fatal(err)
	}
	defer file.Close()

	// デコード
	srcImg, _, err := image.Decode(file)
	if err != nil {
		log.Println("image.Decode error")
		fmt.Fprintln(os.Stderr, err)
	}

	srcRct := srcImg.Bounds()
	dstImg := image.NewRGBA(image.Rect(0, 0, srcRct.Dx()*resizeFileWidth/100, srcRct.Dy()*resizeFileHeight/100)) // 669KB -> 142KB
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcRct, draw.Over, nil)

	dst, err := os.Create(resizeFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	// エンコード
	switch encodeFormat {
	case "jpeg":
		err = jpeg.Encode(dst, dstImg, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatal(err)
		}
	case "png":
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
