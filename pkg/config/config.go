package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is the configuration of the `huectl` CLI.
type Config struct {
	BridgeID        string `yaml:"bridge_id"`
	BridgeURL       string `yaml:"bridge_url"`
	ClientID        string `yaml:"client_id"`
	CertFingerprint string `yaml:"cert_fingerprint"`
}

// Read reads the CLI configuration from the user configuration directory.
func Read() (*Config, error) {
	cfgPath, err := AbsolutePath()
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, err
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	return &cfg, nil
}

// AbsolutePath returns the absolute path of the configuration file.
func AbsolutePath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to find user config dir: %w", err)
	}

	return filepath.Join(cfgDir, "huectl", "config.yml"), nil
}
