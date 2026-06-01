package network

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	if err := logger.InitializeGlobal("./unit-test.log"); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(NetworkInfo{
			Network: "buildnet",
			URL:     "https://buildnet.massa.net/api/v2",
			ChainID: defaultChainId,
		})
	}))
	t.Cleanup(server.Close)

	networkInfo, err := fetchNetworkInfo(server.URL)

	assert.NoError(t, err)
	assert.Equal(t, "buildnet", networkInfo.Network)
	assert.Equal(t, "https://buildnet.massa.net/api/v2", networkInfo.URL)
	assert.Equal(t, defaultChainId, networkInfo.ChainID)
}

func TestFallbackNetworkInfo(t *testing.T) {
	networkInfo := fallbackNetworkInfo()

	assert.Equal(t, "buildnet", networkInfo.Network)
	assert.Equal(t, "https://buildnet.massa.net/api/v2", networkInfo.URL)
	assert.Equal(t, defaultChainId, networkInfo.ChainID)
}
