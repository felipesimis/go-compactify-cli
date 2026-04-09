package image

import (
	_ "embed"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/suite"
)

type BimgImageTestSuite struct {
	suite.Suite
	img            *BimgImageWrapper
	originalWidth  int
	originalHeight int
	originalLength int
}

//go:embed testdata/sample.jpeg
var tinyJPEG []byte

func (suite *BimgImageTestSuite) SetupTest() {
	suite.img = NewBimgImage(tinyJPEG)
	size, err := suite.img.Size()
	suite.Require().NoError(err)
	suite.originalWidth = size.Width
	suite.originalHeight = size.Height
	suite.originalLength = suite.img.Length()
	suite.Require().Greater(suite.originalLength, 0)
}

func (suite *BimgImageTestSuite) TestSize() {
	size, err := suite.img.Size()
	suite.NoError(err)
	suite.Equal(suite.originalWidth, size.Width)
	suite.Equal(suite.originalHeight, size.Height)
}

func (suite *BimgImageTestSuite) TestSizingOperations() {
	tests := []struct {
		name           string
		operation      func() ([]byte, error)
		expectedWidth  int
		expectedHeight int
	}{
		{
			name: "Resize",
			operation: func() ([]byte, error) {
				return suite.img.Resize(300, 200)
			},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name: "Crop",
			operation: func() ([]byte, error) {
				return suite.img.Crop(300, 200, GravitySmart)
			},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name: "Enlarge",
			operation: func() ([]byte, error) {
				return suite.img.Enlarge(1200, 800)
			},
			expectedWidth:  1200,
			expectedHeight: 800,
		},
		{
			name: "Thumbnail",
			operation: func() ([]byte, error) {
				return suite.img.Thumbnail(300)
			},
			expectedWidth:  300,
			expectedHeight: 300,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			processedImg, err := tt.operation()
			suite.NoError(err)
			suite.NotEmpty(processedImg)

			size, err := NewBimgImage(processedImg).Size()
			suite.NoError(err)
			suite.Equal(tt.expectedWidth, size.Width)
			suite.Equal(tt.expectedHeight, size.Height)
		})
	}
}

func (suite *BimgImageTestSuite) TestConvert() {
	suite.Run("Convert to PNG", func() {
		convertedImg, err := suite.img.Convert("png")
		suite.NoError(err)
		suite.NotEmpty(convertedImg)
		suite.Equal("png", NewBimgImage(convertedImg).ImageType())
	})

	suite.Run("Error on unsupported type", func() {
		convertedImg, err := suite.img.Convert("invalid_format")
		suite.ErrorIs(err, ErrUnsupportedImageType)
		suite.Empty(convertedImg)
	})
}

func (suite *BimgImageTestSuite) TestMapStringToImageType() {
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
		suite.Run(tt.name, func() {
			result, err := mapStringToImageType(tt.input)
			if tt.expected == bimg.UNKNOWN {
				suite.Equal(ErrUnsupportedImageType, err)
			} else {
				suite.NoError(err)
				suite.Equal(tt.expected, result)
			}
		})
	}
}

func (suite *BimgImageTestSuite) TestMapGravityToBimg() {
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
		suite.Run(tt.name, func() {
			result := mapGravityToBimg(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *BimgImageTestSuite) TestInvalidImageBuffer() {
	invalidBuffer := []byte("not an image")
	img := NewBimgImage(invalidBuffer)

	_, err := img.Size()
	suite.Error(err)

	metadata, err := img.Metadata()
	suite.Error(err)
	suite.Empty(metadata.Type)
}

func (suite *BimgImageTestSuite) TestFlip() {
	flippedImg, err := suite.img.Flip()
	suite.NoError(err)
	suite.NotEmpty(flippedImg)

	originalSize, err := suite.img.Size()
	suite.NoError(err)
	flippedImgSize, err := NewBimgImage(flippedImg).Size()
	suite.NoError(err)

	suite.Equal(originalSize.Width, flippedImgSize.Width)
	suite.Equal(originalSize.Height, flippedImgSize.Height)
}

func (suite *BimgImageTestSuite) TestLength() {
	suite.Equal(suite.originalLength, suite.img.Length())
}

func (suite *BimgImageTestSuite) TestGrayscale() {
	grayscaleImg, err := suite.img.Grayscale()
	suite.NoError(err)
	suite.NotEmpty(grayscaleImg)
}

func (suite *BimgImageTestSuite) TestEnablePalette() {
	initialImgLength := suite.img.Length()

	paletteImg, err := suite.img.EnablePalette()
	suite.NoError(err)
	suite.NotEmpty(paletteImg)

	paletteImgLength := NewBimgImage(paletteImg).Length()
	suite.NotZero(paletteImgLength)
	suite.NotEqual(initialImgLength, paletteImgLength, "Expected image data to change after applying palette")
}

func (suite *BimgImageTestSuite) TestLosslessCompress_Integrity() {
	compressedImg, err := suite.img.LosslessCompress()
	suite.NoError(err)
	suite.NotEmpty(compressedImg)

	metadata, err := NewBimgImage(compressedImg).Metadata()
	suite.NoError(err)
	suite.Equal(suite.originalWidth, metadata.Size.Width)
	suite.Equal(suite.originalHeight, metadata.Size.Height)
}

func (suite *BimgImageTestSuite) TestMetadata() {
	metadata, err := suite.img.Metadata()
	suite.NoError(err)
	suite.NotEmpty(metadata)
	suite.Equal(suite.originalWidth, metadata.Size.Width)
	suite.Equal(suite.originalHeight, metadata.Size.Height)
	suite.Equal("jpeg", metadata.Type)
}

func TestBimgImageTestSuite(t *testing.T) {
	suite.Run(t, new(BimgImageTestSuite))
}
