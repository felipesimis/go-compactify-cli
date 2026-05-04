package cmd

import (
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type ConvertTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *ConvertTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewConvertCmd)
}

func (suite *ConvertTestSuite) TearDownTest() {
	suite.config.MockProcessor.AssertExpectations(suite.T())
}

func (suite *ConvertTestSuite) TestConvertShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--format", "png")
}

func (suite *ConvertTestSuite) TestConvert_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-converted-png"

	suite.config.MockProcessor.On("Convert", "png").Return([]byte("fake-bytes"), nil).Once()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--format", "png"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.png")
}

func (suite *ConvertTestSuite) TestConvert_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.config.MockProcessor.On("Convert", "png").Return([]byte("fake-bytes"), nil).Once()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--format", "png"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.png")
}

func (suite *ConvertTestSuite) TestConvert_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.webp", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-converted-png"

	suite.config.MockProcessor.On("Convert", "png").Return([]byte("fake-bytes"), nil).Times(3)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--format", "png"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.png", "img2.png", "img3.png")
}

func (suite *ConvertTestSuite) TestConvert_ShouldWorkWithSupportedFormats() {
	tests := []struct {
		format         string
		expectedSuffix string
		expectedExt    string
	}{
		{"png", "-converted-png", ".png"},
		{"webp", "-converted-webp", ".webp"},
		{"jpg", "-converted-jpg", ".jpg"},
		{"jpeg", "-converted-jpeg", ".jpeg"},
	}

	for _, tt := range tests {
		suite.Run(tt.format, func() {
			suite.SetupTest()

			inputDir := PrepareTestImages(suite.T(), "test.jpg")
			expectedOutputDir := inputDir + tt.expectedSuffix

			suite.config.MockProcessor.On("Convert", tt.format).Return([]byte("fake-bytes"), nil).Once()

			suite.cmd.SetArgs([]string{"--input", inputDir, "--format", tt.format})
			suite.NoError(suite.cmd.Execute())

			AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test"+tt.expectedExt)
		})
	}
}

func (suite *ConvertTestSuite) TestConvert_ShouldReturnError_When_InvalidFormat() {
	tests := []struct {
		name        string
		format      string
		expectedErr error
	}{
		{name: "empty_format", format: "", expectedErr: validation.ErrFormatRequired},
		{name: "invalid_format", format: "unsupported", expectedErr: validation.ErrInvalidFormat},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--format", tt.format})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, tt.expectedErr)
			suite.config.MockProcessor.AssertNotCalled(suite.T(), "Convert")
		})
	}
}

func (suite *ConvertTestSuite) TestModifyOutputPath_ShouldReturnFormattedPath_WhenConditionsMet() {
	tests := []struct {
		name         string
		format       string
		originalPath string
		outputDir    string
		expected     string
	}{
		{
			name:         "ShouldReplaceExtension_WhenFormatIsProvided",
			format:       "png",
			originalPath: "/path/to/image.jpg",
			outputDir:    "/output",
			expected:     filepath.Join("/output", "image.png"),
		},
		{
			name:         "ShouldReturnEmpty_WhenFormatIsEmpty",
			format:       "",
			originalPath: "/input/image.jpg",
			outputDir:    "/output",
			expected:     "",
		},
		{
			name:         "ShouldHandleMultipleDots_WhenFileNameIsComplex",
			format:       "webp",
			originalPath: "my.awesome.image.jpg",
			outputDir:    "/out",
			expected:     filepath.Join("/out", "my.awesome.image.webp"),
		},
		{
			name:         "ShouldAddExtension_WhenOriginalFileHasNone",
			format:       "jpg",
			originalPath: "image",
			outputDir:    "/out",
			expected:     filepath.Join("/out", "image.jpg"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			c := ConvertParams{Format: tt.format}
			result := c.ModifyOutputPath(tt.originalPath, tt.outputDir)
			suite.Equal(tt.expected, result)
		})
	}
}

func TestConvertSuite(t *testing.T) {
	suite.Run(t, new(ConvertTestSuite))
}
