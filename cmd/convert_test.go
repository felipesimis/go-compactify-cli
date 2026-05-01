package cmd

import (
	"os"
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

func (suite *ConvertTestSuite) TestConvertShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--format", "png")
}

func (suite *ConvertTestSuite) TestConvert_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--format", "png"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.png")
	suite.True(suite.config.MockProcessor.convertCalled)
}

func (suite *ConvertTestSuite) TestConvert_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.webp", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-converted-png"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--format", "png"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.png", "img2.png", "img3.png")
	suite.True(suite.config.MockProcessor.convertCalled)
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
			suite.config.MockProcessor.convertCalled = false

			inputDir := PrepareTestImages(suite.T(), "test.jpg")
			expectedOutputDir := inputDir + tt.expectedSuffix
			defer os.RemoveAll(expectedOutputDir)

			suite.cmd.SetArgs([]string{"--input", inputDir, "--format", tt.format})
			suite.NoError(suite.cmd.Execute())

			AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test"+tt.expectedExt)
			suite.True(suite.config.MockProcessor.convertCalled)
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
			suite.config.MockProcessor.convertCalled = false
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--format", tt.format})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, tt.expectedErr)
			suite.False(suite.config.MockProcessor.convertCalled)
		})
	}
}

func TestConvertSuite(t *testing.T) {
	suite.Run(t, new(ConvertTestSuite))
}
