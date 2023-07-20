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
	defaultNetwork           = "buildnet"
	defaultNodeUrl           = "https://buildnet.massa.net/api/v2"
	fallbackMsg              = "Fallback to default network. "
)

type NetworkInfo struct {
	Network string `json:"network"`
	URL     string `json:"url"`
	DNS     string `json:"dns"`
}

// retrieve network info from endpoint
func GetNetworkInfo() (*NetworkInfo, error) {
	resp, err := http.Get(massaStationNodeEndpoint)
	if err != nil {
		log.Warnf(fallbackMsg+"Failed to make GET request: %v", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf(fallbackMsg+"Failed to read response body: %v", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	var data NetworkInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Warnf(fallbackMsg+"Failed to parse JSON: %v", err)
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	return &data, nil
}
