package monitor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/provider"
	cf "github.com/milgradesec/ddns/internal/provider/cloudflare"
)

const (
	defaultInterval = 5 * time.Minute
)

// Monitor runs in a infinite loop and triggers provider zone updates
// every 3 min interval, can be triggered at any time by sending a
// SIGHUP signal.
type Monitor struct {
	ConfigFile string

	cfg *config.Configuration
	api provider.DNSProvider
}

// Start implements the service.Service interface.
func (m *Monitor) Start(s service.Service) error {
	cfg, err := config.New(m.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	m.cfg = cfg
	log.Infof("Configuration loaded from file: %s", m.ConfigFile)

	cfAPI, err := cf.New(cfg)
	if err != nil {
		return fmt.Errorf("error creating Cloudflare API client: %w", err)
	}
	m.api = cfAPI

	go func() {
		m.Run()
	}()
	return nil
}

// Run implements the service.Service interface.
func (m *Monitor) Run() {
	var interval time.Duration
	if m.cfg.Interval != 0 {
		interval = time.Duration(m.cfg.Interval) * time.Minute
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
				log.Infof("SIGHUP received: updating records for %s", m.cfg.Zone)
				m.callProvider()
			}
		}
	}()
	<-stop
}

func (m *Monitor) callProvider() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.api.UpdateZone(ctx); err != nil {
		log.Errorf("error updating zone %s: %v", m.cfg.Zone, err)
	}
}

// Stop implements the service.Service interface.
func (m *Monitor) Stop(s service.Service) error {
	return nil
}
