package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	binaryPath string
}

func (suite *E2ETestSuite) SetupSuite() {
	tmpDir := suite.T().TempDir()

	binaryName := "compactify"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	suite.binaryPath = filepath.Join(tmpDir, binaryName)

	buildCmd := exec.Command("go", "build", "-ldflags", "-w -s", "-o", suite.binaryPath, "--trimpath", "../../main.go")

	output, err := buildCmd.CombinedOutput()
	suite.Require().NoError(err, "Failed to build binary for E2E tests: %s", string(output))
}

func (suite *E2ETestSuite) TestCLI_ConvertCommand_ShouldProcessImagesFillsPipeline() {
	testDataDir := filepath.Join("..", "..", "internal", "image", "testdata")

	_, err := os.Stat(filepath.Join(testDataDir, "sample.jpeg"))
	suite.Require().NoError(err, "Real test image not found in testdata directory")

	outputDir := filepath.Join(suite.T().TempDir(), "e2e-output")

	cmd := exec.Command(suite.binaryPath, "convert",
		"--input", testDataDir,
		"--output", outputDir,
		"--format", "webp",
	)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	suite.NoError(err, "CLI command should execute without errors. Output:\n%s", outputStr)
	suite.Contains(outputStr, "Converting images", "Output should contain the progress message")
	suite.Contains(outputStr, "OUTPUT DIRECTORY", "Output should indicate where files are saved")

	files, err := filepath.Glob(filepath.Join(outputDir, "*.webp"))
	suite.NoError(err, "Failed to read output directory")
	suite.NotEmpty(files, "Expected at least one .webp file to be generated")

	fileInfo, err := os.Stat(files[0])
	suite.NoError(err)
	suite.Greater(fileInfo.Size(), int64(0), "Generated image file should not be empty")
}

func (suite *E2ETestSuite) TestCLI_ShouldFailAndReturnExitCode1_WhenInputIsMissing() {
	cmd := exec.Command(suite.binaryPath, "convert", "--format", "webp")

	output, err := cmd.CombinedOutput()
	suite.Require().Error(err, "Command should fail when --input is missing")

	outputStr := string(output)
	suite.Contains(outputStr, "required flag \"input\" (-i) not set")

	exitError, ok := err.(*exec.ExitError)
	suite.True(ok, "Error should be an ExitError")
	suite.Equal(1, exitError.ExitCode(), "CLI should return exit code 1 on fatal error")
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
