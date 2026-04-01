package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	concurrency int
)

var rootCmd = &cobra.Command{
	Use:   "compactify",
	Short: "Compactify: A versatile image compression and manipulation tool",
	Long:  `Compactify is your complete solution for optimizing images. With fast and intuitive commands, you can easily compress, resize, and convert your images, saving time and space.`,
}

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 20, "Number of concurrent operations")
}
