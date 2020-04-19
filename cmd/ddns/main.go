package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/milgradesec/ddns/monitor"
	cf "github.com/milgradesec/ddns/provider/cloudflare"
)

var (
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
		log.Fatalf("extra command line arguments: %s", flag.Args())
	}

	provider, err := cf.New()
	if err != nil {
		log.Fatal(err)
	}

	monitor := monitor.New(os.Getenv("CF_ZONE_NAME"), provider)
	monitor.Run()
}

var (
	version bool
)
