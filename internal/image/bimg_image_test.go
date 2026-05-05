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

func (suite *BimgImageTestSuite) TestProcessSizingOperations_ShouldTransformImageCorrectly_WhenGivenOptions() {
	tests := []struct {
		name           string
		options        []ProcessOption
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "Resize",
			options:        []ProcessOption{WithResize(300, 200)},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name:           "Crop",
			options:        []ProcessOption{WithCrop(300, 200, GravitySmart)},
			expectedWidth:  300,
			expectedHeight: 200,
		},
		{
			name:           "Enlarge",
			options:        []ProcessOption{WithEnlarge(1200, 800)},
			expectedWidth:  1200,
			expectedHeight: 800,
		},
		{
			name:           "Thumbnail",
			options:        []ProcessOption{WithThumbnail(300)},
			expectedWidth:  300,
			expectedHeight: 300,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			processedImg, err := suite.img.Process(tt.options...)
			suite.NoError(err)
			suite.NotEmpty(processedImg)

			size, err := NewProcessor(processedImg).Size()
			suite.NoError(err)
			suite.Equal(tt.expectedWidth, size.Width)
			suite.Equal(tt.expectedHeight, size.Height)
		})
	}
}

func (suite *BimgImageTestSuite) TestProcess_ShouldNotPanic_WhenNilOptionIsProvided() {
	processedImg, err := suite.img.Process(WithResize(100, 100), nil, WithGrayscale())

	suite.NoError(err)
	suite.NotEmpty(processedImg)

	size, _ := NewProcessor(processedImg).Size()
	suite.Equal(100, size.Width)
}

func (suite *BimgImageTestSuite) TestProcess_ShouldChangeImageType_WhenConvertOptionIsProvided() {
	convertedImg, err := suite.img.Process(WithConvert("png"))
	suite.NoError(err)
	suite.NotEmpty(convertedImg)
	suite.Equal("png", NewProcessor(convertedImg).ImageType())
}

func (suite *BimgImageTestSuite) TestProcess_ShouldReturnError_WhenUnsupportedFormatIsProvided() {
	convertedImg, err := suite.img.Process(WithConvert("invalid_format"))
	suite.ErrorIs(err, ErrUnsupportedImageType)
	suite.Empty(convertedImg)
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

func (suite *BimgImageTestSuite) TestProcess_ShouldMaintainDimensions_WhenImageIsFlipped() {
	flippedImg, err := suite.img.Process(WithFlip())
	suite.NoError(err)
	suite.NotEmpty(flippedImg)

	originalSize, err := suite.img.Size()
	suite.NoError(err)
	flippedImgSize, err := NewProcessor(flippedImg).Size()
	suite.NoError(err)

	suite.Equal(originalSize.Width, flippedImgSize.Width)
	suite.Equal(originalSize.Height, flippedImgSize.Height)
}

func (suite *BimgImageTestSuite) TestProcess_ShouldReturnProcessedImage_WhenGrayscaleOptionProvided() {
	grayscaleImg, err := suite.img.Process(WithGrayscale())
	suite.NoError(err)
	suite.NotEmpty(grayscaleImg)
}

func (suite *BimgImageTestSuite) TestProcess_ShouldChangeImageLength_WhenPaletteOptionProvided() {
	initialImgLength := suite.img.Length()

	paletteImg, err := suite.img.Process(WithConvert("png"), WithPalette())
	suite.NoError(err)
	suite.NotEmpty(paletteImg)

	paletteImgLength := NewProcessor(paletteImg).Length()
	suite.NotZero(paletteImgLength)
	suite.NotEqual(initialImgLength, paletteImgLength, "Expected image data to change after applying palette")
}

func (suite *BimgImageTestSuite) TestProcess_ShouldPreserveDimensions_WhenLosslessOptionProvided() {
	compressedImg, err := suite.img.Process(WithConvert("webp"), WithLosslessCompress())
	suite.NoError(err)
	suite.NotEmpty(compressedImg)

	metadata, err := NewProcessor(compressedImg).Metadata()
	suite.NoError(err)
	suite.Equal(suite.originalWidth, metadata.Size.Width)
	suite.Equal(suite.originalHeight, metadata.Size.Height)
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

func (suite *BimgImageTestSuite) TestProcess_ShouldCombineMultipleOperations_WhenMultipleOptionsAreProvided() {
	opts := []ProcessOption{
		WithThumbnail(200),
		WithGrayscale(),
		WithConvert("png"),
	}

	processedImg, err := suite.img.Process(opts...)
	suite.NoError(err)
	suite.NotEmpty(processedImg)

	newProc := NewProcessor(processedImg)

	size, _ := newProc.Size()
	suite.Equal(200, size.Width)
	suite.Equal(200, size.Height)
	suite.Equal("png", newProc.ImageType())
}

func (suite *BimgImageTestSuite) TestProcess_ThumbnailShouldOverrideExplicitResize_WhenBothAreProvided() {
	opts := []ProcessOption{
		WithResize(800, 600),
		WithThumbnail(100),
	}

	processedImg, err := suite.img.Process(opts...)
	suite.NoError(err)

	size, _ := NewProcessor(processedImg).Size()
	suite.Equal(100, size.Width)
	suite.NotEqual(800, size.Width)
}

func (suite *BimgImageTestSuite) TestProcess_ShouldApplyStripMetadata_WhenOptionIsProvided() {
	processedImg, err := suite.img.Process(WithStripMetadata())
	suite.NoError(err)
	suite.NotEmpty(processedImg)
	suite.Less(len(processedImg), suite.originalLength)
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

func (suite *BimgImageTestSuite) TestMetadata_ShouldReturnCorrectMetadata_WhenImageIsValid() {
	metadata, err := suite.img.Metadata()
	suite.NoError(err)
	suite.NotEmpty(metadata)
	suite.Equal(suite.originalWidth, metadata.Size.Width)
	suite.Equal(suite.originalHeight, metadata.Size.Height)
	suite.Equal("jpeg", metadata.Type)
}

func (suite *BimgImageTestSuite) TestLength_ShouldReturnCorrectByteLength() {
	suite.Equal(suite.originalLength, suite.img.Length())
}

func TestBimgImageTestSuite(t *testing.T) {
	suite.Run(t, new(BimgImageTestSuite))
}
