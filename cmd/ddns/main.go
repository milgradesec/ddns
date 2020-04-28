package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/monitor"
	cf "github.com/milgradesec/ddns/internal/provider/cloudflare"
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
	if len(flag.Args()) > 0 {
		fmt.Printf("extra command line arguments: %s", flag.Args())
		os.Exit(1)
	}

	if *version {
		os.Exit(0)
	}

	/*_, set := os.LookupEnv("PROVIDER")
	if set == true {
		// load from env
	}*/

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("error loading config: %v", err)
		os.Exit(1)
	}

	p, err := cf.New(cfg)
	if err != nil {
		fmt.Printf("Cloudflare API login failed: %v\n", err)
		os.Exit(1)
	}

	monitor := monitor.New(os.Getenv("CF_ZONE_NAME"), p)
	monitor.Run()
}
