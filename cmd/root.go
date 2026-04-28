package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"charm.land/lipgloss/v2"
	"github.com/felipesimis/go-compactify-cli/internal/ui"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

var (
	concurrency int
	inputDir    string
	outputDir   string
	dryRun      bool
	Version     = "dev"
)

var rootCmd = &cobra.Command{
	Use:           "compactify",
	Short:         "Compactify: A versatile image compression and manipulation tool",
	Long:          `Compactify is your complete solution for optimizing images. With fast and intuitive commands, you can easily compress, resize, and convert your images, saving time and space.`,
	Version:       Version,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" {
			return nil
		}

		if inputDir == "" {
			return fmt.Errorf(ui.Error("required flag \"input\" (-i) not set"))
		}

		defaultWorkers := runtime.NumCPU()
		if concurrency > defaultWorkers*2 {
			fmt.Println(ui.Warn("⚠️  WARNING: Concurrency set very high. This may cause high memory usage and slow down your system."))
		}
		return nil
	},
}

func Execute() error {
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)
	defer bimg.Shutdown()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	displayVersion := Version
	if len(Version) > 0 && Version[0] == 'v' {
		displayVersion = Version[1:]
	}

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		Bold(true)

	rootCmd.SetVersionTemplate(fmt.Sprintf("Compactify %s\n", versionStyle.Render("v"+displayVersion)))
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	defaultWorkers := runtime.NumCPU()
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", defaultWorkers, "Number of concurrent operations")
	rootCmd.PersistentFlags().StringVarP(&inputDir, "input", "i", "", "Input directory containing the images to process")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "Output directory for processed images (default: auto-creates a sibling directory, e.g., '<input>-resized')")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without processing images, showing what would be done")
}
