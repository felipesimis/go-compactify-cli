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
