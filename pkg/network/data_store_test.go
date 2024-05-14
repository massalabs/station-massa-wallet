package network

import (
	"log"
	"testing"

	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestAssetExistInNetwork(t *testing.T) {
	if err := logger.InitializeGlobal("./unit-test.log"); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	nodeFetcher := NewNodeFetcher()

	notExists := nodeFetcher.AssetExistInNetwork("nonexistentContractAddress")
	assert.False(t, notExists, "Expected asset not to exist in the network")
}
