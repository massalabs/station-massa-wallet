package network

import (
	"log"
	"testing"

	"github.com/massalabs/station-massa-wallet/internal/initialize"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	err := initialize.Logger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	networkInfo, err := GetNetworkInfo()

	assert.NoError(t, err)

	assert.Equal(t, "buildnet", networkInfo.Network)
	assert.Equal(t, "https://buildnet.massa.net/api/v2", networkInfo.URL)
}
