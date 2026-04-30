package cmd

import (
	"bytes"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type mockImageProcessor struct {
	image.ImageProcessor
	paletteCalled  bool
	losslessCalled bool
}

func (m *mockImageProcessor) EnablePalette() ([]byte, error) {
	m.paletteCalled = true
	return []byte("fake-processed-bytes"), nil
}

func (m *mockImageProcessor) LosslessCompress() ([]byte, error) {
	m.losslessCalled = true
	return []byte("fake-processed-bytes"), nil
}

type TestConfig struct {
	FS            filesystem.FileSystem
	MockProcessor *mockImageProcessor
	OutBuf        *bytes.Buffer
}

func SetupTestConfig(createCmd func(filesystem.FileSystem, image.ProcessorFactory) *cobra.Command) (*cobra.Command, *TestConfig) {
	fs := filesystem.NewFileSystem()
	mockProcessor := &mockImageProcessor{}
	outBuf := new(bytes.Buffer)

	factory := func([]byte) image.ImageProcessor {
		return mockProcessor
	}

	cmd := createCmd(fs, factory)
	cmd.Flags().StringP("input", "i", "", "Input directory")
	cmd.SetOut(outBuf)
	cmd.SetErr(outBuf)

	return cmd, &TestConfig{
		FS:            fs,
		MockProcessor: mockProcessor,
		OutBuf:        outBuf,
	}
}

func AssertCommonCommandBehaviors(suite *suite.Suite, cmd *cobra.Command, config *TestConfig) {
	// Invalid directory
	cmd.SetArgs([]string{"--input", "./invalid_path_name_123"})
	err := cmd.Execute()
	suite.Error(err)
	suite.Contains(err.Error(), "failed to open directory")

	// Empty directory
	tmpDir := suite.T().TempDir()
	config.OutBuf.Reset()
	cmd.SetArgs([]string{"--input", tmpDir})
	err = cmd.Execute()
	suite.NoError(err)
	suite.Contains(config.OutBuf.String(), "No files found in directory")
}
