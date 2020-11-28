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
}

type APIAuthType int

const (
	APIToken APIAuthType = iota
	APIKey
)

// IsExcluded determines if a domain is excluded from changes.
func (c *Configuration) IsExcluded(s string) bool {
	for _, e := range c.Exclude {
		if s == e {
			return true
		}
	}
	return false
}

func (c *Configuration) isValid() (bool, error) {
	if c.Zone == "" {
		return false, errors.New("zone is empty")
	}
	if c.Email == "" {
		return false, errors.New("email is empty")
	}
	if c.APIKey == "" && c.APIToken == "" {
		return false, errors.New("apiKey and apiToken are empty")
	}
	return true, nil
}

func (c *Configuration) GetAuthType() APIAuthType {
	if c.APIKey == "" {
		return APIToken
	}
	return APIKey
}

// New configuration from file.
func New(file string) (cfg *Configuration, err error) {
	f, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	_, err = cfg.isValid()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
