package cmd

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	suite.Suite
	oldWd      string
	configName string
}

func (suite *InitTestSuite) SetupTest() {
	suite.oldWd, _ = os.Getwd()
	suite.configName = "config.yaml"
	tmpDir := suite.T().TempDir()
	suite.Require().NoError(os.Chdir(tmpDir))

	viper.Reset()

	initCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})
	rootCmd.SetArgs([]string{})
}

func (suite *InitTestSuite) TearDownTest() {
	suite.Require().NoError(os.Chdir(suite.oldWd))
}

func (s *InitTestSuite) assertConfigContent(expectedSubstring string) {
	content, err := os.ReadFile(s.configName)
	s.Require().NoError(err, "should be able to read the config file")
	s.Contains(string(content), expectedSubstring)
}

func (suite *InitTestSuite) TestInitShould_CreateDefaultConfig_When_FileDoesNotExist() {
	rootCmd.SetArgs([]string{"init"})
	err := rootCmd.Execute()

	suite.NoError(err)
	suite.FileExists(suite.configName, "file config.yaml should be created")
	suite.assertConfigContent("concurrency:")
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_FileAlreadyExists() {
	importantContent := "user-custom-config: true"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(importantContent), 0644))

	rootCmd.SetArgs([]string{"init"})
	err := rootCmd.Execute()

	suite.Require().Error(err)
	suite.Contains(err.Error(), "already exists")
	suite.Contains(err.Error(), "Use --force to overwrite")
	suite.assertConfigContent(importantContent)
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_WriteFileFails() {
	suite.Require().NoError(os.Mkdir(suite.configName, 0755))

	rootCmd.SetArgs([]string{"init", "--force"})
	err := rootCmd.Execute()

	suite.Error(err)
	suite.Contains(err.Error(), "failed to create config file")
}

func (suite *InitTestSuite) TestInitShould_OverwriteConfig_When_FileExistsAndForceFlagIsUsed() {
	oldContent := "profile: old"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(oldContent), 0644))

	rootCmd.SetArgs([]string{"init", "--force"})
	err := rootCmd.Execute()
	suite.NoError(err)

	newContent, _ := os.ReadFile(suite.configName)
	suite.NotContains(string(newContent), oldContent)
	suite.Contains(string(newContent), "concurrency:")
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_ArgumentsAreProvided() {
	rootCmd.SetArgs([]string{"init", "unexpected-arg"})
	err := rootCmd.Execute()
	suite.Error(err)
}

func (suite *InitTestSuite) TestInitShould_WorkWithAliases() {
	rootCmd.SetArgs([]string{"config"})
	err := rootCmd.Execute()
	suite.NoError(err)
	suite.FileExists(suite.configName, "file config.yaml should be created")
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}
