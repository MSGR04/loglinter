package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func FindConfig() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		configPath := filepath.Join(dir, ".loglinter.json")
		if _, err := os.Stat(configPath); err == nil {
			return LoadConfig(configPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	defaultConfig := Config{
		EnableLowercase:     true,
		EnableEnglishOnly:   true,
		EnableSpecialChars:  true,
		EnableSensitiveData: true,
		SensitivePatterns:   []string{},
	}
	return &defaultConfig, nil
}

func DefaultConfig() Config {
	return Config{
		EnableLowercase:     true,
		EnableEnglishOnly:   true,
		EnableSpecialChars:  true,
		EnableSensitiveData: true,
		SensitivePatterns:   []string{},
	}
}
