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
// every 5 min interval, can be triggered at any time by sending a
// SIGHUP signal.
type Monitor struct {
	cfg      *config.Configuration
	provider provider.DNSProvider
}

func New(cfg *config.Configuration) *Monitor {
	return &Monitor{
		cfg: cfg,
	}
}

// Start implements the service.Service interface.
func (m *Monitor) Start(s service.Service) error {
	cloudflareDNS, err := cf.New()
	if err != nil {
		return fmt.Errorf("error creating Cloudflare API client: %w", err)
	}
	m.provider = cloudflareDNS

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
		m.providerUpdateZone()
		for {
			select {
			case <-ticker.C:
				m.providerUpdateZone()

			case <-sighup:
				log.Infof("SIGHUP received: updating records for %s", m.provider.GetZoneName())
				m.providerUpdateZone()
			}
		}
	}()
	<-stop
}

func (m *Monitor) providerUpdateZone() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.provider.UpdateZone(ctx); err != nil {
		log.Errorf("error updating zone %s: %v", m.provider.GetZoneName(), err)
	}
}

// Stop implements the service.Service interface.
func (m *Monitor) Stop(s service.Service) error {
	return nil
}
