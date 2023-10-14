package imagerepository

import (
	"errors"

	"github.com/Kimoto-Norihiro/image-converter/module/imagemodule/model/imagemodel"
)

func ImageStatusFromStatusID(statusID int) (imagemodel.ImageStatus, error) {
	switch statusID {
	case 1:
		return imagemodel.Waiting, nil
	case 2:
		return imagemodel.Succeeded, nil
	case 3:
		return imagemodel.Failed, nil
	}
	return 0, errors.New("invalid status id")
}

func ImageStatusToStatusID(status imagemodel.ImageStatus) (int, error) {
	switch status {
	case imagemodel.Waiting:
		return 1, nil
	case imagemodel.Succeeded:
		return 2, nil
	case imagemodel.Failed:
		return 3, nil
	}
	return 0, errors.New("invalid status")
}
