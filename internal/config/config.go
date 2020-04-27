package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config stores Provider configuration
type Config struct {
	Provider string `json:"provider"`
	Zone     string `json:"zone"`
	Email    string `json:"email"`
	APIKey   string `json:"apikey"`
}

// Load config from file
func Load(file string) (Config, error) {
	var cfg Config

	f, err := os.Open(file)
	if err != nil {
		return cfg, fmt.Errorf("failed to open %s: %v", file, err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to load config from %s: %v", file, err)
	}
	return cfg, nil
}
