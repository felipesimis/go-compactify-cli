package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
)

type mockImageProcessor struct {
	image.ImageProcessor
	paletteCalled   bool
	losslessCalled  bool
	grayscaleCalled bool
	flipCalled      bool
}

func (m *mockImageProcessor) EnablePalette() ([]byte, error) {
	m.paletteCalled = true
	return []byte("fake-processed-bytes"), nil
}

func (m *mockImageProcessor) LosslessCompress() ([]byte, error) {
	m.losslessCalled = true
	return []byte("fake-processed-bytes"), nil
}

func (m *mockImageProcessor) Grayscale() ([]byte, error) {
	m.grayscaleCalled = true
	return []byte("fake-processed-bytes"), nil
}

func (m *mockImageProcessor) Flip() ([]byte, error) {
	m.flipCalled = true
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

	if rootCmd != nil {
		cmd.Flags().AddFlagSet(rootCmd.PersistentFlags())
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
	cmd.SetOut(outBuf)
	cmd.SetErr(outBuf)

	return cmd, &TestConfig{
		FS:            fs,
		MockProcessor: mockProcessor,
		OutBuf:        outBuf,
	}
}

func AssertCommonCommandBehaviors(suite *suite.Suite, cmd *cobra.Command, config *TestConfig) {
	cmd.SetArgs([]string{"--input", "./invalid_path_name_123"})
	err := cmd.Execute()
	suite.Error(err, "should return an error for invalid input directory")
	suite.Contains(err.Error(), "failed to open directory")

	tmpDir := suite.T().TempDir()
	config.OutBuf.Reset()
	cmd.SetArgs([]string{"--input", tmpDir})
	err = cmd.Execute()
	suite.NoError(err)
	suite.Contains(config.OutBuf.String(), "No files found in directory")
}

func PrepareTestImages(t *testing.T, filenames ...string) string {
	tmpDir := t.TempDir()
	for _, name := range filenames {
		path := filepath.Join(tmpDir, name)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte("fake-image"), 0644)
	}
	return tmpDir
}

func AssertImageProcessed(suite *suite.Suite, config *TestConfig, expectedOutputDir string, filenames ...string) {
	for _, filename := range filenames {
		expectedOutputFile := filepath.Join(expectedOutputDir, filename)
		_, err := os.Stat(expectedOutputFile)
		suite.NoError(err)
	}

	output := config.OutBuf.String()
	suite.Contains(output, expectedOutputDir)

	expectedCountStr := fmt.Sprintf("%d images", len(filenames))
	suite.Contains(output, expectedCountStr)
}
