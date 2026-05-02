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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockImageProcessor struct {
	image.ImageProcessor
	mock.Mock
}

func (m *mockImageProcessor) EnablePalette() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) LosslessCompress() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Grayscale() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Flip() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Thumbnail(width int) ([]byte, error) {
	args := m.Called(width)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Resize(width, height int) ([]byte, error) {
	args := m.Called(width, height)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Enlarge(width, height int) ([]byte, error) {
	args := m.Called(width, height)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Convert(format string) ([]byte, error) {
	args := m.Called(format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockImageProcessor) Crop(width int, height int, gravity image.Gravity) ([]byte, error) {
	args := m.Called(width, height, gravity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
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

func AssertCommonCommandBehaviors(suite *suite.Suite, cmd *cobra.Command, config *TestConfig, extraArgs ...string) {
	invalidArgs := append([]string{"--input", "./invalid_path_name"}, extraArgs...)
	cmd.SetArgs(invalidArgs)
	err := cmd.Execute()
	suite.Error(err, "should return an error for invalid input directory")
	suite.Contains(err.Error(), "failed to open directory")

	tmpDir := suite.T().TempDir()
	config.OutBuf.Reset()

	emptyArgs := append([]string{"--input", tmpDir}, extraArgs...)
	cmd.SetArgs(emptyArgs)
	err = cmd.Execute()
	suite.NoError(err)
	suite.Contains(config.OutBuf.String(), "No files found in directory")
}

func PrepareTestImages(t *testing.T, filenames ...string) string {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	os.MkdirAll(inputDir, 0755)

	for _, name := range filenames {
		path := filepath.Join(inputDir, name)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte("fake-image"), 0644)
	}
	return inputDir
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
