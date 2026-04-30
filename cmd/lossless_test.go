package cmd

import (
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
	suite.cmd.SetArgs([]string{"--input", "./invalid_path_name_123"})
	err := suite.cmd.Execute()

	suite.Error(err)
	suite.Contains(err.Error(), "failed to open directory")
}

func TestLosslessSuite(t *testing.T) {
	suite.Run(t, new(LosslessTestSuite))
}
