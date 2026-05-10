package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type FakeImageProcessor struct {
	ResultBytes []byte
	ProcessErr  error
}

func (f *FakeImageProcessor) Size() (image.ImageSize, error) { return image.ImageSize{}, nil }
func (f *FakeImageProcessor) ImageType() string              { return "" }
func (f *FakeImageProcessor) Length() int                    { return 0 }
func (f *FakeImageProcessor) Metadata() (image.ImageMetadata, error) {
	return image.ImageMetadata{}, nil
}
func (f *FakeImageProcessor) Process(opts ...image.ProcessOption) ([]byte, error) {
	return f.ResultBytes, f.ProcessErr
}

type TestConfig struct {
	FS               filesystem.FileSystem
	ProcessorFactory image.ProcessorFactory
	OutBuf           *bytes.Buffer
}

func SetupTestConfig(createCmd func(filesystem.FileSystem, image.ProcessorFactory) *cobra.Command) (*cobra.Command, *TestConfig) {
	fs := filesystem.NewFileSystem()
	outBuf := new(bytes.Buffer)

	testCfg := &TestConfig{
		FS:     fs,
		OutBuf: outBuf,
		ProcessorFactory: func([]byte) image.ImageProcessor {
			return &FakeImageProcessor{ResultBytes: []byte("fake-bytes")}
		},
	}

	factoryProxy := func(data []byte) image.ImageProcessor {
		return testCfg.ProcessorFactory(data)
	}

	cmd := createCmd(fs, factoryProxy)

	rootCmd := NewRootCmd()
	if rootCmd != nil {
		cmd.Flags().AddFlagSet(rootCmd.PersistentFlags())
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
	viper.BindPFlags(cmd.Flags())

	cmd.SetOut(outBuf)
	cmd.SetErr(outBuf)

	return cmd, testCfg
}

func AssertCommonCommandBehaviors(suite *suite.Suite, cmd *cobra.Command, config *TestConfig, extraArgs ...string) {
	sandboxDir := suite.T().TempDir()
	invalidInput := filepath.Join(sandboxDir, "invalid_path_name")

	invalidArgs := append([]string{"--input", invalidInput}, extraArgs...)
	cmd.SetArgs(invalidArgs)
	err := cmd.Execute()
	suite.Error(err, "should return an error for invalid input directory")

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

func AssertEncodingFlagsBehaviors(suite *suite.Suite, cmdFactory func(filesystem.FileSystem, image.ProcessorFactory) *cobra.Command, validBaseArgs ...string) {
	cmd, _ := SetupTestConfig(cmdFactory)
	flag := cmd.Flags().Lookup("quality")
	suite.NotNil(flag, "Expected 'quality' flag to be injected by addEncodingFlags")
	if flag != nil {
		suite.Equal("75", flag.DefValue, "Expected default quality to be 75")
	}

	cmdParsing, _ := SetupTestConfig(cmdFactory)
	argsParsing := append(validBaseArgs, "--quality", "abc")
	cmdParsing.SetArgs(argsParsing)
	errParsing := cmdParsing.Execute()
	suite.Error(errParsing)
	suite.Contains(errParsing.Error(), "invalid argument \"abc\" for \"-q, --quality\"")

	cmdValidation, _ := SetupTestConfig(cmdFactory)
	argsValidation := append(validBaseArgs, "--quality", "150")
	cmdValidation.SetArgs(argsValidation)
	errValidation := cmdValidation.Execute()
	suite.Error(errValidation)
	suite.ErrorIs(errValidation, validation.ErrInvalidQuality)
}
