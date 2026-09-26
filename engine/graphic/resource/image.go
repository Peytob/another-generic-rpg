package resource

type ImageType int

const (
	ImageTypeUnknown = ImageType(iota)
	ImageTypeRGB
	ImageTypeRGBA
)

type Image struct {
	Name   string
	Type   ImageType
	Width  int
	Height int
	Data   []byte
}
