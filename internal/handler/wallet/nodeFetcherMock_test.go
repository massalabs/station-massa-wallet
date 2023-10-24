package wallet

import (
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

type NodeFetcherMock struct {
	client *node.Client
}

func NewNodeFetcherMock() *NodeFetcherMock {
	return &NodeFetcherMock{client: nil}
}

// returns dummy balances
func (n *NodeFetcherMock) GetAccountsInfos(accounts []*account.Account) ([]network.AccountInfos, error) {
	infos := make([]network.AccountInfos, len(accounts))

	for i, acc := range accounts {
		textAddress, err := acc.Address.MarshalText()
		if err != nil {
			return nil, err
		}
		infos[i] = network.AccountInfos{
			Address:          string(textAddress),
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

// DatastoreAssetName returns a dummy Asset Name.
func (n *NodeFetcherMock) DatastoreAssetName(contractAddress string) (string, error) {
	return "TestToken", nil
}

// DatastoreAssetSymbol  returns a dummy Asset Symbol.
func (n *NodeFetcherMock) DatastoreAssetSymbol(contractAddress string) (string, error) {
	return "TST", nil
}

// DatastoreAssetDecimals  returns a dummy Asset Decimals.
func (n *NodeFetcherMock) DatastoreAssetDecimals(contractAddress string) (uint8, error) {
	return 9, nil
}

// DatastoreAssetBalanceMock returns a dummy Asset Balance.
func (n *NodeFetcherMock) DatastoreAssetBalance(assetContractAddress, userAddress string) (string, error) {
	// Return a balance value of 10000 as a string
	return "10000", nil
}

func (n *NodeFetcherMock) AssetExistInNetwork(contractAddress string) bool {
	return contractAddress != "AS12GwD3UEk2BP1zMx2zSdvKov97z8gs1MtsoN4u4C9emLBbhYa3U"
}

// Verifies at compilation time that NodeFetcherMock implements NodeFetcherInterface interface.
var _ network.NodeFetcherInterface = &NodeFetcherMock{}
