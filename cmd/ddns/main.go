package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/monitor"
)

var (
	// Version set at build
	Version string
)

func main() {
	fmt.Printf("DDNS-%s\n", Version)
	fmt.Printf("%s/%s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version(), Version)

	var (
		version    = flag.Bool("version", false, "Show version")
		configFile = flag.String("config", "config.json", "Configuration file")
	)
	flag.Parse()

	if *version {
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
