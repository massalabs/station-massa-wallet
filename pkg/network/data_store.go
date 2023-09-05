package network

import (
	"fmt"

	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
)

const (
	NAME_KEY           = "NAME"
	SYMBOL_KEY         = "SYMBOL"
	DECIMALS_KEY       = "DECIMALS"
	BALANCE_KEY_PREFIX = "BALANCE"
)

func (n *NodeFetcher) AssetExistInNetwork(contractAddress string) bool {
	// Call DatastoreEntry to retrieve the asset name data for the given contractAddress
	nameData, err := DatastoreEntry(contractAddress, convert.ToBytes(SYMBOL_KEY))
	if err != nil {
		// If there is an error, we assume the asset does not exist in the network
		return false
	}
	// Check if the length of nameData is greater than 0 to determine if the asset exists
	return len(nameData) > 0
}

// DatastoreAssetName retrieves the asset name for a given contract address from the Massa node.
func (n *NodeFetcher) DatastoreAssetName(contractAddress string) (string, error) {
	nameData, err := DatastoreEntry(contractAddress, convert.ToBytes(NAME_KEY))
	if err != nil {
		return "", fmt.Errorf("failed to fetch asset name: %w", err)
	}
	return string(nameData), nil
}

// DatastoreAssetSymbol retrieves the asset symbol for a given contract address from the Massa node.
func (n *NodeFetcher) DatastoreAssetSymbol(contractAddress string) (string, error) {
	symbolData, err := DatastoreEntry(contractAddress, convert.ToBytes(SYMBOL_KEY))
	if err != nil {
		return "", fmt.Errorf("failed to fetch asset symbol: %w", err)
	}
	return string(symbolData), nil
}

// DatastoreAssetDecimals retrieves the asset decimals for a given contract address from the Massa node.
func (n *NodeFetcher) DatastoreAssetDecimals(contractAddress string) (uint8, error) {
	decimalsData, err := DatastoreEntry(contractAddress, convert.ToBytes(DECIMALS_KEY))
	if err != nil {
		return 0, fmt.Errorf("failed to fetch asset decimals: %w", err)
	}
	return uint8(decimalsData[0]), nil
}

// DatastoreAssetBalance retrieves the balance of a user for a given asset contract address from the Massa node.
func (n *NodeFetcher) DatastoreAssetBalance(assetContractAddress, userAddress string) (string, error) {
	balanceData, err := DatastoreEntry(assetContractAddress, balanceKey(userAddress))
	if err != nil {
		return "", fmt.Errorf("failed to fetch user balance: %w", err)
	}

	balanceValue, err := convert.BytesToU256(balanceData)
	if err != nil {
		return "", fmt.Errorf("failed to parse user balance: %w", err)
	}

	return fmt.Sprint(balanceValue.String()), nil
}

// Function to convert an address to a storage key using the balance key prefix
func balanceKey(address string) []byte {
	return convert.ToBytes(BALANCE_KEY_PREFIX + address)
}

// DatastoreEntry is a helper function to fetch datastore entry from the Massa node.
func DatastoreEntry(contractAddress string, key []byte) ([]byte, error) {
	client, err := NewMassaClient()
	if err != nil {
		return nil, err
	}

	entry, err := node.FetchDatastoreEntry(client, contractAddress, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch datastore entry: %w", err)
	}
	return entry.CandidateValue, nil
}
