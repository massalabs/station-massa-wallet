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

	// Call the function to get network info
	networkInfo, err := GetNetworkInfo()

	// Ensure no error occurred
	assert.NoError(t, err)

	// Assert that the returned networkInfo matches the expected values
	assert.Equal(t, "buildnet", networkInfo.Network)
	assert.Equal(t, "https://buildnet.massa.net/api/v2", networkInfo.URL)
}
