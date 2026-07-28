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

func (suite *E2ETestSuite) TestCLI_Recursive_ShouldNotLoopInfinitely_WhenOutputIsInsideInput() {
	inputDir := suite.T().TempDir()
	outputDir := filepath.Join(inputDir, "output_bomb")
	os.WriteFile(filepath.Join(inputDir, "image.jpg"), []byte("fake image data"), 0644)

	cmd, cancel := suite.buildCommand(5*time.Second, "", nil, "convert", "--input", inputDir, "--output", outputDir, "--format", "webp", "--recursive")
	defer cancel()

	_, err := cmd.CombinedOutput()
	suite.NoError(err)

	nestedBombPath := filepath.Join(outputDir, "output_bomb")
	suite.NoDirExists(nestedBombPath, "nested output directory was created, indicating a potential infinite loop")
}

func (suite *E2ETestSuite) TestCLI_Recursive_ShouldNotCreateDirectories_WhenNoImagesArePresent() {
	inputDir := suite.T().TempDir()
	outputDir := suite.T().TempDir()

	validDir := filepath.Join(inputDir, "valid_photos")
	ghostDir := filepath.Join(inputDir, "ghost_docs")
	os.Mkdir(validDir, 0755)
	os.Mkdir(ghostDir, 0755)

	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	realImage, _ := filepath.Glob(filepath.Join(testDataDir, "*.*"))
	suite.Require().NotEmpty(realImage)

	imgContent, _ := os.ReadFile(realImage[0])
	os.WriteFile(filepath.Join(validDir, "photo.jpg"), imgContent, 0644)
	os.WriteFile(filepath.Join(ghostDir, "doc.txt"), []byte("text file"), 0644)

	cmd, cancel := suite.buildCommand(5*time.Second, "", nil, "convert", "--input", inputDir, "--output", outputDir, "--format", "webp", "--recursive")
	defer cancel()

	suite.NoError(cmd.Run())
	suite.DirExists(filepath.Join(outputDir, "valid_photos"))
	suite.NoDirExists(filepath.Join(outputDir, "ghost_docs"))
}

func (suite *E2ETestSuite) TestCLI_Recursive_ShouldProcessSubdirectoriesAndReplicateStructure() {
	inputDir := suite.T().TempDir()
	outputDir := suite.T().TempDir()

	subDir1 := filepath.Join(inputDir, "subdir1")
	subDir2 := filepath.Join(inputDir, "subdir2")
	nestedSubDir := filepath.Join(subDir1, "nested")
	os.Mkdir(subDir1, 0755)
	os.Mkdir(subDir2, 0755)
	os.Mkdir(nestedSubDir, 0755)

	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	realImages, _ := filepath.Glob(filepath.Join(testDataDir, "*.*"))
	suite.Require().NotEmpty(realImages)

	imgContent, _ := os.ReadFile(realImages[0])
	os.WriteFile(filepath.Join(inputDir, "root_photo.jpg"), imgContent, 0644)
	os.WriteFile(filepath.Join(subDir1, "photo1.jpg"), imgContent, 0644)
	os.WriteFile(filepath.Join(subDir2, "photo2.jpg"), imgContent, 0644)
	os.WriteFile(filepath.Join(nestedSubDir, "photo3.jpg"), imgContent, 0644)

	cmd, cancel := suite.buildCommand(15*time.Second, "", nil, "convert", "--input", inputDir, "--output", outputDir, "--format", "webp", "--recursive")
	defer cancel()

	_, err = cmd.CombinedOutput()
	suite.NoError(err)

	suite.FileExists(filepath.Join(outputDir, "root_photo.webp"))
	suite.FileExists(filepath.Join(outputDir, "subdir1", "photo1.webp"))
	suite.FileExists(filepath.Join(outputDir, "subdir2", "photo2.webp"))
	suite.FileExists(filepath.Join(outputDir, "subdir1", "nested", "photo3.webp"))
}

func (suite *E2ETestSuite) TestCLI_ShouldGenerateProfileFiles_WhenProfilingFlagsAreSet() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)

	outputDir := suite.T().TempDir()
	profileDir := suite.T().TempDir()

	cpuProfile := filepath.Join(profileDir, "cpu.prof")
	memProfile := filepath.Join(profileDir, "mem.prof")

	cmd, cancel := suite.buildCommand(10*time.Second, "", nil, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp", "--cpuprofile", cpuProfile, "--memprofile", memProfile)
	defer cancel()

	output, err := cmd.CombinedOutput()
	suite.NoError(err, "CLI failed to run with profiling flags. Output:\n%s", string(output))

	suite.FileExists(cpuProfile)
	suite.FileExists(memProfile)

	cpuInfo, err := os.Stat(cpuProfile)
	suite.Require().NoError(err)
	suite.Greater(cpuInfo.Size(), int64(0))

	memInfo, err := os.Stat(memProfile)
	suite.Require().NoError(err)
	suite.Greater(memInfo.Size(), int64(0))
}

func (suite *E2ETestSuite) TestCLI_Idempotency_ShouldSkipExistingFilesOnConsecutiveRuns() {
	testDataDir, err := filepath.Abs(filepath.Join("..", "..", "internal", "image", "testdata"))
	suite.Require().NoError(err)
	outputDir := suite.T().TempDir()

	cmd1, cancel1 := suite.buildCommand(10*time.Second, "", nil, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp")
	defer cancel1()
	output1, err := cmd1.CombinedOutput()
	suite.NoError(err, "First run failed: %s", string(output1))
	suite.NotContains(string(output1), "Skipped")

	generatedFiles, err := filepath.Glob(filepath.Join(outputDir, "*.webp"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(generatedFiles)

	fileStats := make(map[string]time.Time)
	for _, file := range generatedFiles {
		info, err := os.Stat(file)
		suite.Require().NoError(err)
		fileStats[file] = info.ModTime()
	}

	time.Sleep(1 * time.Second)

	cmd2, cancel2 := suite.buildCommand(10*time.Second, "", nil, "convert", "--input", testDataDir, "--output", outputDir, "--format", "webp")
	defer cancel2()
	output2, err := cmd2.CombinedOutput()
	suite.NoError(err, "Second run failed: %s", string(output2))
	suite.Contains(string(output2), "Skipped")

	for _, file := range generatedFiles {
		info, err := os.Stat(file)
		suite.Require().NoError(err)
		suite.Equal(fileStats[file], info.ModTime(), "File %s was modified on second run", file)
	}
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
