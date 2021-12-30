package monitor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/provider"
	cf "github.com/milgradesec/ddns/internal/provider/cloudflare"
)

const (
	defaultInterval = 5 * time.Minute
)

// Monitor runs in a infinite loop and triggers provider zone updates
// every 5 min interval, can be triggered at any time by sending a
// SIGHUP signal.
type Monitor struct {
	config   *config.Configuration
	provider provider.DNSProvider
}

// Start implements the service.Service interface.
func (m *Monitor) Start(s service.Service) error {
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	m.config = config

	log.Info().Msgf("Using %s provider", config.Provider)

	cfAPI, err := cf.New()
	if err != nil {
		return fmt.Errorf("error creating Cloudflare API client: %w", err)
	}
	m.provider = cfAPI

	go func() {
		m.Run()
	}()
	return nil
}

// Run implements the service.Service interface.
func (m *Monitor) Run() {
	var interval time.Duration
	if m.config.Interval != 0 {
		interval = time.Duration(m.config.Interval) * time.Minute
	} else {
		interval = defaultInterval
	}

	ticker := time.NewTicker(interval)
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	stop := make(chan bool)
	go func() {
		m.callProvider()
		for {
			select {
			case <-ticker.C:
				m.callProvider()

			case <-sighup:
				log.Info().Msgf("SIGHUP: updating records for %s", m.config.Zone)
				m.callProvider()
			}
		}
	}()
	<-stop
}

func (m *Monitor) callProvider() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.provider.UpdateZone(ctx); err != nil {
		log.Error().Msgf("error updating zone %s: %v", m.config.Zone, err)
	}
}

// Stop implements the service.Service interface.
func (m *Monitor) Stop(s service.Service) error {
	return nil
}
