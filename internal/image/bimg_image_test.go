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

func TestBimgImageWrapper_SizingOperations(t *testing.T) {
	img := NewBimgImage(mockedImage())

	tests := []struct {
		name           string
		operation      func() ([]byte, error)
		expectedWidth  int
		expectedHeight int
	}{
		{
			name: "Resize",
			operation: func() ([]byte, error) {
				return img.Resize(300, 200)
			},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name: "Crop",
			operation: func() ([]byte, error) {
				return img.Crop(300, 200, bimg.GravitySmart)
			},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name: "Enlarge",
			operation: func() ([]byte, error) {
				return img.Enlarge(1200, 800)
			},
			expectedWidth:  1200,
			expectedHeight: 800,
		},
		{
			name: "Thumbnail",
			operation: func() ([]byte, error) {
				return img.Thumbnail(300)
			},
			expectedWidth:  300,
			expectedHeight: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedImg, err := tt.operation()
			assert.Nil(t, err)
			assert.NotEmpty(t, processedImg)
			size, err := NewBimgImage(processedImg).Size()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedWidth, size.Width)
			assert.Equal(t, tt.expectedHeight, size.Height)
		})
	}
}

func TestBimgImageWrapper_Convert(t *testing.T) {
	img := NewBimgImage(mockedImage())
	convertedImg, err := img.Convert("png")
	assert.Nil(t, err)
	assert.NotEmpty(t, convertedImg)
	assert.Equal(t, "png", NewBimgImage(convertedImg).ImageType())

	convertedImg, err = img.Convert("unknown")
	assert.Equal(t, ErrUnsupportedImageType, err)
	assert.Empty(t, convertedImg)
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

func TestBimgImageWrapper_Flip(t *testing.T) {
	img := NewBimgImage(mockedImage())
	flippedImg, err := img.Flip()
	assert.Nil(t, err)
	assert.NotEmpty(t, flippedImg)

	originalSize, err := img.Size()
	assert.Nil(t, err)
	flippedImgSize, err := NewBimgImage(flippedImg).Size()
	assert.Nil(t, err)

	assert.Equal(t, originalSize.Width, flippedImgSize.Width)
	assert.Equal(t, originalSize.Height, flippedImgSize.Height)
}

func TestBimgImageWrapper_Length(t *testing.T) {
	img := NewBimgImage(mockedImage())
	assert.Equal(t, 3773, img.Length())
}

func TestBimgImageWrapper_ImageInterpretation(t *testing.T) {
	img := NewBimgImage(mockedImage())
	interpretation, err := img.ImageInterpretation()
	assert.Nil(t, err)
	assert.Equal(t, bimg.InterpretationSRGB, interpretation)
}

func TestBimgImageWrapper_Grayscale(t *testing.T) {
	img := NewBimgImage(mockedImage())
	initialInterpretation, err := img.ImageInterpretation()
	assert.Nil(t, err)
	assert.Equal(t, bimg.InterpretationSRGB, initialInterpretation)

	grayscaleImg, err := img.Grayscale()
	assert.Nil(t, err)
	assert.NotEmpty(t, grayscaleImg)

	grayscaleImgInterpretation, err := NewBimgImage(grayscaleImg).ImageInterpretation()
	assert.Nil(t, err)
	assert.Equal(t, bimg.InterpretationBW, grayscaleImgInterpretation)
}

func TestBimgImageWrapper_EnablePalette(t *testing.T) {
	img := NewBimgImage(mockedImage())
	initialImgLength := img.Length()

	paletteImg, err := img.EnablePalette()
	assert.Nil(t, err)
	assert.NotEmpty(t, paletteImg)

	paletteImgLength := NewBimgImage(paletteImg).Length()
	assert.NotZero(t, paletteImgLength)
	assert.NotEqual(t, initialImgLength, paletteImgLength, "Expected image data to change after applying palette")
}

func TestBimgImageWrapper_LosslessCompress(t *testing.T) {
	img := NewBimgImage(mockedImage())
	compressedImg, err := img.LosslessCompress()
	assert.Nil(t, err)
	assert.NotEmpty(t, compressedImg)

	compressedImgLength := NewBimgImage(compressedImg).Length()
	assert.NotZero(t, compressedImgLength)
}

func TestBimgImageWrapper_Metadata(t *testing.T) {
	img := NewBimgImage(mockedImage())
	metadata, err := img.Metadata()
	assert.Nil(t, err)
	assert.NotEmpty(t, metadata)

	assert.Equal(t, 600, metadata.Size.Width)
	assert.Equal(t, 400, metadata.Size.Height)
	assert.Equal(t, "jpeg", metadata.Type)
}
