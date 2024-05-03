package network

import (
	"log"
	"testing"

	testutils "github.com/massalabs/station-massa-wallet/pkg/testUtils"
	"github.com/stretchr/testify/assert"
)

func TestAssetExistInNetwork(t *testing.T) {
	err := testutils.LoggerTest()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	nodeFetcher := NewNodeFetcher()

	notExists := nodeFetcher.AssetExistInNetwork("nonexistentContractAddress")
	assert.False(t, notExists, "Expected asset not to exist in the network")
}
