package network

import (
	"log"
	"testing"

	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	if err := logger.InitializeGlobal("./unit-test.log"); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	networkInfo, err := GetNetworkInfo()

	assert.NoError(t, err)

	assert.Equal(t, "buildnet", networkInfo.Network)
	assert.Equal(t, "https://buildnet.massa.net/api/v2", networkInfo.URL)
}
