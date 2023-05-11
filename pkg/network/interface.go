package network

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
)

type NodeFetcher struct{}

func NewNodeFetcher() *NodeFetcher {
	return &NodeFetcher{}
}

type NodeFetcherInterface interface {
	GetAccountsInfos(wlt []wallet.Wallet) ([]node.Address, error)
	MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error)
	MakeRPCCall(msg []byte, signature []byte, publicKey string) ([]string, error)
}

// Verifies at compilation time that NodeFetcher implements NodeFetcherInterface
var _ NodeFetcherInterface = &NodeFetcher{}
