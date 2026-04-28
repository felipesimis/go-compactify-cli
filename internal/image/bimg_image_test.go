package image

import (
	_ "embed"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/suite"
)

type BimgImageTestSuite struct {
	suite.Suite
	img            ImageProcessor
	originalWidth  int
	originalHeight int
	originalLength int
}

//go:embed testdata/sample.jpeg
var tinyJPEG []byte

func (suite *BimgImageTestSuite) SetupTest() {
	suite.img = NewProcessor(tinyJPEG)
	size, err := suite.img.Size()
	suite.Require().NoError(err)
	suite.originalWidth = size.Width
	suite.originalHeight = size.Height
	suite.originalLength = suite.img.Length()
	suite.Require().Greater(suite.originalLength, 0)
}

func (suite *BimgImageTestSuite) TestSize_ShouldReturnCorrectDimensions_WhenImageIsValid() {
	size, err := suite.img.Size()
	suite.NoError(err)
	suite.Equal(suite.originalWidth, size.Width)
	suite.Equal(suite.originalHeight, size.Height)
}

func (suite *BimgImageTestSuite) TestSizingOperations_ShouldTransformImageCorrectly_WhenOperationIsCalled() {
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

			size, err := NewProcessor(processedImg).Size()
			suite.NoError(err)
			suite.Equal(tt.expectedWidth, size.Width)
			suite.Equal(tt.expectedHeight, size.Height)
		})
	}
}

func (suite *BimgImageTestSuite) TestConvert_ShouldChangeImageType_WhenValidFormatIsProvided() {
	convertedImg, err := suite.img.Convert("png")
	suite.NoError(err)
	suite.NotEmpty(convertedImg)
	suite.Equal("png", NewProcessor(convertedImg).ImageType())
}

func (suite *BimgImageTestSuite) TestConvert_ShouldReturnError_WhenUnsupportedFormatIsProvided() {
	convertedImg, err := suite.img.Convert("invalid_format")
	suite.ErrorIs(err, ErrUnsupportedImageType)
	suite.Empty(convertedImg)
}

func (suite *BimgImageTestSuite) TestMapStringToImageType_ShouldReturnCorrectBimgType_WhenInputIsValid() {
	tests := []struct {
		name     string
		input    string
		expected bimg.ImageType
	}{
		{"JPEG", "jpeg", bimg.JPEG},
		{"JPG", "jpg", bimg.JPEG},
		{"WEBP", "webp", bimg.WEBP},
		{"PNG", "png", bimg.PNG},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := mapStringToImageType(tt.input)
			suite.NoError(err)
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *BimgImageTestSuite) TestMapStringToImageType_ShouldReturnError_WhenInputIsInvalid() {
	result, err := mapStringToImageType("unknown")
	suite.ErrorIs(err, ErrUnsupportedImageType)
	suite.Equal(bimg.UNKNOWN, result)
}

func (suite *BimgImageTestSuite) TestMapGravityToBimg_ShouldMapCorrectlyAndFallbackToSmartOnUnknown() {
	tests := []struct {
		name     string
		input    Gravity
		expected bimg.Gravity
	}{
		{"GravityCentre", GravityCentre, bimg.GravityCentre},
		{"GravityNorth", GravityNorth, bimg.GravityNorth},
		{"GravityEast", GravityEast, bimg.GravityEast},
		{"GravitySouth", GravitySouth, bimg.GravitySouth},
		{"GravityWest", GravityWest, bimg.GravityWest},
		{"GravitySmart", GravitySmart, bimg.GravitySmart},
		{"Fallback on Unknown Gravity", Gravity(999), bimg.GravitySmart},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := mapGravityToBimg(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *BimgImageTestSuite) TestInvalidImageBuffer_ShouldReturnError_WhenBufferIsNotAnImage() {
	invalidBuffer := []byte("not an image")
	img := NewProcessor(invalidBuffer)

	_, err := img.Size()
	suite.Error(err)

	metadata, err := img.Metadata()
	suite.Error(err)
	suite.Empty(metadata.Type)
}

func (suite *BimgImageTestSuite) TestFlip_ShouldMaintainDimensions_WhenImageIsFlipped() {
	flippedImg, err := suite.img.Flip()
	suite.NoError(err)
	suite.NotEmpty(flippedImg)

	originalSize, err := suite.img.Size()
	suite.NoError(err)
	flippedImgSize, err := NewProcessor(flippedImg).Size()
	suite.NoError(err)

	suite.Equal(originalSize.Width, flippedImgSize.Width)
	suite.Equal(originalSize.Height, flippedImgSize.Height)
}

func (suite *BimgImageTestSuite) TestLength_ShouldReturnCorrectByteLength() {
	suite.Equal(suite.originalLength, suite.img.Length())
}

func (suite *BimgImageTestSuite) TestGrayscale_ShouldReturnProcessedImage_WhenCalled() {
	grayscaleImg, err := suite.img.Grayscale()
	suite.NoError(err)
	suite.NotEmpty(grayscaleImg)
}

func (suite *BimgImageTestSuite) TestEnablePalette_ShouldChangeImageLength_WhenPaletteIsApplied() {
	initialImgLength := suite.img.Length()

	paletteImg, err := suite.img.EnablePalette()
	suite.NoError(err)
	suite.NotEmpty(paletteImg)

	paletteImgLength := NewProcessor(paletteImg).Length()
	suite.NotZero(paletteImgLength)
	suite.NotEqual(initialImgLength, paletteImgLength, "Expected image data to change after applying palette")
}

func (suite *BimgImageTestSuite) TestLosslessCompress_ShouldPreserveDimensions_WhenCompressionIsApplied() {
	compressedImg, err := suite.img.LosslessCompress()
	suite.NoError(err)
	suite.NotEmpty(compressedImg)

	metadata, err := NewProcessor(compressedImg).Metadata()
	suite.NoError(err)
	suite.Equal(suite.originalWidth, metadata.Size.Width)
	suite.Equal(suite.originalHeight, metadata.Size.Height)
}

func (suite *BimgImageTestSuite) TestMetadata_ShouldReturnCorrectMetadata_WhenImageIsValid() {
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
