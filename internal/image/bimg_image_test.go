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
	converted, err := img.Convert(bimg.PNG)

	assert.Nil(t, err)
	assert.NotEmpty(t, converted)
	assert.Equal(t, "png", NewBimgImage(converted).ImageType())
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
