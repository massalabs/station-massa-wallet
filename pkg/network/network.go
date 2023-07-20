package network

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/massalabs/station-massa-hello-world/pkg/plugin"
)

const (
	massaStationNodeEndpoint = plugin.MassaStationBaseURL + "/massa/node"
	defaultNetwork           = "testnet"
	defaultNodeUrl           = "https://test.massa.net/api/v2"
)

type NetworkInfo struct {
	Network string `json:"network"`
	URL     string `json:"url"`
	DNS     string `json:"dns"`
}

func logFallback(action string, err error) {
	log.Warnf("Fallback to default network: failed to %s: %v", action, err)
}

// retrieve network info from endpoint
func GetNetworkInfo() (*NetworkInfo, error) {
	resp, err := http.Get(massaStationNodeEndpoint)
	if err != nil {
		logFallback("GET massa station node endpoint", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logFallback("read response body", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	var data NetworkInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		logFallback("parse JSON", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	return &data, nil
}
