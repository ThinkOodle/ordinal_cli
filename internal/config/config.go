// Package config handles loading configuration from CLI flags, environment
// variables, and the config file at ~/.config/ordinal/config.yaml.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

const (
	// EnvPrefix is the prefix for environment variables.
	EnvPrefix = "ORDINAL"

	// DefaultOutputFormat is the default output format.
	DefaultOutputFormat = "json"

	// configDirName is the directory name under ~/.config/.
	configDirName = "ordinal"

	// configFileName is the config file name without extension.
	configFileName = "config"

	// configFileType is the config file format.
	configFileType = "yaml"
)

// Config holds the application configuration.
type Config struct {
	APIKey       string `mapstructure:"api_key" yaml:"api_key,omitempty"`
	OutputFormat string `mapstructure:"output_format" yaml:"output_format,omitempty"`
	Verbose      bool   `mapstructure:"verbose" yaml:"verbose,omitempty"`
}

// ConfigDir returns the path to the configuration directory.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determining home directory: %w", err)
	}
	return filepath.Join(home, ".config", configDirName), nil
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName+"."+configFileType), nil
}

// Load reads configuration from the config file and environment variables.
// CLI flags should override these values after Load is called.
//
// A malformed or unreadable config file is non-fatal: Load still returns a
// *Config populated from env vars and defaults, along with the read error so
// the caller can surface a warning. This preserves ORDINAL_* env values
// (most importantly ORDINAL_API_KEY) when the on-disk file is broken.
func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("output_format", DefaultOutputFormat)
	v.SetDefault("verbose", false)

	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()

	v.BindEnv("api_key", "ORDINAL_API_KEY")
	v.BindEnv("output_format", "ORDINAL_OUTPUT_FORMAT")
	v.BindEnv("verbose", "ORDINAL_VERBOSE")

	var readErr error
	configDir, err := ConfigDir()
	if err == nil {
		v.AddConfigPath(configDir)
		v.SetConfigName(configFileName)
		v.SetConfigType(configFileType)

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				configFile := filepath.Join(configDir, configFileName+"."+configFileType)
				if _, statErr := os.Stat(configFile); statErr == nil {
					readErr = fmt.Errorf("reading config file: %w", err)
				}
			}
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, readErr
}

// SaveAPIKey writes the API key to the config file, preserving any other
// existing settings. Creates the config directory and file if they don't exist.
func SaveAPIKey(apiKey string) error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	configFile, err := ConfigFilePath()
	if err != nil {
		return err
	}

	cfg := &Config{}
	data, err := os.ReadFile(configFile)
	switch {
	case err == nil:
		if err := yaml.Unmarshal(data, cfg); err != nil {
			// Refuse to overwrite a file we can't parse — otherwise the
			// caller's other settings get silently wiped on save.
			return fmt.Errorf("parsing existing config at %s (fix or remove the file before saving): %w", configFile, err)
		}
	case !os.IsNotExist(err):
		return fmt.Errorf("reading existing config: %w", err)
	}

	cfg.APIKey = apiKey

	data, err = yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key, checking the provided flag value first, then
// falling back to the config. Returns an error if no API key is found.
func GetAPIKey(flagValue string, cfg *Config) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if cfg != nil && cfg.APIKey != "" {
		return cfg.APIKey, nil
	}
	return "", fmt.Errorf("api key is required: set via --api-key flag, ORDINAL_API_KEY env var, or `ordinal auth`")
}
