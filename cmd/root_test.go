package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	oldWd      string
	configName string
}

func (suite *RootTestSuite) SetupTest() {
	suite.oldWd, _ = os.Getwd()
	tmpDir := suite.T().TempDir()
	suite.configName = "config.yaml"
	suite.Require().NoError(os.Chdir(tmpDir))

	viper.Reset()

	resetFlags := func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	}
	rootCmd.Flags().VisitAll(resetFlags)
	rootCmd.PersistentFlags().VisitAll(resetFlags)
}

func (suite *RootTestSuite) TearDownTest() {
	suite.Require().NoError(os.Chdir(suite.oldWd))
}

func (suite *RootTestSuite) TestShould_UseConcurrencyFromFile_When_NoFlagIsProvided() {
	configContent := "concurrency: 4\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))

	initConfig()
	concurrency, _ := rootCmd.Flags().GetInt("concurrency")
	suite.Equal(4, concurrency)
}

func (suite *RootTestSuite) TestShould_PrioritizeEnvVar_Over_ConfigFile() {
	configContent := "concurrency: 4\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))
	os.Setenv("COMPACTIFY_CONCURRENCY", "12")
	defer os.Unsetenv("COMPACTIFY_CONCURRENCY")

	initConfig()

	concurrency, _ := rootCmd.Flags().GetInt("concurrency")
	suite.Equal(12, concurrency)
}

func (suite *RootTestSuite) TestShould_PrioritizeFlag_Over_EnvVar_And_ConfigFile() {
	configContent := "concurrency: 4\n"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(configContent), 0644))
	os.Setenv("COMPACTIFY_CONCURRENCY", "12")
	defer os.Unsetenv("COMPACTIFY_CONCURRENCY")

	rootCmd.Flags().Set("concurrency", "2")
	initConfig()

	concurrency, _ := rootCmd.Flags().GetInt("concurrency")
	suite.Equal(2, concurrency, "should prioritize command-line flag over environment variable and config file")
}

func (suite *RootTestSuite) TestShould_ReturnError_When_InputFlagIsMissing() {
	rootCmd.Flags().Set("input", "")
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	suite.Error(err)
	suite.Contains(err.Error(), "required flag \"input\" (-i) not set")
}

func (suite *RootTestSuite) TestShould_ShowWarning_When_ConcurrencyIsTooHigh() {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	defer rootCmd.SetOut(os.Stdout)

	rootCmd.Flags().Set("input", "./fake-dir")
	rootCmd.Flags().Set("concurrency", "1000")

	err := rootCmd.PersistentPreRunE(rootCmd, []string{})

	suite.NoError(err)
	suite.Contains(buf.String(), "WARNING: Concurrency set very high. This may cause high memory usage and slow down your system.")
}

func (suite *RootTestSuite) TestShould_LoadSpecificConfigFile_When_ConfigFlagIsProvided() {
	customConfigFile := "custom_config.yaml"
	suite.Require().NoError(os.WriteFile(customConfigFile, []byte("concurrency: 5\n"), 0644))

	rootCmd.Flags().Set("config", customConfigFile)
	initConfig()

	concurrency, _ := rootCmd.Flags().GetInt("concurrency")
	suite.Equal(5, concurrency)
}

func (suite *RootTestSuite) TestShould_PrintError_When_ConfigFileIsCorrupted() {
	corruptedContent := "concurrency: [invalid-systax"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(corruptedContent), 0644))

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	defer rootCmd.SetOut(os.Stdout)
	defer rootCmd.SetErr(os.Stderr)

	initConfig()

	suite.Contains(buf.String(), "Error reading config file")
	suite.Contains(buf.String(), "yaml: line 1")
}

func (suite *RootTestSuite) TestExecute_ShouldInitializeConfigAndVersion() {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	defer rootCmd.SetOut(os.Stdout)

	rootCmd.SetArgs([]string{"--version"})

	err := Execute()
	suite.NoError(err)

	output := buf.String()
	suite.Contains(output, "Compactify")
	suite.Contains(output, "v")
}

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
