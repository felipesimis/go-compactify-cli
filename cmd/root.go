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
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:           "compactify",
	Short:         "Compactify: A versatile image compression and manipulation tool",
	Long:          `Compactify is your complete solution for optimizing images. With fast and intuitive commands, you can easily compress, resize, and convert your images, saving time and space.`,
	Version:       Version,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		isHelp := cmd.Name() == "help" || cmd.Flags().Changed("help")
		isVersion := cmd.Flags().Changed("version")
		if isHelp || isVersion {
			return nil
		}

		cfg := loadGlobalConfig(cmd)
		if cfg.InputDir == "" {
			return fmt.Errorf(ui.Error("required flag \"input\" (-i) not set"))
		}

		defaultWorkers := runtime.NumCPU()
		if cfg.Concurrency > defaultWorkers*2 {
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

	cobra.OnInitialize(initConfig)
	return rootCmd.ExecuteContext(ctx)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")

		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(utils.GetConfigDir(home))
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("COMPACTIFY")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, ui.Error("Error reading config file: %v\n"), err)
		}
	}

	bindFlags(rootCmd)
}

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed && viper.IsSet(flag.Name) {
			value := viper.Get(flag.Name)
			cmd.Flags().Set(flag.Name, fmt.Sprintf("%v", value))
		}
	})

	for _, child := range cmd.Commands() {
		bindFlags(child)
	}
}

func init() {
	defaultWorkers := runtime.NumCPU()
	rootCmd.PersistentFlags().IntP("concurrency", "c", defaultWorkers, "Number of concurrent operations")
	rootCmd.PersistentFlags().StringP("input", "i", "", "Input directory containing the images to process")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output directory for processed images (default: auto-creates a sibling directory, e.g., '<input>-resized')")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Perform a dry run without processing images, showing what would be done")
}
