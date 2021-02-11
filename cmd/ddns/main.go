package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/monitor"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		versionFlag = flag.Bool("version", false, "Show version information.")
		serviceFlag = flag.String("service", "", "Manage DDNS as a system service")
		configFlag  = flag.String("config", "config.json", "Set configuration file.")
		helpFlag    = flag.Bool("help", false, "Show help.")
	)
	flag.Parse()

	if *versionFlag {
		fmt.Println("DDNS-" + Version)
		fmt.Printf("%s/%s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version(), Version)
		return
	}

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	m := &monitor.Monitor{
		ConfigFile: *configFlag,
	}

	svcConfig := &service.Config{
		Name:        "ddns",
		DisplayName: "ddns",
		Description: "Dynamic DNS service",
		Arguments:   []string{"-config", *configFlag},
	}

	svc, err := service.New(m, svcConfig)
	if err != nil {
		log.Fatalf("%v.", err)
	}

	if *serviceFlag != "" {
		if err := service.Control(svc, *serviceFlag); err != nil {
			log.Fatalf("%v", err)
		}

		switch *serviceFlag {
		case "install":
			log.Infoln("service created successfully")
		case "uninstall":
			log.Infoln("service removed successfully")
		case "start":
			log.Infoln("service started")
		case "stop":
			log.Infoln("service stopped")
		case "restart":
			log.Infoln("service restarted")
		default:
			log.Errorf("invalid argument: %s", *serviceFlag)
		}
		return
	}

	if err := svc.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

var (
	Version string
)
