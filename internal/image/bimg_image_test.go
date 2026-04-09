package image

import (
	_ "embed"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.originalWidth, size.Width)
	assert.Equal(suite.T(), suite.originalHeight, size.Height)
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
			assert.NoError(suite.T(), err)
			assert.NotEmpty(suite.T(), processedImg)

			size, err := NewBimgImage(processedImg).Size()
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expectedWidth, size.Width)
			assert.Equal(suite.T(), tt.expectedHeight, size.Height)
		})
	}
}

func (suite *BimgImageTestSuite) TestConvert() {
	suite.Run("Convert to PNG", func() {
		convertedImg, err := suite.img.Convert("png")
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), convertedImg)
		assert.Equal(suite.T(), "png", NewBimgImage(convertedImg).ImageType())
	})

	suite.Run("Error on unsupported type", func() {
		convertedImg, err := suite.img.Convert("invalid_format")
		assert.ErrorIs(suite.T(), err, ErrUnsupportedImageType)
		assert.Empty(suite.T(), convertedImg)
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
				assert.Equal(suite.T(), ErrUnsupportedImageType, err)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expected, result)
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
			assert.Equal(suite.T(), tt.expected, result)
		})
	}
}

func (suite *BimgImageTestSuite) TestInvalidImageBuffer() {
	invalidBuffer := []byte("not an image")
	img := NewBimgImage(invalidBuffer)

	_, err := img.Size()
	assert.Error(suite.T(), err)

	metadata, err := img.Metadata()
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), metadata.Type)
}

func (suite *BimgImageTestSuite) TestFlip() {
	flippedImg, err := suite.img.Flip()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), flippedImg)

	originalSize, err := suite.img.Size()
	assert.NoError(suite.T(), err)
	flippedImgSize, err := NewBimgImage(flippedImg).Size()
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), originalSize.Width, flippedImgSize.Width)
	assert.Equal(suite.T(), originalSize.Height, flippedImgSize.Height)
}

func (suite *BimgImageTestSuite) TestLength() {
	assert.Equal(suite.T(), suite.originalLength, suite.img.Length())
}

func (suite *BimgImageTestSuite) TestGrayscale() {
	grayscaleImg, err := suite.img.Grayscale()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), grayscaleImg)
}

func (suite *BimgImageTestSuite) TestEnablePalette() {
	initialImgLength := suite.img.Length()

	paletteImg, err := suite.img.EnablePalette()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), paletteImg)

	paletteImgLength := NewBimgImage(paletteImg).Length()
	assert.NotZero(suite.T(), paletteImgLength)
	assert.NotEqual(suite.T(), initialImgLength, paletteImgLength, "Expected image data to change after applying palette")
}

func (suite *BimgImageTestSuite) TestLosslessCompress_Integrity() {
	compressedImg, err := suite.img.LosslessCompress()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), compressedImg)

	metadata, err := NewBimgImage(compressedImg).Metadata()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.originalWidth, metadata.Size.Width)
	assert.Equal(suite.T(), suite.originalHeight, metadata.Size.Height)
}

func (suite *BimgImageTestSuite) TestMetadata() {
	metadata, err := suite.img.Metadata()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), metadata)

	assert.Equal(suite.T(), suite.originalWidth, metadata.Size.Width)
	assert.Equal(suite.T(), suite.originalHeight, metadata.Size.Height)
	assert.Equal(suite.T(), "jpeg", metadata.Type)
}

func TestBimgImageTestSuite(t *testing.T) {
	suite.Run(t, new(BimgImageTestSuite))
}
