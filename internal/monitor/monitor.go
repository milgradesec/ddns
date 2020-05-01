package monitor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/provider"
	cf "github.com/milgradesec/ddns/internal/provider/cloudflare"
)

const defaultInterval = 3 * time.Minute

// Monitor runs in a infinite loop and triggers provider zone updates
// every 3 min interval, can be triggered at any time by sending a
// SIGHUP signal
type Monitor struct {
	Config string

	cfg config.Config
	api provider.API
}

// Start implements service.Interface
func (m *Monitor) Start(s service.Service) error {
	cfg, err := config.Load(m.Config)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	cfAPI, err := cf.New(cfg)
	if err != nil {
		return fmt.Errorf("error creating Cloudflare API client: %v", err)
	}
	m.api = cfAPI
	go m.Run()
	return nil
}

// Stop implements service.Interface
func (m *Monitor) Stop(s service.Service) error {
	return nil
}

// Run starts monitoring
func (m *Monitor) Run() {
	ticker := time.NewTicker(defaultInterval)
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
				log.Printf("[INFO] SIGHUP received: updating records for %s\n", m.cfg.Zone)
				m.callProvider()
			}
		}
	}()
	<-stop
}

func (m *Monitor) callProvider() {
	if err := m.api.UpdateZone(); err != nil {
		log.Printf("[ERROR] error updating zone %s: %v\n", m.cfg.Zone, err)
	}
}
