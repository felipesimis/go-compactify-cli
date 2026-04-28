package main

import (
	"fmt"
	"os"

	"github.com/felipesimis/go-compactify-cli/cmd"
	"github.com/felipesimis/go-compactify-cli/internal/ui"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, ui.Error("Error: %v\n"), err)
		os.Exit(1)
	}
}
