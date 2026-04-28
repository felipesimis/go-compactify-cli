package image

import (
	"errors"
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

	maxGravity
)

func (g Gravity) IsValid() bool {
	return g >= 0 && g < maxGravity
}

type ImageSize struct {
	Width  int
	Height int
}

type ImageMetadata struct {
	Size ImageSize
	Type string
}

type ImageProcessor interface {
	Size() (ImageSize, error)
	Resize(width, height int) ([]byte, error)
	Convert(format string) ([]byte, error)
	ImageType() string
	Crop(width, height int, gravity Gravity) ([]byte, error)
	Flip() ([]byte, error)
	Enlarge(width, height int) ([]byte, error)
	Thumbnail(width int) ([]byte, error)
	Grayscale() ([]byte, error)
	Length() int
	EnablePalette() ([]byte, error)
	LosslessCompress() ([]byte, error)
	Metadata() (ImageMetadata, error)
}
