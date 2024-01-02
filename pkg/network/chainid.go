package network

import (
	"errors"

	"github.com/massalabs/station/pkg/node"
)

var ErrChainIDNotInStatus = errors.New("chain id not in status")

func getChainID() (uint64, error) {
	client, err := NewMassaClient()
	if err != nil {
		return 0, err
	}

	status, err := node.Status(client)
	if err != nil {
		return 0, err
	}

	chainID := status.ChainID

	if chainID == nil {
		return 0, ErrChainIDNotInStatus
	}

	return uint64(*chainID), nil
}
