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
	suite.cmd.SetArgs([]string{"--input", "./invalid_path_name_123"})
	err := suite.cmd.Execute()

	suite.Error(err)
	suite.Contains(err.Error(), "failed to open directory")
}

func (suite *PaletteTestSuite) TestPaletteShould_Warn_When_DirectoryIsEmpty() {
	tmpDir := suite.T().TempDir()

	suite.cmd.SetArgs([]string{"--input", tmpDir})
	err := suite.cmd.Execute()

	suite.NoError(err)
	suite.Contains(suite.config.OutBuf.String(), "No files found in directory")
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
