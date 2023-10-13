package converter

import (
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/image/draw"
)

func Converter(file *os.File) {
	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	srcRct := srcImg.Bounds()
	dstImg := image.NewRGBA(image.Rect(0, 0, srcRct.Dx()/4, srcRct.Dy()/4)) // 669KB -> 142KB
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcRct, draw.Over, nil)

	dst, err := os.Create("./img/resize_sample.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()
}
