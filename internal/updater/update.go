package updater

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/inconshreveable/go-update"
	log "github.com/sirupsen/logrus"
)

const baseURL = "https://dl.paesacybersecurity.eu/ddns/"

var ErrNotAvailable = errors.New("no update available")

func Update(version string) error {
	resp, err := checkForUpdate(version)
	if err != nil {
		if errors.Is(err, ErrNotAvailable) {
			log.Println("no updates available")
			return nil
		}
		return err
	}

	log.Infof("new version %s is available", resp.ReleaseVersion)
	err = resp.Apply()
	if err != nil {
		return err
	}

	log.Infof("updated to version %s", resp.ReleaseVersion)
	return nil
}

type Response struct {
	ReleaseVersion string
	checksum       []byte
}

type serverResponse struct {
	Version string `json:"version"`
	Sha256  string `json:"sha256"`
}

func checkForUpdate(version string) (r Response, err error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(baseURL + runtime.GOOS + "-" + runtime.GOARCH + ".json")
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("server responded with status code %s", resp.Status)
	}

	var serverResp serverResponse
	err = json.NewDecoder(resp.Body).Decode(&serverResp)
	if err != nil {
		return r, err
	}

	if serverResp.Version == version {
		return r, ErrNotAvailable
	}

	r.ReleaseVersion = serverResp.Version
	r.checksum, err = hex.DecodeString(serverResp.Sha256)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (r Response) Apply() error {
	opts := update.Options{
		Checksum: r.checksum,
		Hash:     crypto.SHA256,
	}

	resp, err := r.fetchUpdate()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r Response) fetchUpdate() (*http.Response, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(baseURL + r.ReleaseVersion + "/" + runtime.GOOS + "-" + runtime.GOARCH)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New("server responded with code: %s" + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}
