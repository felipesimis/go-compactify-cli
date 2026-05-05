package cmd

import (
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/spf13/cobra"
)

func addEncodingFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("quality", "q", image.DefaultQuality, "Compression quality (1-100)")
}
