package cmd

import (
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type EnlargeTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *EnlargeTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewEnlargeCmd)
}

func (suite *EnlargeTestSuite) TearDownTest() {
	suite.config.MockProcessor.AssertExpectations(suite.T())
}

func (suite *EnlargeTestSuite) TestEnlargeShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config, "--width", "150", "--height", "150")
}

func (suite *EnlargeTestSuite) TestEnlarge_ShouldWorkWithDefaultOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	expectedOutputDir := inputDir + "-enlarged-150x150"

	suite.config.MockProcessor.On("Enlarge", 150, 150).Return([]byte("fake-bytes"), nil).Once()
	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "test.jpg")
}

func (suite *EnlargeTestSuite) TestEnlarge_ShouldWorkWithCustomOutput() {
	inputDir := PrepareTestImages(suite.T(), "test.jpg")
	customOutputDir := suite.T().TempDir()

	suite.config.MockProcessor.On("Enlarge", 150, 150).Return([]byte("fake-bytes"), nil).Once()
	suite.cmd.SetArgs([]string{"--input", inputDir, "--output", customOutputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, customOutputDir, "test.jpg")
}

func (suite *EnlargeTestSuite) TestEnlarge_ShouldProcessMultipleImages() {
	inputDir := PrepareTestImages(suite.T(), "img1.jpg", "img2.jpg", "img3.jpg")
	expectedOutputDir := inputDir + "-enlarged-150x150"

	suite.config.MockProcessor.On("Enlarge", 150, 150).Return([]byte("fake-bytes"), nil).Times(3)
	suite.cmd.SetArgs([]string{"--input", inputDir, "--width", "150", "--height", "150"})
	suite.NoError(suite.cmd.Execute())

	AssertImageProcessed(&suite.Suite, suite.config, expectedOutputDir, "img1.jpg", "img2.jpg", "img3.jpg")
}

func (suite *EnlargeTestSuite) TestEnlarge_ShouldReturnError_When_DimensionsAreInvalid() {
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
			suite.SetupTest()
			suite.cmd.SetArgs([]string{"--input", "some/dir", "--width", tt.width, "--height", tt.height})

			err := suite.cmd.Execute()
			suite.Error(err)
			suite.ErrorIs(err, validation.ErrInvalidDimensions)
			suite.config.MockProcessor.AssertNotCalled(suite.T(), "Enlarge")
		})
	}
}

func TestEnlargeSuite(t *testing.T) {
	suite.Run(t, new(EnlargeTestSuite))
}
