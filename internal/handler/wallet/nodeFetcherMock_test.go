package wallet

import (
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

// returns dummy balances
func (n *NodeFetcherMock) GetAccountsInfos(wlt []wallet.Wallet) ([]network.AccountInfos, error) {
	infos := make([]network.AccountInfos, len(wlt))
	for i, addr := range wlt {
		infos[i] = network.AccountInfos{
			Address:          addr.Address,
			CandidateBalance: uint64(i + 1*1000000),
			Balance:          uint64(i + 1*1000000),
		}
	}
	return infos, nil
}

// Verifies at compilation time that NodeFetcherMock implements NodeFetcherInterface interface.
var _ network.NodeFetcherInterface = &NodeFetcherMock{}
