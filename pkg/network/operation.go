package network

import (
	"fmt"

	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
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

func SendOperation(wlt *wallet.Wallet, massaClient NodeFetcherInterface, operation sendOperation.Operation, fee uint64) (*sendOperation.OperationResponse, error) {
	msg, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, fmt.Errorf("Error while making operation: %w", err)
	}

	// TODO: we do not implement the handling of the correlation id for now
	signature, err := wlt.Sign(msg)
	if err != nil {
		return nil, fmt.Errorf("Error sign: %w", err)
	}

	// send the msg to the network
	resp, err := massaClient.MakeRPCCall(msg, signature, wlt.GetPupKey())
	if err != nil {
		return nil, fmt.Errorf("Error during RPC call: %w", err)
	}

	return &sendOperation.OperationResponse{CorrelationID: "", OperationID: resp[0]}, nil
}
