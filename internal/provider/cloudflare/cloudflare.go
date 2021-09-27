package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/pkg/ip"

	httpc "github.com/milgradesec/go-libs/http"
	log "github.com/sirupsen/logrus"
)

// API implements provider.API interface.
type API struct {
	Zone   string
	Email  string
	zoneID string
	Key    string
	Token  string
	cf     *cloudflare.API
	cfg    *config.Configuration
}

// New creates a Cloudflare DNS provider.
func New(cfg *config.Configuration) (*API, error) { //nolint
	var (
		zone     string
		email    string
		apiKey   string
		apiToken string
	)

	zone, found := os.LookupEnv("CLOUDFLARE_ZONE")
	if !found {
		return nil, errors.New("cloudflare zone not set")
	}

	email, found = os.LookupEnv("CLOUDFLARE_MAIL") // not needed for token
	if !found {
		return nil, errors.New("cloudflare mail not set")
	}

	apiKey, found = os.LookupEnv("CLOUDFLARE_API_KEY")
	if !found {
		keyFile, found := os.LookupEnv("CLOUDFLARE_API_KEY_FILE")
		if found {
			buf, err := ioutil.ReadFile(keyFile)
			if err != nil {
				return nil, err
			}
			apiKey = strings.TrimSpace(string(buf))
		}
	}

	if apiKey != "" {
		cfAPI, err := newWithAPIKey(apiKey, email)
		if err != nil {
			return nil, err
		}

		return &API{
			Zone:  zone,
			Email: email,
			Key:   apiKey,
			cf:    cfAPI,
		}, nil
	}

	apiToken, found = os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if !found {
		tokenFile, found := os.LookupEnv("CLOUDFLARE_API_TOKEN_FILE")
		if found {
			buf, err := ioutil.ReadFile(tokenFile)
			if err != nil {
				return nil, err
			}
			apiToken = strings.TrimSpace(string(buf))
		}
	}

	if apiToken != "" {
		cfAPI, err := newWithAPIToken(apiToken)
		if err != nil {
			return nil, err
		}

		return &API{
			Zone:  zone,
			Email: email,
			Token: apiKey,
			cf:    cfAPI,
		}, nil
	}

	return nil, errors.New("no Cloudflare API credentials found")
}

func newWithAPIKey(key string, email string) (*cloudflare.API, error) {
	return cloudflare.New(key, email, cloudflare.HTTPClient(httpc.NewHTTPClient()))
}

func newWithAPIToken(token string) (*cloudflare.API, error) {
	return cloudflare.NewWithAPIToken(token, cloudflare.HTTPClient(httpc.NewHTTPClient()))
}

// Name implements the provider.API interface.
func (api *API) Name() string {
	return "Cloudflare"
}

// UpdateZone implements the provider.API interface.
func (api *API) UpdateZone(ctx context.Context) error {
	if api.zoneID == "" {
		id, err := api.cf.ZoneIDByName(api.Zone)
		if err != nil {
			return fmt.Errorf("cloudflare api error: failed to retrieve zone id: %w", err)
		}
		api.zoneID = id
	}

	publicIP, err := ip.GetIP(ctx)
	if err != nil {
		return err
	}

	records, err := api.cf.DNSRecords(ctx, api.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("cloudflare api error: failed to list dns records: %w", err)
	}

	for i := range records {
		r := records[i]

		if api.cfg.IsExcluded(r.Name) {
			continue
		}

		if r.Type == "A" {
			if r.Content != publicIP {
				rr := cloudflare.DNSRecord{
					Type:    r.Type,
					Name:    r.Name,
					Content: publicIP,
					Proxied: r.Proxied,
				}

				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				if err := api.cf.UpdateDNSRecord(ctx, api.zoneID, r.ID, rr); err != nil {
					return fmt.Errorf("error updating %s: %w", r.Name, err)
				}
				log.Infof("updated %s from %s to %s", r.Name, r.Content, publicIP)
			}
		}
	}
	return nil
}
