package image

import (
	"github.com/h2non/bimg"
)

type bimgImageWrapper struct {
	image *bimg.Image
}

func NewProcessor(buffer []byte) ImageProcessor {
	return &bimgImageWrapper{image: bimg.NewImage(buffer)}
}

func (b *bimgImageWrapper) Size() (ImageSize, error) {
	size, err := b.image.Size()
	if err != nil {
		return ImageSize{}, err
	}
	return ImageSize{Width: size.Width, Height: size.Height}, nil
}

func (b *bimgImageWrapper) Resize(width, height int) ([]byte, error) {
	return b.image.Resize(width, height)
}

func (b *bimgImageWrapper) Convert(format string) ([]byte, error) {
	bimgFormat, err := mapStringToImageType(format)
	if err != nil {
		return nil, err
	}
	return b.image.Convert(bimgFormat)
}

func (b *bimgImageWrapper) ImageType() string {
	return b.image.Type()
}

func (b *bimgImageWrapper) Crop(width, height int, gravity Gravity) ([]byte, error) {
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

func (b *bimgImageWrapper) Flip() ([]byte, error) {
	return b.image.Flip()
}

func (b *bimgImageWrapper) Enlarge(width, height int) ([]byte, error) {
	return b.image.Enlarge(width, height)
}

func (b *bimgImageWrapper) Thumbnail(width int) ([]byte, error) {
	return b.image.Thumbnail(width)
}

func (b *bimgImageWrapper) Grayscale() ([]byte, error) {
	return b.image.Colourspace(bimg.InterpretationBW)
}

func (b *bimgImageWrapper) Length() int {
	return b.image.Length()
}

func (b *bimgImageWrapper) EnablePalette() ([]byte, error) {
	return b.image.Process(bimg.Options{Palette: true})
}

func (b *bimgImageWrapper) LosslessCompress() ([]byte, error) {
	return b.image.Process(bimg.Options{Lossless: true})
}

func (b *bimgImageWrapper) Metadata() (ImageMetadata, error) {
	size, err := b.Size()
	if err != nil {
		return ImageMetadata{}, err
	}
	return ImageMetadata{
		Size: size,
		Type: b.ImageType(),
	}, nil
}
