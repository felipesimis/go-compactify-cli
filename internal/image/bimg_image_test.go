package image

import (
	"os"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func mockedImage() []byte {
	buffer, err := os.ReadFile("../../test/testdata/sample.jpeg")
	if err != nil {
		panic(err)
	}
	return buffer
}

func TestBimgImageWrapper_Size(t *testing.T) {
	img := NewBimgImage(mockedImage())
	size, err := img.Size()

	assert.Nil(t, err)
	assert.Equal(t, 600, size.Width)
	assert.Equal(t, 400, size.Height)
}

func TestBimgImageWrapper_Resize(t *testing.T) {
	img := NewBimgImage(mockedImage())
	resized, err := img.Resize(300, 200)

	assert.Nil(t, err)
	assert.NotEmpty(t, resized)

	size, err := NewBimgImage(resized).Size()
	assert.Nil(t, err)
	assert.Equal(t, 300, size.Width)
	assert.Equal(t, 200, size.Height)
}

func TestBimgImageWrapper_Convert(t *testing.T) {
	img := NewBimgImage(mockedImage())
	converted, err := img.Convert("png")

	assert.Nil(t, err)
	assert.NotEmpty(t, converted)
	assert.Equal(t, "png", NewBimgImage(converted).ImageType())

	converted, err = img.Convert("unknown")
	assert.Equal(t, ErrUnsupportedImageType, err)
	assert.Empty(t, converted)
}

func TestBimgImageWrapper_mapStringToImageType(t *testing.T) {
	imageType, err := mapStringToImageType("jpeg")
	assert.Nil(t, err)
	assert.Equal(t, bimg.JPEG, imageType)

	imageType, err = mapStringToImageType("jpg")
	assert.Nil(t, err)
	assert.Equal(t, bimg.JPEG, imageType)

	imageType, err = mapStringToImageType("webp")
	assert.Nil(t, err)
	assert.Equal(t, bimg.WEBP, imageType)

	imageType, err = mapStringToImageType("png")
	assert.Nil(t, err)
	assert.Equal(t, bimg.PNG, imageType)

	imageType, err = mapStringToImageType("unknown")
	assert.Equal(t, ErrUnsupportedImageType, err)
	assert.Equal(t, bimg.UNKNOWN, imageType)
}

func TestBimgImageWrapper_Crop(t *testing.T) {
	img := NewBimgImage(mockedImage())
	cropped, err := img.Crop(300, 200, bimg.GravitySmart)

	assert.Nil(t, err)
	assert.NotEmpty(t, cropped)

	size, err := NewBimgImage(cropped).Size()
	assert.Nil(t, err)
	assert.Equal(t, 300, size.Width)
	assert.Equal(t, 200, size.Height)
}

func TestBimgImageWrapper_Flip(t *testing.T) {
	img := NewBimgImage(mockedImage())
	flipped, err := img.Flip()
	assert.Nil(t, err)
	assert.NotEmpty(t, flipped)

	originalSize, err := img.Size()
	assert.Nil(t, err)
	flippedSize, err := NewBimgImage(flipped).Size()
	assert.Nil(t, err)

	assert.Equal(t, originalSize.Width, flippedSize.Width)
	assert.Equal(t, originalSize.Height, flippedSize.Height)
}

func TestBimgImageWrapper_Enlarge(t *testing.T) {
	img := NewBimgImage(mockedImage())
	enlarged, err := img.Enlarge(1200, 800)
	assert.Nil(t, err)
	assert.NotEmpty(t, enlarged)

	size, err := NewBimgImage(enlarged).Size()
	assert.Nil(t, err)
	assert.Equal(t, 1200, size.Width)
	assert.Equal(t, 800, size.Height)
}
