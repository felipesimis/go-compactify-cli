package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	oldWd      string
	configName string
	rootCmd    *cobra.Command
}

func (suite *RootTestSuite) SetupTest() {
	suite.oldWd, _ = os.Getwd()
	tmpDir := suite.T().TempDir()
	suite.configName = "config.yaml"
	suite.Require().NoError(os.Chdir(tmpDir))

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "COMPACTIFY_") {
			parts := strings.SplitN(env, "=", 2)[0]
			os.Unsetenv(parts)
		}
	}

	viper.Reset()
	cfgFile = ""

	suite.rootCmd = NewRootCmd()

	suite.rootCmd.AddCommand(&cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	})

	suite.rootCmd.SetOut(io.Discard)
	suite.rootCmd.SetErr(io.Discard)
}

func (suite *RootTestSuite) TearDownTest() {
	suite.Require().NoError(os.Chdir(suite.oldWd))
}

func (suite *RootTestSuite) TestShould_UseConcurrencyFromFile_When_NoFlagIsProvided() {
	configContent := "concurrency: 4\ninput: ./fake-dir\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))

	suite.rootCmd.SetArgs([]string{"dummy"})
	err := suite.rootCmd.Execute()
	suite.NoError(err)

	cfg := loadAppConfig()
	suite.Equal(4, cfg.Concurrency)
}

func (suite *RootTestSuite) TestShould_PrioritizeEnvVar_Over_ConfigFile() {
	configContent := "concurrency: 4\ninput: ./fake-dir\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))

	os.Setenv("COMPACTIFY_CONCURRENCY", "12")
	defer os.Unsetenv("COMPACTIFY_CONCURRENCY")

	suite.rootCmd.SetArgs([]string{"dummy"})
	err := suite.rootCmd.Execute()
	suite.NoError(err)

	cfg := loadAppConfig()
	suite.Equal(12, cfg.Concurrency, "should prioritize environment variable over config file")
}

func (suite *RootTestSuite) TestShould_PrioritizeFlag_Over_EnvVar_And_ConfigFile() {
	configContent := "concurrency: 4\ninput: ./fake-dir\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))

	os.Setenv("COMPACTIFY_CONCURRENCY", "12")
	defer os.Unsetenv("COMPACTIFY_CONCURRENCY")

	suite.rootCmd.SetArgs([]string{"dummy", "--concurrency", "2"})
	err := suite.rootCmd.Execute()
	suite.NoError(err)

	cfg := loadAppConfig()
	suite.Equal(2, cfg.Concurrency, "should prioritize command-line flag over environment variable and config file")
}

func (suite *RootTestSuite) TestShould_ReturnError_When_InputFlagIsMissing() {
	suite.rootCmd.SetArgs([]string{"dummy"})
	err := suite.rootCmd.Execute()

	suite.Error(err)
	suite.Contains(err.Error(), "required flag \"input\" (-i) not set")
}

func (suite *RootTestSuite) TestShould_ShowWarning_When_ConcurrencyIsTooHigh() {
	buf := new(bytes.Buffer)
	suite.rootCmd.SetOut(buf)
	suite.rootCmd.SetErr(buf)

	suite.rootCmd.SetArgs([]string{"dummy", "--input", "./fake-dir", "--concurrency", "1000"})
	err := suite.rootCmd.Execute()

	suite.NoError(err)
	suite.Contains(buf.String(), "WARNING: Concurrency set very high")
}

func (suite *RootTestSuite) TestShould_LoadSpecificConfigFile_When_ConfigFlagIsProvided() {
	customConfigFile := "custom_config.yaml"
	suite.Require().NoError(os.WriteFile(customConfigFile, []byte("concurrency: 5\ninput: ./fake-dir\n"), 0644))

	suite.rootCmd.SetArgs([]string{"dummy", "--config", customConfigFile})
	err := suite.rootCmd.Execute()
	suite.NoError(err)

	cfg := loadAppConfig()
	suite.Equal(5, cfg.Concurrency, "should load concurrency from specified config file")
}

func (suite *RootTestSuite) TestShould_PrintError_When_ConfigFileIsCorrupted() {
	corruptedContent := "concurrency: [invalid-syntax\ninput: ./fake-dir\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(corruptedContent), 0644))

	buf := new(bytes.Buffer)
	suite.rootCmd.SetOut(buf)
	suite.rootCmd.SetErr(buf)

	suite.rootCmd.SetArgs([]string{"dummy", "--input", "./fake-dir"})
	err := suite.rootCmd.Execute()

	suite.NoError(err)
	suite.Contains(buf.String(), "Error reading config file")
	suite.Contains(buf.String(), "yaml: line 1")
}

func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func (suite *RootTestSuite) TestExecute_ShouldFormatVersionWithVPrefix() {
	oldVersion := Version
	Version = "v2.0.0"
	defer func() { Version = oldVersion }()

	oldArgs := os.Args
	os.Args = []string{"compactify", "--version"}
	defer func() { os.Args = oldArgs }()

	output, err := captureOutput(Execute)

	suite.NoError(err)
	suite.Contains(output, "v2.0.0")
	suite.NotContains(output, "vv2.0.0")
}

func (suite *RootTestSuite) TestPersistentPreRunE_ShouldBypassValidation_ForSpecificCommands() {
	initCmd := &cobra.Command{Use: "init", Run: func(cmd *cobra.Command, args []string) {}}
	suite.rootCmd.AddCommand(initCmd)
	suite.rootCmd.SetArgs([]string{"init"})
	err := suite.rootCmd.Execute()
	suite.NoError(err, "should not require --input for init command")

	suite.rootCmd.SetArgs([]string{"help"})
	err = suite.rootCmd.Execute()
	suite.NoError(err, "should not require --input for help command")

	suite.rootCmd.SetArgs([]string{"--version"})
	err = suite.rootCmd.Execute()
	suite.NoError(err, "should not require --input when --version is used")
}

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
