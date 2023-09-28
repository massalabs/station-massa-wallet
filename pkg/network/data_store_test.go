package network

import (
	"log"
	"testing"

	"github.com/massalabs/station-massa-wallet/internal/initialize"
	"github.com/stretchr/testify/assert"
)

func TestAssetExistInNetwork(t *testing.T) {
	err := initialize.Logger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	nodeFetcher := NewNodeFetcher()

	notExists := nodeFetcher.AssetExistInNetwork("nonexistentContractAddress")
	assert.False(t, notExists, "Expected asset not to exist in the network")
}
