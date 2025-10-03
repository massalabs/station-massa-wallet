package network

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/massalabs/station/pkg/logger"
	pluginKit "github.com/massalabs/station/plugin-kit"
)

const (
	massaStationNodeEndpoint = pluginKit.MassaStationBaseURL + "/massa/node"
	defaultNetwork           = "buildnet"
	defaultNodeUrl           = "https://buildnet.massa.net/api/v2"
	defaultChainId           = 77658366
)

type NetworkInfo struct {
	Network string `json:"network"`
	URL     string `json:"url"`
	ChainID int    `json:"chainId"`
}

func logFallback(action string, err error) {
	logger.Warnf("Fallback to default network: failed to %s: %v", action, err)
}

// retrieve network info from endpoint
func GetNetworkInfo() (*NetworkInfo, error) {
	resp, err := http.Get(massaStationNodeEndpoint)
	if err != nil {
		logFallback("GET massa station node endpoint", err)

		return fallbackNetworkInfo(), nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logFallback("read response body", err)
		return fallbackNetworkInfo(), nil
	}

	var data NetworkInfo

	err = json.Unmarshal(body, &data)
	if err != nil {
		logFallback("parse JSON", err)
		return fallbackNetworkInfo(), nil
	}

	return &data, nil
}

func fallbackNetworkInfo() *NetworkInfo {
	return &NetworkInfo{
		Network: defaultNetwork,
		URL:     defaultNodeUrl,
		ChainID: defaultChainId,
	}
}
