package network

import (
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

type NodeFetcher struct{}

func NewNodeFetcher() *NodeFetcher {
	return &NodeFetcher{}
}

type NodeFetcherInterface interface {
	GetAccountsInfos(wlt []wallet.Wallet) ([]AccountInfos, error)
	MakeOperation(fee uint64, operation sendOperation.Operation) ([]byte, error)
	MakeRPCCall(msg []byte, signature []byte, publicKey string) ([]string, error)
	AssetExistInNetwork(contractAddress string) bool
	DatastoreAssetName(contractAddress string) (string, error)
	DatastoreAssetSymbol(contractAddress string) (string, error)
	DatastoreAssetDecimals(contractAddress string) (uint8, error)
	DatastoreAssetBalance(assetContractAddress, userAddress string) (string, error)
}

// Verifies at compilation time that NodeFetcher implements NodeFetcherInterface
var _ NodeFetcherInterface = &NodeFetcher{}
