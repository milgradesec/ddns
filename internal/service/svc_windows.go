package service

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc/mgr"
)

const serviceName = "ddns"

// Install as a service
func Install() error {
	bin, err := os.Executable()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(serviceName)
	if err == nil {
		srv.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}

	srvConfig := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: serviceName,
		Description: "Dynamic DNS client",
	}

	srv, err = m.CreateService(serviceName, bin, srvConfig)
	if err != nil {
		return err
	}
	defer srv.Close()

	return nil
}

// Remove service
func Remove() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer srv.Close()

	err = srv.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Start service
func Start() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	err = srv.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

// Run the service
func Run(name string) error {
	//svc.Run(name, &service{})
	return nil
}

// Execute implements svc.Handler interface
/*func (m *Monitor) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go func() {
		m.Run()
	}()

	go func() {
		for {
			select {
			case c := <-r:
				switch c.Cmd {
				case svc.Interrogate:
					changes <- c.CurrentStatus
					time.Sleep(100 * time.Millisecond)
					changes <- c.CurrentStatus
				case svc.Stop, svc.Shutdown:
					break
				case svc.Pause:
					changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				case svc.Continue:
					changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				default:
				}
			}
		}
	}()

	changes <- svc.Status{State: svc.StopPending}
	return
}*/
