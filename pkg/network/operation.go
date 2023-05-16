package network

import (
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
)

// MakeOperation makes a new operation by calling Thyra source code sendOperation.MakeOperation function.
func (n *NodeFetcher) MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error) {
	client, err := newMassaClient()
	if err != nil {
		return nil, err
	}

	msg, _, err := sendOperation.MakeOperation(client, sendOperation.DefaultSlotsDuration, fee, operation)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
