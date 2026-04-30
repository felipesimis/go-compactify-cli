package cmd

import (
	"os"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type ThumbnailTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *ThumbnailTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewThumbnailCmd)
}

func (suite *ThumbnailTestSuite) TestThumbnailShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--width", "150")
}

func (suite *ThumbnailTestSuite) TestThumbnail_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-thumbnail"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.thumbnailCalled)
}

func (suite *ThumbnailTestSuite) TestThumbnail_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--width", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.thumbnailCalled)
}

func (suite *ThumbnailTestSuite) TestThumbnail_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-thumbnail"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
	suite.True(suite.config.MockProcessor.thumbnailCalled)
}

func (suite *ThumbnailTestSuite) TestThumbnail_ShouldReturnError_When_WidthIsTooSmall() {
	suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", "49"})
	err := suite.cmd.Execute()

	suite.Error(err)
	suite.ErrorIs(err, validation.ErrWidthTooSmall)
	suite.False(suite.config.MockProcessor.thumbnailCalled)
}

func (suite *ThumbnailTestSuite) TestThumbnail_ShouldReturnError_When_WidthIsTooLarge() {
	suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", "1025"})
	err := suite.cmd.Execute()

	suite.Error(err)
	suite.ErrorIs(err, validation.ErrWidthTooLarge)
	suite.False(suite.config.MockProcessor.thumbnailCalled)
}

func TestThumbnailSuite(t *testing.T) {
	suite.Run(t, new(ThumbnailTestSuite))
}
