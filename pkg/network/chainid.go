package network

import (
	"errors"

	"github.com/massalabs/station/pkg/node"
)

var ErrChainIDNotInStatus = errors.New("chain id not in status")

func GetNodeInfo() (uint64, string, error) {
	client, err := NewMassaClient()
	if err != nil {
		return 0, "", err
	}

	status, err := node.Status(client)
	if err != nil {
		return 0, "", err
	}

	chainID := status.ChainID

	if chainID == nil {
		return 0, "", ErrChainIDNotInStatus
	}

	var minimalFees string
	if status.MinimalFees == nil {
		minimalFees = "0"
	} else {
		minimalFees = *status.MinimalFees
	}

	return uint64(*chainID), minimalFees, nil
}
}
