package image

import "github.com/h2non/bimg"

type BimgImage interface {
	Resize(width, height int) ([]byte, error)
	Size() (bimg.ImageSize, error)
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
