package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ersanisk/sieve/internal/app"
	"github.com/ersanisk/sieve/internal/config"
	"github.com/ersanisk/sieve/internal/theme"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	themeName string
	follow    bool
)

// NewRootCmd creates the root cobra command.
func NewRootCmd(version, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sieve [file...]",
		Short: "Terminal JSON log viewer",
		Long:  "Sieve is a blazing-fast, terminal-based JSON log viewer with filtering, searching, and live tailing.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var filePath string
			if len(args) > 0 {
				filePath = args[0]
			}

			appCfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if cmd.Flags().Changed("theme") {
				appCfg.Theme = themeName
			}
			if cmd.Flags().Changed("follow") {
				appCfg.Follow = follow
			}

			if theme.Get(appCfg.Theme) == nil {
				appCfg.Theme = "default"
			}

			model := app.NewModel(filePath, appCfg.Theme, appCfg.Follow)
			program := tea.NewProgram(model)

			if _, err := program.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", config.DefaultTheme, "color theme (monokai, dracula, gruvbox, nord)")
	rootCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow file for new lines (like tail -f)")

	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, buildTime)

	return rootCmd
}
