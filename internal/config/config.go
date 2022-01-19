package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// Configuration stores Provider configuration.
type Configuration struct {
	Provider string
	Zone     string
	Exclude  []string
	Interval int
}

// New configuration from environment.
func Load() (*Configuration, error) {
	config := &Configuration{}
	if err := config.LoadFromEnv(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Configuration) LoadFromEnv() error {
	provider, found := os.LookupEnv("DDNS_PROVIDER")
	if !found {
		return errors.New("DDNS_PROVIDER not set")
	}
	c.Provider = provider
	log.Info().Msgf("DDNS_PROVIDER ==> %s", c.Provider)

	zone, found := os.LookupEnv("DDNS_ZONE")
	if !found {
		return errors.New("DDNS_ZONE not set")
	}
	c.Zone = zone
	log.Info().Msgf("DDNS_ZONE ==> %s", c.Zone)

	interval, found := os.LookupEnv("DDNS_UPDATE_INTERVAL")
	if found {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return fmt.Errorf("invalid update interval value: %s", interval)
		}
		c.Interval = i
	}

	return nil
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
