package network

import "github.com/massalabs/station/pkg/node"

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

	return uint64(*chainID), nil
}
