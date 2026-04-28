package main

import (
	"fmt"
	"os"

	"github.com/felipesimis/go-compactify-cli/cmd"
	"github.com/felipesimis/go-compactify-cli/internal/ui"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, ui.Error(err.Error()))
		os.Exit(1)
	}
}
