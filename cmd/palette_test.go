package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type PaletteTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *PaletteTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewPaletteCmd)
}

func (suite *PaletteTestSuite) TestPaletteShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config)
}

func (suite *PaletteTestSuite) TestPalette_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-palette"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.paletteCalled)
}

func (suite *PaletteTestSuite) TestPalette_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
}

func (suite *PaletteTestSuite) TestPalette_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-palette"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
	suite.True(suite.config.MockProcessor.paletteCalled)
}

func TestPaletteSuite(t *testing.T) {
	suite.Run(t, new(PaletteTestSuite))
}
