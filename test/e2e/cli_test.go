//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	_ "golang.org/x/image/webp"
)

type E2ETestSuite struct {
	suite.Suite
	binaryPath string
}

func (suite *E2ETestSuite) cleanEnv() []string {
	keep := []string{"PATH", "SystemRoot", "TMP", "TEMP", "HOME"}
	var clean []string
	for _, k := range keep {
		if val, ok := os.LookupEnv(k); ok {
			clean = append(clean, k+"="+val)
		}
	}
	return clean
}

func (suite *E2ETestSuite) buildCommand(timeout time.Duration, dir string, env []string, args ...string) (*exec.Cmd, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, suite.binaryPath, args...)

	if dir != "" {
		cmd.Dir = dir
	}

	if env != nil {
		cmd.Env = env
	} else {
		cmd.Env = suite.cleanEnv()
	}

	return cmd, cancel
}

func (suite *E2ETestSuite) SetupSuite() {
	tmpDir := suite.T().TempDir()
	binaryName := "compactify"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	suite.binaryPath = filepath.Join(tmpDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", suite.binaryPath, "../../main.go")
	output, err := buildCmd.CombinedOutput()
	suite.Require().NoError(err, "Failed to build binary for E2E tests: %s", string(output))
}

func (suite *E2ETestSuite) TestCLI_InitCommand_ShouldGenerateValidConfig() {
	tmpDir := suite.T().TempDir()

	cmd, cancel := suite.buildCommand(5*time.Second, tmpDir, nil, "init")
	defer cancel()
	suite.NoError(cmd.Run())

	configPath := filepath.Join(tmpDir, "config.yaml")
	suite.FileExists(configPath)

	content, err := os.ReadFile(configPath)
	suite.Require().NoError(err)
	contentStr := string(content)

	suite.Contains(contentStr, "# quality: 75")
	suite.Contains(contentStr, "# strip-metadata: false")
	suite.Contains(contentStr, "concurrency:")
}

func (suite *E2ETestSuite) TestCLI_Init_ShouldProtectExistingConfig() {
	tmpDir := suite.T().TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	originalContent := "custom: data"
	err := os.WriteFile(configPath, []byte(originalContent), 0644)
	suite.Require().NoError(err)

	cmd, cancel := suite.buildCommand(5*time.Second, tmpDir, nil, "init")
	defer cancel()
	output, err := cmd.CombinedOutput()

	suite.Error(err)
	suite.Contains(string(output), "already exists")

	cmdForce, cancelForce := suite.buildCommand(5*time.Second, tmpDir, nil, "init", "--force")
	defer cancelForce()
	suite.NoError(cmdForce.Run())

	newContent, err := os.ReadFile(configPath)
	suite.Require().NoError(err)
	suite.NotContains(string(newContent), originalContent)
}

func (suite *E2ETestSuite) TestCLI_Init_ShouldInjectRealCPUCount() {
	tmpDir := suite.T().TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cmd, cancel := suite.buildCommand(5*time.Second, tmpDir, nil, "init")
	defer cancel()
	suite.NoError(cmd.Run())

	content, err := os.ReadFile(configPath)
	suite.Require().NoError(err)

	expectedLine := fmt.Sprintf("concurrency: %d", runtime.NumCPU())
	suite.Contains(string(content), expectedLine)
}

func (suite *E2ETestSuite) TestCLI_ConvertCommand_ShouldProcessImages() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	outputDir := filepath.Join(suite.T().TempDir(), "e2e-output")

	cmd, cancel := suite.buildCommand(15*time.Second, "", nil, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp")
	defer cancel()

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	suite.NoError(err, "CLI command failed: %s", outputStr)
	suite.Contains(outputStr, "Converting images")

	images, err := filepath.Glob(filepath.Join(outputDir, "*.webp"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(images)

	suite.assertIsValidImage(images[0], "image/webp")
}

func (suite *E2ETestSuite) TestCLI_ShouldFail_WhenInputIsMissing() {
	cmd, cancel := suite.buildCommand(5*time.Second, "", nil, "convert", "--format", "webp")
	defer cancel()

	output, err := cmd.CombinedOutput()

	suite.Require().Error(err)
	suite.Contains(string(output), "required flag \"input\"")

	exitError, ok := err.(*exec.ExitError)
	suite.Require().True(ok)
	suite.Equal(1, exitError.ExitCode())
}

func (suite *E2ETestSuite) TestCLI_ResizeCommand_ShouldExecuteCGOWithoutSegmentationFault() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)

	inputFiles, err := filepath.Glob(filepath.Join(testDataDir, "*.*"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(inputFiles)

	originalFile, err := os.Open(inputFiles[0])
	suite.Require().NoError(err)
	defer originalFile.Close()

	originalImg, _, err := image.Decode(originalFile)
	suite.Require().NoError(err)

	targetWidth, targetHeight := 100, 100
	origBounds := originalImg.Bounds()
	suite.Require().False(origBounds.Dx() == targetWidth && origBounds.Dy() == targetHeight)

	outputDir := filepath.Join(suite.T().TempDir(), "e2e-resize-output")

	cmd, cancel := suite.buildCommand(15*time.Second, "", nil, "resize", "--input", testDataDir, "--output", outputDir, "--width", "100", "--height", "100")
	defer cancel()

	output, err := cmd.CombinedOutput()
	suite.NoError(err, "CGO interaction failed. Output:\n%s", string(output))

	images, err := filepath.Glob(filepath.Join(outputDir, "*.jpeg"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(images)

	resizedImg := suite.assertIsValidImage(images[0], "image/jpeg")
	suite.Equal(targetWidth, resizedImg.Bounds().Dx())
	suite.Equal(targetHeight, resizedImg.Bounds().Dy())
}

func (suite *E2ETestSuite) TestCLI_DryRun_ShouldNotProduceAnySideEffectsOnDisk() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	outputDir := filepath.Join(suite.T().TempDir(), "e2e-dryrun")

	cmd, cancel := suite.buildCommand(5*time.Second, "", nil, "convert", "--input", testDataDir, "--output", outputDir, "--format", "png", "--dry-run")
	defer cancel()
	suite.NoError(cmd.Run())

	images, err := filepath.Glob(filepath.Join(outputDir, "*.png"))
	suite.Require().NoError(err)
	suite.Empty(images)
}

func (suite *E2ETestSuite) TestCLI_ShouldPrioritizeFlag_OverEnvironmentVariable() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	outputDir := suite.T().TempDir()

	env := append(suite.cleanEnv(), "COMPACTIFY_CONCURRENCY=1000")
	cmd, cancel := suite.buildCommand(10*time.Second, "", env, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp", "--concurrency", "2")
	defer cancel()

	output, err := cmd.CombinedOutput()

	suite.NoError(err)
	suite.NotContains(string(output), "WARNING: Concurrency set very high")
}

func (suite *E2ETestSuite) TestCLI_ShouldRespectEnvironmentVariablePrecedence() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	outputDir := suite.T().TempDir()

	env := append(suite.cleanEnv(), "COMPACTIFY_CONCURRENCY=1000")
	cmd, cancel := suite.buildCommand(10*time.Second, "", env, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp")
	defer cancel()

	output, err := cmd.CombinedOutput()

	suite.NoError(err)
	suite.Contains(string(output), "WARNING: Concurrency set very high")
}

func (suite *E2ETestSuite) TestCLI_ShouldHandleCorruptedImagesGracefully() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)

	inputDir := suite.T().TempDir()
	outputDir := suite.T().TempDir()

	validImages, err := filepath.Glob(filepath.Join(testDataDir, "*.*"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(validImages)

	validData, err := os.ReadFile(validImages[0])
	suite.Require().NoError(err)
	suite.Require().NoError(os.WriteFile(filepath.Join(inputDir, "valid.jpg"), validData, 0644))

	fakeImage := filepath.Join(inputDir, "corrupted.jpg")
	suite.Require().NoError(os.WriteFile(fakeImage, []byte("not an image"), 0644))

	emptyImage := filepath.Join(inputDir, "empty.png")
	suite.Require().NoError(os.WriteFile(emptyImage, []byte(""), 0644))

	cmd, cancel := suite.buildCommand(10*time.Second, "", nil, "convert", "--input", inputDir, "--output", outputDir, "--format", "webp")
	defer cancel()

	output, err := cmd.CombinedOutput()

	exitError, isExitError := err.(*exec.ExitError)
	if isExitError {
		suite.NotEqual(139, exitError.ExitCode())
		suite.NotEqual(2, exitError.ExitCode())
	}

	outputStr := string(output)
	suite.Contains(outputStr, "Image buffer is empty")
	suite.Contains(outputStr, "Unsupported image format")
}

func (suite *E2ETestSuite) TestCLI_ShouldFail_When_QualityIsInvalid() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)

	cmd, cancel := suite.buildCommand(5*time.Second, "", nil, "convert", "--input", testDataDir, "--format", "webp", "--quality", "150")
	defer cancel()

	output, err := cmd.CombinedOutput()

	suite.Error(err)
	suite.Contains(string(output), "quality must be between 1 and 100")
}

func (suite *E2ETestSuite) assertIsValidImage(filePath string, expectedMimeType string) image.Image {
	f, err := os.Open(filePath)
	suite.Require().NoError(err)
	defer f.Close()

	img, format, err := image.Decode(f)
	suite.Require().NoError(err)
	suite.Require().NotNil(img)

	expectedFormat := strings.TrimPrefix(expectedMimeType, "image/")
	suite.Equal(expectedFormat, format)

	return img
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
