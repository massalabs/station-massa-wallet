package network

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/massalabs/thyra/pkg/node"
)

type NodeFetcher struct {
}

func NewNodeFetcher() *NodeFetcher {
	return &NodeFetcher{}
}

type NodeFetcherInterface interface {
	GetAccountsInfos(wlt []wallet.Wallet) ([]node.Address, error)
}

// Verifies at compilation time that NodeFetcher implements NodeFetcherInterface
var _ NodeFetcherInterface = &NodeFetcher{}
