package cmd

import (
	"bytes"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
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
