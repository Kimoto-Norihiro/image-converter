package imagemodel

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
	ID int64
	ObjectName string
	ResizeWidthPercent int
	ResizeHeightPercent int
	EncodeFormat EncodeFormat
	Status ImageStatus
}