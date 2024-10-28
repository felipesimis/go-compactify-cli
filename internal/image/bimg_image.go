package image

import (
	"errors"

	"github.com/h2non/bimg"
)

var ErrUnsupportedImageType = errors.New("unsupported image type")

type BimgImage interface {
	Resize(width, height int) ([]byte, error)
	Size() (bimg.ImageSize, error)
	Convert(format string) ([]byte, error)
	ImageType() string
	Crop(width, height int, gravity bimg.Gravity) ([]byte, error)
	Flip() ([]byte, error)
	Enlarge(width, height int) ([]byte, error)
	Thumbnail(width int) ([]byte, error)
	ImageInterpretation() (bimg.Interpretation, error)
	Grayscale() ([]byte, error)
	Length() int
}

type BimgImageWrapper struct {
	image *bimg.Image
}

func NewBimgImage(buffer []byte) BimgImage {
	return &BimgImageWrapper{image: bimg.NewImage(buffer)}
}

func (b *BimgImageWrapper) Size() (bimg.ImageSize, error) {
	return b.image.Size()
}

func (b *BimgImageWrapper) Resize(width, height int) ([]byte, error) {
	return b.image.Resize(width, height)
}

func (b *BimgImageWrapper) Convert(format string) ([]byte, error) {
	bimgFormat, err := mapStringToImageType(format)
	if err != nil {
		return nil, err
	}
	return b.image.Convert(bimgFormat)
}

func (b *BimgImageWrapper) ImageType() string {
	return b.image.Type()
}

func (b *BimgImageWrapper) Crop(width, height int, gravity bimg.Gravity) ([]byte, error) {
	return b.image.Crop(width, height, gravity)
}

func mapStringToImageType(format string) (bimg.ImageType, error) {
	switch format {
	case "jpeg", "jpg":
		return bimg.JPEG, nil
	case "webp":
		return bimg.WEBP, nil
	case "png":
		return bimg.PNG, nil
	default:
		return bimg.UNKNOWN, ErrUnsupportedImageType
	}
}

func (b *BimgImageWrapper) Flip() ([]byte, error) {
	return b.image.Flip()
}

func (b *BimgImageWrapper) Enlarge(width, height int) ([]byte, error) {
	return b.image.Enlarge(width, height)
}

func (b *BimgImageWrapper) Thumbnail(width int) ([]byte, error) {
	return b.image.Thumbnail(width)
}

func (b *BimgImageWrapper) ImageInterpretation() (bimg.Interpretation, error) {
	return b.image.Interpretation()
}

func (b *BimgImageWrapper) Grayscale() ([]byte, error) {
	return b.image.Colourspace(bimg.InterpretationBW)
}

func (b *BimgImageWrapper) Length() int {
	return b.image.Length()
}
