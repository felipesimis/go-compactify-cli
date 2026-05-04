package cmd

import (
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

func (suite *PaletteTestSuite) TearDownTest() {
	suite.config.MockProcessor.AssertExpectations(suite.T())
}

func (suite *PaletteTestSuite) TestPalette_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-palette"

	suite.config.MockProcessor.On("EnablePalette").Return([]byte("fake-bytes"), nil).Once()
	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
}

func (suite *PaletteTestSuite) TestPalette_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.config.MockProcessor.On("EnablePalette").Return([]byte("fake-bytes"), nil).Once()
	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
}

func (suite *PaletteTestSuite) TestPalette_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-palette"

	suite.config.MockProcessor.On("EnablePalette").Return([]byte("fake-bytes"), nil).Times(3)
	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
}

func TestPaletteSuite(t *testing.T) {
	suite.Run(t, new(PaletteTestSuite))
}
