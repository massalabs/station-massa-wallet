package network

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/massalabs/thyra/pkg/node"
)

func newMassaClient() (*node.Client, error) {

	networkInfo, err := GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	return node.NewClient(networkInfo.URL), nil
}

func (n *NodeFetcher) GetAccountsInfos(wlt []wallet.Wallet) ([]node.Address, error) {

	client, err := newMassaClient()
	if err != nil {
		return nil, err
	}

	addresses := make([]string, len(wlt))
	for i, addr := range wlt {
		addresses[i] = addr.Address
	}

	infos, err := node.Addresses(client, addresses)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
