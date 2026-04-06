package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

var (
	concurrency int
)

var rootCmd = &cobra.Command{
	Use:           "compactify",
	Short:         "Compactify: A versatile image compression and manipulation tool",
	Long:          `Compactify is your complete solution for optimizing images. With fast and intuitive commands, you can easily compress, resize, and convert your images, saving time and space.`,
	SilenceErrors: true,
}

func Execute() error {
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)
	defer bimg.Shutdown()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 20, "Number of concurrent operations")
}
