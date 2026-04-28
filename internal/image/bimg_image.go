package image

import (
	"errors"

	"github.com/h2non/bimg"
)

var ErrUnsupportedImageType = errors.New("unsupported image type")

type Gravity int

const (
	GravityCentre Gravity = iota
	GravityNorth
	GravityEast
	GravitySouth
	GravityWest
	GravitySmart
)

type ImageSize struct {
	Width  int
	Height int
}

type ImageMetadata struct {
	Size ImageSize
	Type string
}

type BimgImageWrapper struct {
	image *bimg.Image
}

func NewBimgImage(buffer []byte) *BimgImageWrapper {
	return &BimgImageWrapper{image: bimg.NewImage(buffer)}
}

func (b *BimgImageWrapper) Size() (ImageSize, error) {
	size, err := b.image.Size()
	if err != nil {
		return ImageSize{}, err
	}
	return ImageSize{Width: size.Width, Height: size.Height}, nil
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

func (b *BimgImageWrapper) Crop(width, height int, gravity Gravity) ([]byte, error) {
	return b.image.Crop(width, height, mapGravityToBimg(gravity))
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

func mapGravityToBimg(g Gravity) bimg.Gravity {
	switch g {
	case GravityCentre:
		return bimg.GravityCentre
	case GravityNorth:
		return bimg.GravityNorth
	case GravityEast:
		return bimg.GravityEast
	case GravitySouth:
		return bimg.GravitySouth
	case GravityWest:
		return bimg.GravityWest
	default:
		return bimg.GravitySmart
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

func (b *BimgImageWrapper) Grayscale() ([]byte, error) {
	return b.image.Colourspace(bimg.InterpretationBW)
}

func (b *BimgImageWrapper) Length() int {
	return b.image.Length()
}

func (b *BimgImageWrapper) EnablePalette() ([]byte, error) {
	return b.image.Process(bimg.Options{Palette: true})
}

func (b *BimgImageWrapper) LosslessCompress() ([]byte, error) {
	return b.image.Process(bimg.Options{Lossless: true})
}

func (b *BimgImageWrapper) Metadata() (ImageMetadata, error) {
	size, err := b.Size()
	if err != nil {
		return ImageMetadata{}, err
	}
	return ImageMetadata{
		Size: size,
		Type: b.ImageType(),
	}, nil
}
