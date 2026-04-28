package cmd

import (
	_ "embed"
	"fmt"
	"runtime"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/templates"
	"github.com/spf13/cobra"
)

func initRun(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	configPath := "config.yaml"

	fs := filesystem.NewFileSystem()
	_, err := fs.OpenFile(configPath)
	if err == nil && !force {
		return fmt.Errorf("a configuration file already exists at '%s'. Use --force to overwrite it", configPath)
	}

	content := fmt.Sprintf(templates.ConfigTemplate, runtime.NumCPU())
	err = fs.WriteFile(configPath, []byte(content))
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Println("✓ Configuration file initialized successfully: " + configPath)
	return nil
}

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "config"},
	Args:    cobra.NoArgs,
	Short:   "Initialize a default configuration file",
	Long: `Create a 'config.yaml' file in the current directory with default settings.
This allows you to persist settings like concurrency and default directories 
without having to pass flags every time you run a command.`,
	RunE: initRun,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing config.yaml file")
}
