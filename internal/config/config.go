package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Configuration stores Provider configuration.
type Configuration struct {
	Provider string   `json:"provider"`
	Zone     string   `json:"zone"`
	Email    string   `json:"email"`
	APIKey   string   `json:"apikey"`
	APIToken string   `json:"apitoken"`
	Exclude  []string `json:"exclude"`
	Interval int      `json:"interval"`
}

// IsExcluded determines if a domain is excluded from changes.
func (c *Configuration) IsExcluded(s string) bool {
	for _, e := range c.Exclude {
		if s == e {
			return true
		}
	}
	return false
}

// LoadFromEnv reads configuration from environment variables.
func (c *Configuration) LoadFromEnv() error {
	provider, found := os.LookupEnv("DDNS_PROVIDER")
	if !found {
		return errors.New("no provider is configured")
	}
	c.Provider = provider

	return nil
}

// New configuration from file.
func New(file string) (cfg *Configuration, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
