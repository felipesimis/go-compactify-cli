package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type FlipTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *FlipTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewFlipCmd)
}

func (suite *FlipTestSuite) TestFlipShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config)
}

func (suite *FlipTestSuite) TestFlip_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-flip"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.flipCalled)
}

func (suite *FlipTestSuite) TestFlip_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
	suite.True(suite.config.MockProcessor.flipCalled)
}

func (suite *FlipTestSuite) TestFlip_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-flip"
	defer os.RemoveAll(expectedOutputDir)

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
	suite.True(suite.config.MockProcessor.flipCalled)
}

func TestFlipSuite(t *testing.T) {
	suite.Run(t, new(FlipTestSuite))
}
