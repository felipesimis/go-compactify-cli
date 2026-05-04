package cmd

import (
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type LosslessTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *LosslessTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewLosslessCmd)
	suite.config.ProcessorFactory = func([]byte) image.ImageProcessor {
		return &fakeProcessor{resultBytes: []byte("fake-bytes")}
	}
}

func (suite *LosslessTestSuite) TestLosslessShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config)
}

func (suite *LosslessTestSuite) TestLossless_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-lossless"

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
}

func (suite *LosslessTestSuite) TestLossless_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
}

func (suite *LosslessTestSuite) TestLossless_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-lossless"

	suite.cmd.SetArgs([]string{"--input", inputDir})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
}

func TestLosslessSuite(t *testing.T) {
	suite.Run(t, new(LosslessTestSuite))
}
