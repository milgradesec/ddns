package config

import (
	"encoding/json"
	"os"
	"reflect"
)

// Config stores Provider configuration
type Config struct {
	Provider string   `json:"provider"`
	Zone     string   `json:"zone"`
	Email    string   `json:"email"`
	APIKey   string   `json:"apikey"`
	Exclude  []string `json:"exclude"`
}

// IsEmpty checks if Config struct is empty
func (cfg Config) IsEmpty() bool {
	return reflect.DeepEqual(cfg, Config{})
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
