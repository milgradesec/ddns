package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/monitor"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

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
		log.Fatal().Msgf("%v.", err)
	}

	if *serviceFlag != "" {
		if err := service.Control(svc, *serviceFlag); err != nil {
			log.Fatal().Msgf("%v", err)
		}

		switch *serviceFlag {
		case "install":
			log.Info().Msg("service created successfully")
		case "uninstall":
			log.Info().Msg("service removed successfully")
		case "start":
			log.Info().Msg("service started")
		case "stop":
			log.Info().Msg("service stopped")
		case "restart":
			log.Info().Msg("service restarted")
		default:
			log.Error().Msgf("invalid argument: %s", *serviceFlag)
		}
		return
	}
	if err := svc.Run(); err != nil {
		log.Fatal().Msgf("%v", err)
	}
}

var (
	Version string
)
