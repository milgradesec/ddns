package cloudflare

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/pkg/ip"
)

// API implements provider.API interface.
type API struct {
	cfg *config.Configuration

	api *cloudflare.API
	id  string
}

// New creates a Cloudflare DNS provider.
func New(cfg *config.Configuration) (*API, error) {
	if cfg.GetAuthType() == config.APIKey {
		return newWithAPIKey(cfg)
	}
	return newWithAPIToken(cfg)
}

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			},
		},
		Timeout: 15 * time.Second,
	}
}

// Creates a Clouflare Provider using API Key.
func newWithAPIKey(cfg *config.Configuration) (*API, error) {
	api, err := cloudflare.New(cfg.APIKey, cfg.Email, cloudflare.HTTPClient(newHttpClient()))
	if err != nil {
		return nil, err
	}

	cf := &API{
		api: api,
		cfg: cfg,
	}
	return cf, nil
}

// Creates a Cloudflare Provider using API Token.
func newWithAPIToken(cfg *config.Configuration) (*API, error) {
	api, err := cloudflare.NewWithAPIToken(cfg.APIToken, cloudflare.HTTPClient(newHttpClient()))
	if err != nil {
		return nil, err
	}

	cf := &API{
		api: api,
		cfg: cfg,
	}
	return cf, nil
}

// Name implements the provider.API interface.
func (cf *API) Name() string {
	return "Cloudflare"
}

// UpdateZone implements the provider.API interface.
func (cf *API) UpdateZone() error {
	if cf.id == "" {
		id, err := cf.api.ZoneIDByName(cf.cfg.Zone)
		if err != nil {
			return errors.New("cloudflare api error: failed to retrieve zone id")
		}
		cf.id = id
	}

	publicIP, err := ip.GetIP()
	if err != nil {
		return err
	}

	records, err := cf.api.DNSRecords(cf.id, cloudflare.DNSRecord{})
	if err != nil {
		return errors.New("cloudflare api error: failed to list dns records")
	}

	for i := range records {
		r := records[i]

		if cf.cfg.IsExcluded(r.Name) {
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
				if err := cf.api.UpdateDNSRecord(cf.id, r.ID, rr); err != nil {
					return fmt.Errorf("error updating %s: %v", r.Name, err)
				}
				log.Infof("updated %s from %s to %s", r.Name, r.Content, publicIP)
			}
		}
	}
	return nil
}
