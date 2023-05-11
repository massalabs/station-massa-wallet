package network

import (
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
)

// MakeRPCCall makes a new RPC call by calling Thyra source code sendOperation.MakeRPCCall function.
func (n *NodeFetcher) MakeRPCCall(msg []byte, signature []byte, publicKey string) ([]string, error) {
	client, err := newMassaClient()
	if err != nil {
		return nil, err
	}

	return sendOperation.MakeRPCCall(msg, signature, publicKey, client)
}
