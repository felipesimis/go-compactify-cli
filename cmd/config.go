package cmd

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Concurrency   int
	InputDir      string
	OutputDir     string
	DryRun        bool
	StripMetadata bool
	Quality       int
	Recursive     bool
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Concurrency:   viper.GetInt("concurrency"),
		InputDir:      viper.GetString("input"),
		OutputDir:     viper.GetString("output"),
		DryRun:        viper.GetBool("dry-run"),
		StripMetadata: viper.GetBool("strip-metadata"),
		Quality:       viper.GetInt("quality"),
		Recursive:     viper.GetBool("recursive"),
	}
}
