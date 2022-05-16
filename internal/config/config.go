package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	log.Debug().Msg("Loading configuration from environment variables")
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
	log.Debug().Str("name", "DDNS_PROVIDER").Str("value", provider).Msg("Found environment variable")

	zone, found := os.LookupEnv("DDNS_ZONE")
	if !found {
		return errors.New("DDNS_ZONE not set")
	}
	c.Zone = zone
	log.Debug().Str("name", "DDNS_ZONE").Str("value", zone).Msg("Found environment variable")

	interval, found := os.LookupEnv("DDNS_UPDATE_INTERVAL")
	if found {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return fmt.Errorf("invalid update interval value: %s", interval)
		}
		c.Interval = i
		log.Debug().Str("name", "DDNS_UPDATE_INTERVAL").Str("value", interval).Msg("Found environment variable")
	}

	exclusions, found := os.LookupEnv("DDNS_EXCLUDE")
	if found {
		list := strings.Split(exclusions, ",")
		if len(list) > 0 {
			c.Exclude = list
		}
		log.Debug().Str("name", "DDNS_EXCLUDE").Str("value", exclusions).Msg("Found environment variable")
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
