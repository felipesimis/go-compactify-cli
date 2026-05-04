package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	suite.Suite
	oldWd      string
	configName string
	fs         filesystem.FileSystem
	cmd        *cobra.Command
}

func (suite *InitTestSuite) SetupTest() {
	suite.oldWd, _ = os.Getwd()
	suite.configName = "config.yaml"
	tmpDir := suite.T().TempDir()
	suite.Require().NoError(os.Chdir(tmpDir))

	suite.fs = filesystem.NewFileSystem()
	suite.cmd = NewInitCmd(suite.fs)
}

func (suite *InitTestSuite) TearDownTest() {
	suite.Require().NoError(os.Chdir(suite.oldWd))
}

func (suite *InitTestSuite) assertConfigContent(expectedSubstring string) {
	content, err := os.ReadFile(suite.configName)
	suite.Require().NoError(err, "should be able to read the config file")
	suite.Contains(string(content), expectedSubstring)
}

func (suite *InitTestSuite) TestInitShould_CreateDefaultConfig_When_FileDoesNotExist() {
	suite.cmd.SetArgs([]string{})
	err := suite.cmd.Execute()

	suite.NoError(err)
	suite.FileExists(suite.configName, "file config.yaml should be created")
	suite.assertConfigContent("concurrency:")
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_FileAlreadyExists() {
	importantContent := "user-custom-config: true"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(importantContent), 0644))

	suite.cmd.SetArgs([]string{})
	err := suite.cmd.Execute()

	suite.Require().Error(err)
	suite.Contains(err.Error(), "already exists")
	suite.Contains(err.Error(), "Use --force to overwrite")
	suite.assertConfigContent(importantContent)
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_WriteFileFails() {
	suite.Require().NoError(os.Mkdir(suite.configName, 0755))

	suite.cmd.SetArgs([]string{"--force"})
	err := suite.cmd.Execute()

	suite.Error(err)
	suite.Contains(err.Error(), "failed to create config file")
}

func (suite *InitTestSuite) TestInitShould_OverwriteConfig_When_FileExistsAndForceFlagIsUsed() {
	oldContent := "profile: old"
	suite.Require().NoError(os.WriteFile(suite.configName, []byte(oldContent), 0644))

	suite.cmd.SetArgs([]string{"--force"})
	err := suite.cmd.Execute()
	suite.NoError(err)

	newContent, _ := os.ReadFile(suite.configName)
	suite.NotContains(string(newContent), oldContent)
	suite.Contains(string(newContent), "concurrency:")
}

func (suite *InitTestSuite) TestInitShould_ReturnError_When_ArgumentsAreProvided() {
	suite.cmd.SetArgs([]string{"unexpected-arg"})
	err := suite.cmd.Execute()
	suite.Error(err)
}

func (suite *InitTestSuite) TestInitShould_WorkWithAliases() {
	dummyRoot := &cobra.Command{Use: "dummy"}
	dummyRoot.AddCommand(suite.cmd)

	dummyRoot.SetArgs([]string{"config"})
	err := dummyRoot.Execute()

	suite.NoError(err)
	suite.FileExists(suite.configName, "file config.yaml should be created")
}

func (suite *InitTestSuite) TestInitShould_PrintSuccessMessage_When_Initialized() {
	buf := new(bytes.Buffer)
	suite.cmd.SetOut(buf)

	suite.cmd.SetArgs([]string{})
	err := suite.cmd.Execute()
	suite.NoError(err)

	output := buf.String()
	suite.Contains(output, "Configuration file initialized successfully")
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}
