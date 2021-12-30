package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
	cfg := &Configuration{}
	if err := cfg.LoadFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
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

// type APIAuthType int

// const (
// 	APIToken APIAuthType = iota
// 	APIKey
// )

// LoadFromEnv reads the API Key or Token from environment variables.
func (c *Configuration) LoadFromEnv() error {
	provider, found := os.LookupEnv("DDNS_PROVIDER")
	if !found {
		return errors.New("DDNS_PROVIDER not set")
	}
	c.Provider = provider

	zone, found := os.LookupEnv("DDNS_ZONE")
	if !found {
		return errors.New("DDNS_ZONE not set")
	}
	c.Zone = zone

	interval, found := os.LookupEnv("DDNS_UPDATE_INTERVAL")
	if found {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return fmt.Errorf("invalid update interval value: %s", interval)
		}
		c.Interval = i
	}

	// key, found := os.LookupEnv("CLOUDFLARE_API_KEY")
	// if found {
	// 	c.APIKey = key
	// }

	// token, found := os.LookupEnv("CLOUDFLARE_API_TOKEN")
	// if found {
	// 	c.APIToken = token
	// }

	// keyFile, found := os.LookupEnv("CLOUDFLARE_API_KEY_FILE")
	// if found {
	// 	buf, err := ioutil.ReadFile(keyFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	c.APIKey = strings.TrimSpace(string(buf))
	// }

	// tokenFile, found := os.LookupEnv("CLOUDFLARE_API_TOKEN_FILE")
	// if found {
	// 	buf, err := ioutil.ReadFile(tokenFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	c.APIToken = strings.TrimSpace(string(buf))
	// }
	return nil
}

// func (c *Configuration) isValid() (bool, error) {
// 	if c.Zone == "" {
// 		return false, errors.New("zone is empty")
// 	}
// 	if c.Email == "" {
// 		return false, errors.New("email is empty")
// 	}
// 	if c.APIKey == "" && c.APIToken == "" {
// 		return false, errors.New("no APIKey or APIToken provided")
// 	}
// 	return true, nil
// }

// func (c *Configuration) GetAuthType() APIAuthType {
// 	if c.APIKey == "" {
// 		return APIToken
// 	}
// 	return APIKey
// }

// New configuration from file.
// func New(file string) (cfg *Configuration, err error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	err = json.NewDecoder(f).Decode(&cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if cfg.APIKey == "" && cfg.APIToken == "" {
// 		err = cfg.LoadFromEnv()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	_, err = cfg.isValid()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cfg, nil
// }
