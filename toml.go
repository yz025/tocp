package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Pairs            []Pair `toml:"pairs"`
	Log              bool   `toml:"log"`
	RemoveBeforeCopy bool   `toml:"remove_before_copy"`
}
type Pair struct {
	Src string `toml:"src"`
	Dst string `toml:"dst"`
}

func LoadConfig(givenPath string) (*Config, error) {
	path, found := findConfig(givenPath)
	if !found {
		return nil, errors.New("Failed to find a config")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the config %q: %w", path, err)
	}

	config := Config{
		Log:              true,
		RemoveBeforeCopy: true,
	}
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Failed to parse the config: %w", err)
	}

	return &config, nil
}
func findConfig(pathFlag string) (string, bool) {
	path, found := checkGivenConfig(pathFlag)
	if found {
		return path, true
	}
	path, found = checkLocalConfig()
	if found {
		return path, true
	}
	path, found = checkGlobalConfig()
	if found {
		return path, true
	}

	return "", false
}
func checkGivenConfig(pathFlag string) (string, bool) {
	if pathFlag == "" {
		return "", false
	}

	abs, err := filepath.Abs(pathFlag)
	if err != nil {
		return "", false
	}

	if _, err := os.Stat(abs); err != nil {
		return "", false
	}

	return pathFlag, true
}
func checkLocalConfig() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	localConfig := filepath.Join(cwd, "tocp.toml")
	if _, err := os.Stat(localConfig); err != nil {
		return "", false
	}

	return localConfig, true
}
func checkGlobalConfig() (string, bool) {
	home := os.Getenv("TOCP_CONFIG_HOME")
	if home == "" {
		return "", false
	}

	globalConfig := filepath.Join(home, "tocp.toml")
	if _, err := os.Stat(globalConfig); err != nil {
		return "", false
	}

	return globalConfig, true
}
