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
