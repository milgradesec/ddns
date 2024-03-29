package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/kardianos/service"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/monitor"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var (
		versionFlag = flag.Bool("version", false, "Show version information.")
		serviceFlag = flag.String("service", "", "Manage DDNS as a system service")
		debug       = flag.Bool("debug", false, "Enable debug logging.")
		helpFlag    = flag.Bool("help", false, "Show help.")
	)
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msgf("DDNS %s", Version)
	log.Info().Msgf("%s/%s, %s", runtime.GOOS, runtime.GOARCH, runtime.Version())

	if *versionFlag {
		return
	}

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	svcConfig := &service.Config{
		Name:        "ddns",
		DisplayName: "ddns",
		Description: "Dynamic DNS service",
	}

	svc, err := service.New(monitor.New(config), svcConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create service")
	}

	if *serviceFlag != "" {
		if err := service.Control(svc, *serviceFlag); err != nil {
			log.Fatal().Err(err).Msg("service control error")
		}

		switch *serviceFlag {
		case "install":
			log.Info().Msg("Service installed successfully.")
		case "uninstall":
			log.Info().Msg("Service removed successfully.")
		case "start":
			log.Info().Msg("Service started.")
		case "stop":
			log.Info().Msg("Service stopped.")
		case "restart":
			log.Info().Msg("Service restarted.")
		default:
			log.Error().Msgf("invalid argument: %s", *serviceFlag)
		}
		return
	}
	if err := svc.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start service")
	}
}

var (
	Version = "DEV"
)
