package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/monitor"
	"github.com/milgradesec/ddns/internal/updater"
)

var (
	// Version set at build
	Version string
)

func main() {
	fmt.Printf("DDNS-%s\n", Version)
	fmt.Printf("%s/%s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version(), Version)

	var (
		versionFlag = flag.Bool("version", false, "Only version information.")
		updateFlag  = flag.Bool("update", false, "Updates DDNS binary to latest version available.")
		serviceFlag = flag.String("service", "", "Manage DDNS as a system service")
		configFlag  = flag.String("config", "config.json", "Set configuration file.")
		helpFlag    = flag.Bool("help", false, "Show help.")
	)
	flag.Parse()

	if len(flag.Args()) > 0 {
		log.Fatalf("[ERROR] extra command line arguments.")
	}

	if *versionFlag {
		return
	}

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	if *updateFlag {
		err := updater.Update(Version)
		if err != nil {
			log.Fatalf("[ERROR] update failed: %v.", err)
		}
		return
	}

	m := &monitor.Monitor{
		Config: *configFlag,
	}

	svcConfig := &service.Config{
		Name:        "ddns",
		DisplayName: "Dynamic DNS service",
		Description: "Dynamic DNS service",
		Arguments:   []string{"-config", *configFlag},
	}

	svc, err := service.New(m, svcConfig)
	if err != nil {
		log.Fatalf("[ERROR] %v.", err)
	}

	if *serviceFlag != "" {
		if err := service.Control(svc, *serviceFlag); err != nil {
			log.Fatalf("[ERROR] %v", err)
		}

		switch *serviceFlag {
		case "install":
			log.Println("[INFO] service created successfully.")
		case "uninstall":
			log.Println("[INFO] service removed successfully.")
		case "start":
			log.Println("[INFO] service started.")
		case "stop":
			log.Println("[INFO] service stopped.")
		case "restart":
			log.Println("[INFO] service restarted.")
		default:
			log.Fatalf("[ERROR] invalid argument: %s.", *serviceFlag)
		}
		return
	}

	if err := svc.Run(); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}
