package network

import (
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
)

func NewMassaClient() (*node.Client, error) {
	networkInfo, err := GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Connected to node URL: %s, chain id: %d", networkInfo.URL, networkInfo.ChainID)

	return node.NewClient(networkInfo.URL), nil
}

type AccountInfos struct {
	Address          string
	CandidateBalance uint64
	Balance          uint64
}

func (n *NodeFetcher) GetAccountsInfos(accounts []*account.Account) ([]AccountInfos, error) {
	client, err := NewMassaClient()
	if err != nil {
		return nil, err
	}

	addresses := make([]string, len(accounts))

	for i, acc := range accounts {
		textAddress, err := acc.Address.MarshalText()
		if err != nil {
			return nil, err
		}
		addresses[i] = string(textAddress)
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
