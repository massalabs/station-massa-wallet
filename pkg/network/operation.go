package network

import (
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

// MakeOperation makes a new operation by calling Massa Station source code sendOperation.MakeOperation function.
func (n *NodeFetcher) MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error) {
	client, err := NewMassaClient()
	if err != nil {
		return nil, err
	}

	msg, _, err := sendOperation.MakeOperation(client, sendOperation.DefaultExpiryInSlot, fee, operation)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func SendOperation(
	acc *account.Account,
	guardedPassword *memguard.LockedBuffer,
	massaClient NodeFetcherInterface,
	operation sendOperation.Operation,
	fee uint64,
) (*sendOperation.OperationResponse, *walletmanager.WalletError) {
	operationData, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, &walletmanager.WalletError{Err: fmt.Errorf("Error while making operation: %w", err), CodeErr: utils.ErrNetwork}
	}

	// TODO: we do not implement the handling of the correlation id for now
	signature, err := acc.Sign(guardedPassword, operationData)
	if err != nil {
		return nil, &walletmanager.WalletError{Err: fmt.Errorf("Error while signing operation: %w", err), CodeErr: utils.ErrUnknown}
	}

	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return nil, &walletmanager.WalletError{Err: fmt.Errorf("Unable to marshal public key: %s", err.Error()), CodeErr: utils.ErrUnknown}
	}

	// send the operationData to the network
	resp, err := massaClient.MakeRPCCall(operationData, signature, string(publicKeyBytes))
	if err != nil {
		// unknown error because it could be the signature, the network, the node...
		return nil, &walletmanager.WalletError{Err: fmt.Errorf("Error during RPC call: %w", err), CodeErr: utils.ErrUnknown}
	}

	return &sendOperation.OperationResponse{CorrelationID: "", OperationID: resp[0]}, nil
}
