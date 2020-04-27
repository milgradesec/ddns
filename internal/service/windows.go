// Package service registers ddns to run as a Windows Service
// +build windows
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
		DisplayName: serviceName,
	}

	srv, err = m.CreateService(serviceName, bin, srvConfig, "is", "auto-started")
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
