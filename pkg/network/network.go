package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/massalabs/station-massa-hello-world/pkg/plugin"
)

const (
	massaStationNodeEndpoint = plugin.MassaStationBaseURL + "/massa/node"
	defaultNetwork           = "buildnet"
	defaultNodeUrl           = "https://buildnet.massa.net/api/v2"
)

type NetworkInfo struct {
	Network string `json:"network"`
	URL     string `json:"url"`
	DNS     string `json:"dns"`
}

// retrieve network info from endpoint
func GetNetworkInfo() (*NetworkInfo, error) {
	if os.Getenv("STANDALONE") == "1" {
		return &NetworkInfo{Network: defaultNetwork, URL: defaultNodeUrl}, nil
	}

	resp, err := http.Get(massaStationNodeEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data NetworkInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &data, nil
}
