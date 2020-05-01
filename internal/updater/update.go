package updater

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/inconshreveable/go-update"
)

const baseURL = "https://dl.paesacybersecurity.eu/ddns/"

type updateInfo struct {
	Version string `json:"version"`
	Sha256  string `json:"sha256"`
}

func checkForUpdateAndApply(version string) error {
	resp, err := http.Get(baseURL + runtime.GOOS + "-" + runtime.GOARCH + ".json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid response from server")
	}

	var info updateInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return err
	}

	if info.Version == "" {
		return errors.New("invalid response from server")
	}

	if info.Version == version {
		log.Println("DDNS is up to date.")
		return nil
	}
	log.Printf("New version %s is available.\n", info.Version)

	resp, err = http.Get(baseURL + info.Version + "/" + runtime.GOOS + "-" + runtime.GOARCH)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid response from server")
	}

	checksum, err := hex.DecodeString(info.Sha256)
	if err != nil {
		return err
	}

	opts := update.Options{
		Checksum: []byte(checksum),
		Hash:     crypto.SHA256,
	}
	err = update.Apply(resp.Body, opts)
	if err != nil {
		return err
	}
	log.Printf("DDNS updated to version %s\n", info.Version)

	return nil
}

// Update updates DDNS binary to the latest version available.
func Update(version string) error {
	return checkForUpdateAndApply(version)
}
