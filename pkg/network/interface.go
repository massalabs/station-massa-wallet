package network

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

type NodeFetcher struct{}

func NewNodeFetcher() *NodeFetcher {
	return &NodeFetcher{}
}

type NodeFetcherInterface interface {
	GetAccountsInfos(wlt []wallet.Wallet) ([]AccountInfos, error)
}

// Verifies at compilation time that NodeFetcher implements NodeFetcherInterface
var _ NodeFetcherInterface = &NodeFetcher{}
