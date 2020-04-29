package config

import (
	"encoding/json"
	"os"
)

// Config stores Provider configuration
type Config struct {
	Provider string   `json:"provider"`
	Zone     string   `json:"zone"`
	Email    string   `json:"email"`
	APIKey   string   `json:"apikey"`
	Exclude  []string `json:"exclude"`
}

func (c Config) IsExcluded(s string) bool {
	for _, e := range c.Exclude {
		if s == e {
			return true
		}
	}
	return false
}

// Load config from file
func Load(file string) (Config, error) {
	var cfg Config

	f, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
