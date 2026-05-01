package cmd

import (
	"os"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CropTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *CropTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewCropCmd)
}

func (suite *CropTestSuite) TestCropShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--width", "150", "--height", "150", "--gravity", "0")
}

func (suite *CropTestSuite) TestCrop_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-cropped_150x150"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150", "--gravity", "0"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.cropCalled)
}

func (suite *CropTestSuite) TestCrop_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--width", "150", "--height", "150", "--gravity", "0"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.cropCalled)
}

func (suite *CropTestSuite) TestCrop_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-cropped_150x150"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150", "--gravity", "0"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
	suite.True(suite.config.MockProcessor.cropCalled)
}

func (suite *CropTestSuite) TestCrop_ShouldReturnError_When_DimensionsAreInvalid() {
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
			suite.config.MockProcessor.cropCalled = false
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", tt.width, "--height", tt.height, "--gravity", "0"})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, validation.ErrInvalidDimensions)
			suite.False(suite.config.MockProcessor.cropCalled)
		})
	}
}

func (suite *CropTestSuite) TestCrop_ShouldReturnError_When_GravityIsInvalid() {
	tests := []struct {
		name    string
		gravity string
	}{
		{"gravity_below_min", "-1"},
		{"gravity_above_max", "6"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.config.MockProcessor.cropCalled = false
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", "150", "--height", "150", "--gravity", tt.gravity})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, validation.ErrInvalidGravity)
			suite.False(suite.config.MockProcessor.cropCalled)
		})
	}
}

func (suite *CropTestSuite) TestCrop_ShouldReturnMultipleErrors_When_AllFlagsAreInvalid() {
	suite.config.MockProcessor.cropCalled = false
	suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", "9", "--height", "9", "--gravity", "6"})

	err := suite.cmd.Execute()
	suite.Error(err)
	suite.ErrorIs(err, validation.ErrInvalidDimensions)
	suite.ErrorIs(err, validation.ErrInvalidGravity)
	suite.False(suite.config.MockProcessor.cropCalled)
}

func (suite *CropTestSuite) TestCrop_ShouldWorkWithAllValidFlags() {
	gravities := []string{"0", "1", "2", "3", "4", "5"}
	for _, gravity := range gravities {
		suite.Run("gravity_"+gravity, func() {
			suite.config.MockProcessor.cropCalled = false
			inputDir := PrepareTestImages(suite.T(), "test.jpg")
			expectedOutputDir := inputDir + "-cropped_150x150"
			defer os.RemoveAll(expectedOutputDir)

			suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150", "--gravity", gravity})
			suite.NoError(suite.cmd.Execute())

			AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
			suite.True(suite.config.MockProcessor.cropCalled)
		})
	}
}

func TestCropSuite(t *testing.T) {
	suite.Run(t, new(CropTestSuite))
}
