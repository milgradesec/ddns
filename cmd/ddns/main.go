package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/milgradesec/ddns/internal/monitor"
	cf "github.com/milgradesec/ddns/internal/provider/cloudflare"
)

var (
	// Version set from Makefile
	Version string
)

func init() {
	flag.BoolVar(&version, "version", false, "Show version")
}

func main() {
	fmt.Printf("DDNS-%s\n", Version)
	fmt.Printf("%s/%s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version(), Version)

	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Printf("extra command line arguments: %s", flag.Args())
		os.Exit(1)
	}

	p, err := cf.New()
	if err != nil {
		fmt.Printf("Cloudflare API login failed: %v\n", err)
		os.Exit(1)
	}

	monitor := monitor.New(os.Getenv("CF_ZONE_NAME"), p)
	monitor.Run()
}

var (
	version bool
)
