package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/config"
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
		os.Exit(0)
	}

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *updateFlag {
		err := updater.Update(Version)
		if err != nil {
			log.Fatalf("[ERROR] update failed: %v.", err)
		}
		os.Exit(0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to find the path to the current directory")
	}

	cfg, err := config.Load(*configFlag)
	if err != nil {
		log.Fatalf("[ERROR] failed to load configuration: %v.", err)
	}

	m, err := monitor.New(cfg)
	if err != nil {
		log.Fatalf("[ERROR] %v.", err)
	}

	svcConfig := &service.Config{
		Name:             "ddns",
		Description:      "Dynamic DNS client",
		WorkingDirectory: cwd,
		Arguments:        []string{"-config", *configFlag},
	}
	svc, err := service.New(m, svcConfig)
	if err != nil {
		log.Fatalf("[ERROR] %v.", err)
	}

	if *serviceFlag != "" {
		if err := service.Control(svc, *serviceFlag); err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
	}

	if err := svc.Run(); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	/*if *serviceFlag != "" {
		switch *serviceFlag {
		case "install":
			if err := service.Install(); err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			log.Println("[INFO] service created successfully.")
			return

		case "uninstall":
			if err := service.Uninstall(); err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			log.Println("[INFO] service removed successfully.")
			return

		case "start":
			if err := service.Start(); err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			log.Println("[INFO] service started.")
			return

		case "stop":
			if err := service.Control(svc.Stop, svc.Stopped); err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			log.Println("[INFO] service stopped")
			return
		default:
			log.Fatalf("[ERROR] invalid argument: \"%s\".", *serviceFlag)
		}
	}*/

	//m.Run()

}
