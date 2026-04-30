package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type mockImageProcessor struct {
	image.ImageProcessor
	called bool
}

func (m *mockImageProcessor) EnablePalette() ([]byte, error) {
	m.called = true
	return []byte("fake-processed-bytes"), nil
}

type PaletteTestSuite struct {
	suite.Suite
	fs            filesystem.FileSystem
	cmd           *cobra.Command
	outBuf        *bytes.Buffer
	mockProcessor *mockImageProcessor
}

func (suite *PaletteTestSuite) SetupTest() {
	suite.fs = filesystem.NewFileSystem()
	suite.mockProcessor = &mockImageProcessor{}
	mockFactory := func([]byte) image.ImageProcessor {
		return suite.mockProcessor
	}

	suite.cmd = NewPaletteCmd(suite.fs, mockFactory)
	suite.cmd.Flags().StringP("input", "i", "", "Input directory")

	suite.outBuf = new(bytes.Buffer)
	suite.cmd.SetOut(suite.outBuf)
	suite.cmd.SetErr(suite.outBuf)
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
	suite.Contains(suite.outBuf.String(), "No files found in directory")
}

func (suite *PaletteTestSuite) TestPaletteShould_ProcessImageSuccessfully() {
	tmpDir := suite.T().TempDir()
	testPath := filepath.Join(tmpDir, "test.jpg")
	suite.Require().NoError(os.WriteFile(testPath, []byte("fake-image"), 0644))

	suite.cmd.SetArgs([]string{"--input", tmpDir})
	err := suite.cmd.Execute()

	suite.NoError(err)
	suite.True(suite.mockProcessor.called, "enablePalette should have been called")

	output := suite.outBuf.String()
	suite.Contains(output, "1 images")
	suite.Contains(output, "Processed")
}

func TestPaletteSuite(t *testing.T) {
	suite.Run(t, new(PaletteTestSuite))
}
