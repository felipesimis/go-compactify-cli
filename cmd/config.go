package cmd

import (
	"github.com/spf13/cobra"
)

type GlobalConfig struct {
	Concurrency int
	InputDir    string
	OutputDir   string
	DryRun      bool
}

func loadGlobalConfig(cmd *cobra.Command) GlobalConfig {
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	inputDir, _ := cmd.Flags().GetString("input")
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	return GlobalConfig{
		Concurrency: concurrency,
		InputDir:    inputDir,
		OutputDir:   outputDir,
		DryRun:      dryRun,
	}
}
