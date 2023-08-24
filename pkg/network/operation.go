package network

import (
	"fmt"

	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

// MakeOperation makes a new operation by calling Massa Station source code sendOperation.MakeOperation function.
func (n *NodeFetcher) MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error) {
	client, err := NewMassaClient()
	if err != nil {
		return nil, err
	}

	msg, _, err := sendOperation.MakeOperation(client, sendOperation.DefaultSlotsDuration, fee, operation)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func SendOperation(wlt *wallet.Wallet, massaClient NodeFetcherInterface, operation sendOperation.Operation, fee uint64) (*sendOperation.OperationResponse, *wallet.WalletError) {
	operationData, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, &wallet.WalletError{Err: fmt.Errorf("Error while making operation: %w", err), CodeErr: utils.ErrNetwork}
	}

	// TODO: we do not implement the handling of the correlation id for now
	signature, err := wlt.Sign(true, operationData)
	if err != nil {
		return nil, &wallet.WalletError{Err: fmt.Errorf("Error sign: %w", err), CodeErr: utils.ErrUnknown}
	}

	// send the operationData to the network
	resp, err := massaClient.MakeRPCCall(operationData, signature, wlt.GetPupKey())
	if err != nil {
		// unknown error because it could be the signature, the network, the node...
		return nil, &wallet.WalletError{Err: fmt.Errorf("Error during RPC call: %w", err), CodeErr: utils.ErrUnknown}
	}

	return &sendOperation.OperationResponse{CorrelationID: "", OperationID: resp[0]}, nil
}
