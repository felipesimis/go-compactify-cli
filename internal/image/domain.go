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

type domainOptions struct {
	width          int
	height         int
	format         string
	grayscale      bool
	crop           bool
	gravity        Gravity
	enlarge        bool
	palette        bool
	lossless       bool
	thumbnailWidth int
	flip           bool
}

type ProcessOption func(*domainOptions)

func WithResize(width, height int) ProcessOption {
	return func(o *domainOptions) {
		o.width = width
		o.height = height
	}
}

func WithConvert(format string) ProcessOption {
	return func(o *domainOptions) {
		o.format = format
	}
}

func WithCrop(width, height int, gravity Gravity) ProcessOption {
	return func(o *domainOptions) {
		o.crop = true
		o.width = width
		o.height = height
		o.gravity = gravity
	}
}

func WithGrayscale() ProcessOption {
	return func(o *domainOptions) {
		o.grayscale = true
	}
}

func WithFlip() ProcessOption {
	return func(o *domainOptions) {
		o.flip = true
	}
}

func WithEnlarge(width, height int) ProcessOption {
	return func(o *domainOptions) {
		o.enlarge = true
		o.width = width
		o.height = height
	}
}

func WithThumbnail(width int) ProcessOption {
	return func(o *domainOptions) {
		o.thumbnailWidth = width
	}
}

func WithPalette() ProcessOption {
	return func(o *domainOptions) {
		o.palette = true
	}
}

func WithLosslessCompress() ProcessOption {
	return func(o *domainOptions) {
		o.lossless = true
	}
}

type ImageProcessor interface {
	Size() (ImageSize, error)
	ImageType() string
	Length() int
	Metadata() (ImageMetadata, error)

	Process(opts ...ProcessOption) ([]byte, error)
}

type ProcessorFactory func(imageData []byte) ImageProcessor
