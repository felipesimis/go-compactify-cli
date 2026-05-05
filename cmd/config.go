package cmd

import (
	"github.com/spf13/cobra"
)

type AppConfig struct {
	Concurrency   int
	InputDir      string
	OutputDir     string
	DryRun        bool
	StripMetadata bool
}

func loadAppConfig(cmd *cobra.Command) AppConfig {
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	inputDir, _ := cmd.Flags().GetString("input")
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	stripMetadata, _ := cmd.Flags().GetBool("strip-metadata")

	return AppConfig{
		Concurrency:   concurrency,
		InputDir:      inputDir,
		OutputDir:     outputDir,
		DryRun:        dryRun,
		StripMetadata: stripMetadata,
	}
}
