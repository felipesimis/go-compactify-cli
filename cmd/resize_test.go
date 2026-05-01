package cmd

import (
	"os"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type ResizeTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *ResizeTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewResizeCmd)
}

func (suite *ResizeTestSuite) TestResizeShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--width", "150", "--height", "150")
}

func (suite *ResizeTestSuite) TestResize_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-resized"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.resizeCalled)
}

func (suite *ResizeTestSuite) TestResize_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.resizeCalled)
}

func (suite *ResizeTestSuite) TestResize_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-resized"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
	suite.True(suite.config.MockProcessor.resizeCalled)
}

func (suite *ResizeTestSuite) TestResize_ShouldReturnError_When_DimensionsAreInvalid() {
	tests := []struct {
		name   string
		width  string
		height string
	}{
		{"width_below_min", "9", "150"},
		{"height_below_min", "150", "9"},
		{"both_below_min", "9", "9"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.config.MockProcessor.resizeCalled = false
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", tt.width, "--height", tt.height})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, validation.ErrInvalidDimensions)
			suite.False(suite.config.MockProcessor.resizeCalled)
		})
	}
}

func TestResizeSuite(t *testing.T) {
	suite.Run(t, new(ResizeTestSuite))
}
