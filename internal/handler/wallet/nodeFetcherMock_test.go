package wallet

import (
	"fmt"

	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/massalabs/thyra/pkg/node"
)

type NodeFetcherMock struct {
	client *node.Client
}

func NewNodeFetcherMock() *NodeFetcherMock {
	return &NodeFetcherMock{client: nil}
}

func (n *NodeFetcherMock) GetAccountsInfos(wlt []wallet.Wallet) ([]node.Address, error) {
	infos := make([]node.Address, len(wlt))
	for i, addr := range wlt {
		infos[i] = node.Address{
			Address:          addr.Address,
			CandidateBalance: fmt.Sprint(i + 1*1000000),
			FinalBalance:     fmt.Sprint(i + 1*1000000),
		}
	}
	return infos, nil
}

// Verifies at compilation time that NodeFetcherMock implements NodeFetcherInterface interface.
var _ network.NodeFetcherInterface = &NodeFetcherMock{}
