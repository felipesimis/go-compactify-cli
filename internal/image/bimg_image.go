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

func InitializeProcessor() {
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)
}

func ShutdownProcessor() {
	bimg.Shutdown()
}

func (b *bimgImageWrapper) Size() (ImageSize, error) {
	size, err := b.image.Size()
	if err != nil {
		return ImageSize{}, err
	}
	return ImageSize{Width: size.Width, Height: size.Height}, nil
}

func (b *bimgImageWrapper) ImageType() string {
	return b.image.Type()
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

func (b *bimgImageWrapper) Length() int {
	return b.image.Length()
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

func (b *bimgImageWrapper) Process(opts ...ProcessOption) ([]byte, error) {
	domainOpts := &domainOptions{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(domainOpts)
	}

	bimgOpts := bimg.Options{}

	if domainOpts.thumbnailWidth > 0 {
		bimgOpts.Width = domainOpts.thumbnailWidth
		bimgOpts.Height = domainOpts.thumbnailWidth
		bimgOpts.Crop = true
		bimgOpts.Gravity = bimg.GravitySmart
	} else if domainOpts.width > 0 || domainOpts.height > 0 {
		bimgOpts.Width = domainOpts.width
		bimgOpts.Height = domainOpts.height
		bimgOpts.Enlarge = domainOpts.enlarge
	}

	if domainOpts.crop {
		bimgOpts.Crop = true
		bimgOpts.Gravity = mapGravityToBimg(domainOpts.gravity)
	}

	if domainOpts.format != "" {
		bimgFormat, err := mapStringToImageType(domainOpts.format)
		if err != nil {
			return nil, err
		}
		bimgOpts.Type = bimgFormat
	}

	if domainOpts.grayscale {
		bimgOpts.Interpretation = bimg.InterpretationBW
	}

	if domainOpts.flip {
		bimgOpts.Flip = true
	}

	if domainOpts.palette {
		bimgOpts.Palette = true
	}

	if domainOpts.lossless {
		bimgOpts.Lossless = true
	}

	if domainOpts.stripMetadata {
		bimgOpts.StripMetadata = true
	}

	return b.image.Process(bimgOpts)
}
