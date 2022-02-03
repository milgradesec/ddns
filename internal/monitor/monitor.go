package monitor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kardianos/service"
	"github.com/rs/zerolog/log"

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

	interval time.Duration
	stop     chan struct{}
}

func New(config *config.Configuration) *Monitor {
	var interval time.Duration
	if config.Interval != 0 {
		interval = time.Duration(config.Interval) * time.Minute
	} else {
		interval = defaultInterval
	}

	return &Monitor{
		config:   config,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start implements the service.Service interface.
func (m *Monitor) Start(s service.Service) error {
	m.stop = make(chan struct{})

	cloudflareDNS, err := cf.New(m.config)
	if err != nil {
		return fmt.Errorf("error creating cloudflare API client: %w", err)
	}
	m.provider = cloudflareDNS

	go func() {
		m.Run()
	}()
	return nil
}

// Run implements the service.Service interface.
func (m *Monitor) Run() {
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	ticker := time.NewTicker(m.interval)

	m.providerUpdateZone()
	for {
		select {
		case <-ticker.C:
			m.providerUpdateZone()

		case <-sighup:
			log.Info().Str("provider", m.provider.Name()).Str("zone", m.provider.GetZoneName()).Msg("SIGHUP received, updating records")
			m.providerUpdateZone()

		case <-m.stop:
			ticker.Stop()
			return
		}
	}
}

func (m *Monitor) providerUpdateZone() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.provider.UpdateZone(ctx); err != nil {
		log.Error().Err(err).Str("provider", m.provider.Name()).Str("zone", m.provider.GetZoneName()).Msg("error updating zone")
	}
}

// Stop implements the service.Service interface.
func (m *Monitor) Stop(s service.Service) error {
	log.Info().Msg("Stopping service")
	close(m.stop)
	return nil
}
