package cmd

import (
	"fmt"

	"github.com/ersanisk/sieve/internal/config"
	"github.com/ersanisk/sieve/internal/theme"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	themeName   string
	levelFilter string
	filterExpr  string
	follow      bool
)

// NewRootCmd creates the root cobra command.
func NewRootCmd(version, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sieve [flags] [file...]",
		Short: "Terminal JSON log viewer",
		Long:  "Sieve is a blazing-fast, terminal-based JSON log viewer with filtering, searching, and live tailing.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Override config with CLI flags
			if cmd.Flags().Changed("theme") {
				cfg.Theme = themeName
			}
			if cmd.Flags().Changed("level") {
				cfg.LevelFilter = levelFilter
			}
			if cmd.Flags().Changed("filter") {
				cfg.FilterExpr = filterExpr
			}
			if cmd.Flags().Changed("follow") {
				cfg.Follow = follow
			}
			cfg.FilePaths = args

			t := theme.Get(cfg.Theme)

			// TODO: Replace with Bubble Tea program in Phase 7
			fmt.Printf("sieve %s (built %s)\n", version, buildTime)
			fmt.Printf("Theme: %s\n", t.Name())
			if cfg.Follow {
				fmt.Println("Mode: follow")
			}
			if cfg.LevelFilter != "" {
				fmt.Printf("Level filter: %s\n", cfg.LevelFilter)
			}
			if cfg.FilterExpr != "" {
				fmt.Printf("Filter: %s\n", cfg.FilterExpr)
			}
			for _, path := range cfg.FilePaths {
				fmt.Printf("File: %s\n", path)
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", config.DefaultTheme, "color theme (monokai, dracula, gruvbox, nord)")
	rootCmd.Flags().StringVarP(&levelFilter, "level", "l", "", "minimum log level filter (debug, info, warn, error, fatal)")
	rootCmd.Flags().StringVar(&filterExpr, "filter", "", "filter expression (e.g., '.service == \"auth\"')")
	rootCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow file for new lines (like tail -f)")

	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, buildTime)

	return rootCmd
}

// Execute runs the root command.
func Execute(version, buildTime string) error {
	return NewRootCmd(version, buildTime).Execute()
}
