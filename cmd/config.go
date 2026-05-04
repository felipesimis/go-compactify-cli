package cmd

import (
	"github.com/spf13/cobra"
)

type GlobalConfig struct {
	Concurrency   int
	InputDir      string
	OutputDir     string
	DryRun        bool
	StripMetadata bool
}

func loadGlobalConfig(cmd *cobra.Command) GlobalConfig {
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	inputDir, _ := cmd.Flags().GetString("input")
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	stripMetadata, _ := cmd.Flags().GetBool("strip-metadata")

	return GlobalConfig{
		Concurrency:   concurrency,
		InputDir:      inputDir,
		OutputDir:     outputDir,
		DryRun:        dryRun,
		StripMetadata: stripMetadata,
	}
}
