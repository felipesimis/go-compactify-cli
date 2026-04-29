package cmd

import (
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

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
