package network

import (
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
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
	password *memguard.LockedBuffer,
	massaClient NodeFetcherInterface,
	operation sendOperation.Operation,
	fee uint64,
	chainID uint64,
) (*sendOperation.OperationResponse, error) {
	operationData, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, fmt.Errorf("Error while making operation: %w", err)
	}

	publicKey, err := acc.PublicKey.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %w", err)
	}

	operationDataToSign := utils.PrepareSignData(uint64(chainID), append(publicKey, operationData...))

	signature, err := acc.Sign(password, operationDataToSign)
	if err != nil {
		return nil, fmt.Errorf("Error while signing operation: %w", err)
	}

	publicKeyText, err := acc.PublicKey.MarshalText()
	if err != nil {
		return nil, err
	}

	// send the operationData to the network
	resp, err := massaClient.MakeRPCCall(operationData, signature, string(publicKeyText))
	if err != nil {
		// unknown error because it could be the signature, the network, the node...
		return nil, fmt.Errorf("Error during RPC call: %w", err)
	}

	return &sendOperation.OperationResponse{OperationID: resp[0]}, nil
}
