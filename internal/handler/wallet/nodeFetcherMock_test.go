package wallet

import (
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
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

// MakeOperation returns a dummy operation
func (n *NodeFetcherMock) MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error) {
	msg := []byte{0, 226, 204, 2, 0, 0, 84, 241, 88, 133, 82, 70, 100, 219, 180, 87, 210, 99, 186, 197, 218, 51, 252, 165, 147, 138, 98, 206, 27, 228, 157, 142, 104, 250, 30, 86, 49, 22, 1}
	return msg, nil
}

// MakeRPCCall returns a dummy RPC call response
func (n *NodeFetcherMock) MakeRPCCall(msg []byte, signature []byte, publicKey string) ([]string, error) {
	return []string{"[O1Mw8wdurZphk6VbfB7i4irGwXmGkbpRmrLR84xfP5Ui4qEBy4z]"}, nil
}

// Verifies at compilation time that NodeFetcherMock implements NodeFetcherInterface interface.
var _ network.NodeFetcherInterface = &NodeFetcherMock{}
