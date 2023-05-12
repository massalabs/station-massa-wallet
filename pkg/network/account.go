package network

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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

type AccountInfos struct {
	Address          string
	CandidateBalance uint64
	Balance          uint64
}

func (n *NodeFetcher) GetAccountsInfos(wlt []wallet.Wallet) ([]AccountInfos, error) {
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

	res := make([]AccountInfos, len(infos))
	for i, info := range infos {
		res[i].Address = info.Address
		nano, err := utils.MasToNano(info.CandidateBalance)
		if err != nil {
			return nil, err
		}
		res[i].CandidateBalance = nano
		nano, err = utils.MasToNano(info.FinalBalance)
		if err != nil {
			return nil, err
		}
		res[i].Balance = nano
	}

	return res, nil
}
