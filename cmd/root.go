package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"charm.land/lipgloss/v2"
	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/ui"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	cfgFile string
)

func NewRootCmd() *cobra.Command {
	defaultWorkers := runtime.NumCPU()

	cmd := &cobra.Command{
		Use:           "compactify",
		Short:         "Compactify: A versatile image compression and manipulation tool",
		Long:          `Compactify is your complete solution for optimizing images. With fast and intuitive commands, you can easily compress, resize, and convert your images, saving time and space.`,
		Version:       Version,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initConfig(cmd)

			viper.BindPFlags(cmd.Flags())

			isHelp := cmd.Name() == "help" || cmd.Flags().Changed("help")
			isInit := cmd.Name() == "init"
			isVersion := cmd.Flags().Changed("version")
			if isHelp || isInit || isVersion {
				return nil
			}

			cfg := loadAppConfig()
			if cfg.InputDir == "" {
				return errors.New("required flag \"input\" (-i) not set")
			}

			if cfg.Concurrency > defaultWorkers*2 {
				fmt.Fprintln(cmd.OutOrStdout(), ui.Warn("WARNING: Concurrency set very high. This may cause high memory usage and slow down your system."))
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml or $HOME/.config/compactify/config.yaml)")
	cmd.PersistentFlags().IntP("concurrency", "c", defaultWorkers, "Maximum number of images to process in parallel")
	cmd.PersistentFlags().StringP("input", "i", "", "Input directory containing images to process")
	cmd.PersistentFlags().StringP("output", "o", "", "Output directory (default: auto-creates sibling directory)")
	cmd.PersistentFlags().Bool("dry-run", false, "Preview operations without modifying files")
	cmd.PersistentFlags().Bool("strip-metadata", false, "Strip EXIF data for privacy (GPS, camera info) and reduced file size")

	viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func Execute() error {
	image.InitializeProcessor()
	defer image.ShutdownProcessor()
	fs := filesystem.NewFileSystem()

	rootCmd := NewRootCmd()
	rootCmd.AddCommand(
		NewInitCmd(fs),
		NewPaletteCmd(fs, image.NewProcessor),
		NewGrayscaleCmd(fs, image.NewProcessor),
		NewFlipCmd(fs, image.NewProcessor),
		NewThumbnailCmd(fs, image.NewProcessor),
		NewResizeCmd(fs, image.NewProcessor),
		NewEnlargeCmd(fs, image.NewProcessor),
		NewConvertCmd(fs, image.NewProcessor),
		NewCropCmd(fs, image.NewProcessor),
		NewLosslessCmd(fs, image.NewProcessor),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	displayVersion := Version
	if len(Version) > 0 && Version[0] == 'v' {
		displayVersion = Version[1:]
	}

	versionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true)
	rootCmd.SetVersionTemplate(fmt.Sprintf("Compactify %s\n", versionStyle.Render("v"+displayVersion)))

	return rootCmd.ExecuteContext(ctx)
}

func initConfig(cmd *cobra.Command) {
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

	viper.SetEnvPrefix("COMPACTIFY")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(cmd.OutOrStderr(), "%s: %v\n", ui.Error("Error reading config file"), err)
		}
	}
}
