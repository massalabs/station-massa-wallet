package network

import (
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

// MakeRPCCall makes a new RPC call by calling Massa Station source code sendOperation.MakeRPCCall function.
func (n *NodeFetcher) MakeRPCCall(msg []byte, signature []byte, publicKey string) ([]string, error) {
	client, err := NewMassaClient()
	if err != nil {
		return nil, err
	}

	return sendOperation.MakeRPCCall(msg, signature, publicKey, client)
}
