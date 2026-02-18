package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Theme           string   `mapstructure:"theme"`
	TimestampFormat string   `mapstructure:"timestamp_format"`
	ShowLineNumbers bool     `mapstructure:"show_line_numbers"`
	WrapLines       bool     `mapstructure:"wrap_lines"`
	JSONIndent      int      `mapstructure:"json_indent"`
	MaxBufferSize   int      `mapstructure:"max_buffer_size"`
	WorkerCount     int      `mapstructure:"worker_count"`
	Follow          bool     `mapstructure:"follow"`
	LevelFilter     string   `mapstructure:"level_filter"`
	FilterExpr      string   `mapstructure:"filter_expr"`
	FilePaths       []string `mapstructure:"-"`
}

// Load reads configuration from file and environment, applying defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("theme", DefaultTheme)
	v.SetDefault("timestamp_format", DefaultTimestampFormat)
	v.SetDefault("show_line_numbers", DefaultShowLineNumbers)
	v.SetDefault("wrap_lines", DefaultWrapLines)
	v.SetDefault("json_indent", DefaultJSONIndent)
	v.SetDefault("max_buffer_size", DefaultMaxBufferSize)
	v.SetDefault("worker_count", DefaultWorkerCount)
	v.SetDefault("follow", false)

	// Config file search paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if configDir, err := os.UserConfigDir(); err == nil {
		v.AddConfigPath(filepath.Join(configDir, "sieve"))
	}
	v.AddConfigPath(".")

	// Environment variables
	v.SetEnvPrefix("SIEVE")
	v.AutomaticEnv()

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	return cfg, nil
}
