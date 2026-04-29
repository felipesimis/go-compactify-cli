package cmd

import (
	"bytes"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type PaletteTestSuite struct {
	suite.Suite
	fs     filesystem.FileSystem
	cmd    *cobra.Command
	outBuf *bytes.Buffer
}

func (suite *PaletteTestSuite) SetupTest() {
	suite.fs = filesystem.NewFileSystem()
	suite.cmd = NewPaletteCmd(suite.fs)
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

func TestPaletteSuite(t *testing.T) {
	suite.Run(t, new(PaletteTestSuite))
}
