package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "compactify",
	Short: "Compactify: A versatile image compression and manipulation tool",
	Long: `Compactify is a powerful CLI tool focused on compressing images efficiently.
In addition to compression, it provides various commands for image manipulation,
including resizing, cropping, converting, and more. Ideal for optimizing images
for web and other uses.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
