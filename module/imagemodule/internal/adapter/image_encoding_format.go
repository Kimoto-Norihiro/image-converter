package imagerepository

import (
	"fmt"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

func ImageEncodingFormatFromID(encodeFormatID int) (imagemodel.EncodeFormat, error) {
	switch encodeFormatID {
	case 1:
		return imagemodel.JPEG, nil
	case 2:
		return imagemodel.PNG, nil
	}
	return 0, fmt.Errorf("invalid encode format id %d", encodeFormatID)
}

func ImageEncodingFormatToID(encodeFormat imagemodel.EncodeFormat) (int, error) {
	switch encodeFormat {
	case imagemodel.JPEG:
		return 1, nil
	case imagemodel.PNG:
		return 2, nil
	}
	return 0, fmt.Errorf("invalid encode format %d", encodeFormat)
}
