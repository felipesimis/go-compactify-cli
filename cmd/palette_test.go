package cmd

import (
	"os"
	"path/filepath"
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

func (suite *PaletteTestSuite) TestPaletteShould_ProcessImageSuccessfully() {
	tmpDir := suite.T().TempDir()
	testPath := filepath.Join(tmpDir, "test.jpg")
	suite.Require().NoError(os.WriteFile(testPath, []byte("fake-image"), 0644))

	suite.cmd.SetArgs([]string{"--input", tmpDir})
	err := suite.cmd.Execute()

	suite.NoError(err)
	suite.True(suite.config.MockProcessor.paletteCalled, "enablePalette should have been called")

	output := suite.config.OutBuf.String()
	suite.Contains(output, "1 images")
	suite.Contains(output, "Processed")
}

func TestPaletteSuite(t *testing.T) {
	suite.Run(t, new(PaletteTestSuite))
}
