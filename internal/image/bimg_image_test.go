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
				return img.Crop(300, 200, GravitySmart)
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
	tests := []struct {
		name     string
		input    string
		expected bimg.ImageType
	}{
		{"JPEG", "jpeg", bimg.JPEG},
		{"JPG", "jpg", bimg.JPEG},
		{"WEBP", "webp", bimg.WEBP},
		{"PNG", "png", bimg.PNG},
		{"Unknown", "unknown", bimg.UNKNOWN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapStringToImageType(tt.input)
			if tt.expected == bimg.UNKNOWN {
				assert.Equal(t, ErrUnsupportedImageType, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBimgImageWrapper_mapGravityToBimg(t *testing.T) {
	tests := []struct {
		name     string
		input    Gravity
		expected bimg.Gravity
	}{
		{"GravityCentre", GravityCentre, bimg.GravityCentre},
		{"GravityNorth", GravityNorth, bimg.GravityNorth},
		{"GravityEast", GravityEast, bimg.GravityEast},
		{"GravitySouth", GravitySouth, bimg.GravitySouth},
		{"GravitySmart (default)", GravitySmart, bimg.GravitySmart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapGravityToBimg(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBimgImageWrapper_InvalidImageBuffer(t *testing.T) {
	invalidBuffer := []byte("not an image")
	img := NewBimgImage(invalidBuffer)

	size, err := img.Size()
	assert.NotNil(t, err)
	assert.Equal(t, 0, size.Width)
	assert.Equal(t, 0, size.Height)

	metadata, err := img.Metadata()
	assert.NotNil(t, err)
	assert.Equal(t, 0, metadata.Size.Width)
	assert.Empty(t, metadata.Type)
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

func TestBimgImageWrapper_Grayscale(t *testing.T) {
	img := NewBimgImage(mockedImage())
	grayscaleImg, err := img.Grayscale()
	assert.Nil(t, err)
	assert.NotEmpty(t, grayscaleImg)
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
