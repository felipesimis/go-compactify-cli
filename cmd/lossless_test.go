package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type LosslessTestSuite struct {
	suite.Suite
	cmd    *cobra.Command
	config *TestConfig
}

func (suite *LosslessTestSuite) SetupTest() {
	suite.cmd, suite.config = SetupTestConfig(NewLosslessCmd)
}

func (suite *LosslessTestSuite) TestLosslessShould_ReturnError_When_InputDirectoryDoesNotExist() {
	AssertCommonCommandBehaviors(&suite.Suite, suite.cmd, suite.config)
}

func (suite *LosslessTestSuite) TestLosslessShould_ProcessImageSuccessfully() {
	tmpDir := suite.T().TempDir()
	testPath := filepath.Join(tmpDir, "test.jpg")
	suite.Require().NoError(os.WriteFile(testPath, []byte("fake"), 0644))

	suite.cmd.SetArgs([]string{"--input", tmpDir})
	err := suite.cmd.Execute()

	suite.NoError(err)
	suite.True(suite.config.MockProcessor.losslessCalled)
	suite.Contains(suite.config.OutBuf.String(), "Processed")
}

func TestLosslessSuite(t *testing.T) {
	suite.Run(t, new(LosslessTestSuite))
}
