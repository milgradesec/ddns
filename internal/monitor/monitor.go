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
	config config.Config
	prov   provider.API
}

// New creates a Monitor with the provided configuration
func New(cfg config.Config) (*Monitor, error) {
	p, err := cf.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Cloudflare API client: %v", err)
	}

	return &Monitor{
		config: cfg,
		prov:   p,
	}, nil
}

// Start implements service.Interface
func (m *Monitor) Start(s service.Service) error {
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
				log.Printf("[INFO] SIGHUP received: updating records for %s\n", m.config.Zone)
				m.callProvider()
			}
		}
	}()
	<-stop
}

func (m *Monitor) callProvider() {
	if err := m.prov.UpdateZone(); err != nil {
		log.Printf("[ERROR] error updating zone %s: %v\n", m.config.Zone, err)
	}
}
