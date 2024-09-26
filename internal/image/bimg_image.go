package image

import "github.com/h2non/bimg"

type BimgImage interface {
	Resize(width, height int) ([]byte, error)
	Size() (bimg.ImageSize, error)
	Convert(format bimg.ImageType) ([]byte, error)
	ImageType() string
	Crop(width, height int, gravity bimg.Gravity) ([]byte, error)
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

func (b *BimgImageWrapper) Convert(format bimg.ImageType) ([]byte, error) {
	return b.image.Convert(format)
}

func (b *BimgImageWrapper) ImageType() string {
	return b.image.Type()
}

func (b *BimgImageWrapper) Crop(width, height int, gravity bimg.Gravity) ([]byte, error) {
	return b.image.Crop(width, height, gravity)
}
