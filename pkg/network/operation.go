package network

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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

func SendOperation(wlt *wallet.Wallet, massaClient NodeFetcherInterface, operation sendOperation.Operation, fee uint64) (*sendOperation.OperationResponse, error) {
	msg, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, fmt.Errorf("Error while making operation: %w", err)
	}

	// sign the msg in base64
	// TODO: we do not implement the handling of the correlation id for now
	byteMsgB64 := strfmt.Base64(msg)
	signature, err := wlt.Sign(&byteMsgB64)
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
