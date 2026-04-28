package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupInitTest(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)

	viper.Reset()

	t.Cleanup(func() {
		os.Chdir(oldWd)
	})
	return tmpDir
}

func TestShould_CreateDefaultConfig_When_FileDoesNotExist(t *testing.T) {
	setupInitTest(t)
	configPath := "config.yaml"

	rootCmd.SetArgs([]string{"init"})
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.FileExists(t, configPath, "file config.yaml should be created")

	content, _ := os.ReadFile(configPath)
	assert.Contains(t, string(content), "concurrency:", "The file should contain the concurrency key")
}

func TestShould_ReturnError_When_FileAlreadyExists(t *testing.T) {
	setupInitTest(t)
	configPath := "config.yaml"
	importantContent := "user-custom-config: true"

	err := os.WriteFile(configPath, []byte(importantContent), 0644)
	assert.NoError(t, err)

	rootCmd.SetArgs([]string{"init"})
	err = rootCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists", "The error should inform that the file already exists")

	content, _ := os.ReadFile(configPath)
	assert.Equal(t, importantContent, string(content), "The init command should never overwrite an existing file without permission")
}

func TestShould_ReturnError_When_WriteFileFails(t *testing.T) {
	setupInitTest(t)
	configPath := "config.yaml"

	err := os.Mkdir(configPath, 0755)
	assert.NoError(t, err)

	rootCmd.SetArgs([]string{"init", "--force"})
	err = rootCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create config file", "should capture the file write failure")
}

func TestShould_OverwriteConfig_When_FileExistsAndForceFlagIsUsed(t *testing.T) {
	setupInitTest(t)
	configPath := "config.yaml"
	oldContent := "profile: old"

	os.WriteFile(configPath, []byte(oldContent), 0644)

	rootCmd.SetArgs([]string{"init", "--force"})
	err := rootCmd.Execute()

	assert.NoError(t, err)

	newContent, _ := os.ReadFile(configPath)
	assert.NotEqual(t, oldContent, string(newContent), "The content should be overwritten when --force is used")
	assert.Contains(t, string(newContent), "concurrency:", "The new file should contain the default concurrency key")
}
