package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

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
		version    = flag.Bool("version", false, "Only version information.")
		update     = flag.Bool("update", false, "Updates DDNS binary to latest version available.")
		configFile = flag.String("config", "config.json", "Set configuration file.")
		help       = flag.Bool("help", false, "Show help.")
	)
	flag.Parse()

	if *version {
		os.Exit(0)
	}

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *update {
		err := updater.Update(Version)
		if err != nil {
			log.Fatalf("update error: %v", err)
		}
		os.Exit(0)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	monitor, err := monitor.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	monitor.Run()
}
