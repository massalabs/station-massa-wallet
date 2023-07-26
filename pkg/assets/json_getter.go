package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/pkg/errors"
)

const assetsJSONPath = "assets.json"

// AssetsStore encapsulates the contract assets and mutex.
type AssetsStore struct {
	ContractAssets      map[string]models.AssetInfo
	contractAssetsMutex sync.Mutex
}

// NewAssetsStore creates and initializes a new instance of AssetsStore.
func NewAssetsStore() (*AssetsStore, error) {
	store := &AssetsStore{
		ContractAssets: make(map[string]models.AssetInfo),
	}
	if err := store.loadContractAssets(); err != nil {
		return nil, errors.Wrap(err, "failed to create AssetsStore")
	}
	return store, nil
}

// loadContractAssets loads the JSON file into the ContractAssets map.
func (s *AssetsStore) loadContractAssets() error {
	file, err := os.Open(assetsJSONPath)
	if err != nil {
		return errors.Wrap(err, "failed to open assets JSON file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read assets JSON data")
	}

	// Unmarshal the JSON data into the ContractAssets map
	if err := json.Unmarshal(data, &s.ContractAssets); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data into ContractAssets")
	}

	return nil
}

// getAssetInfoFromJSON retrieves the asset information for a given contract address from the JSON.
func (s *AssetsStore) getAssetInfoFromJSON(contractAddress string) (*models.AssetInfo, error) {
	// Look up the asset information in the ContractAssets map
	assetInfo, found := s.ContractAssets[contractAddress]
	if !found {
		return nil, errors.Wrap(ErrAssetInfoNotFound, fmt.Sprintf("asset information not found for contract address: %s", contractAddress))
	}

	return &assetInfo, nil
}

// ErrAssetInfoNotFound is an error indicating that the asset information is not found in the JSON.
var ErrAssetInfoNotFound = errors.New("asset information not found")

// persistContractAssets saves the ContractAssets map to the specified JSON file.
func (s *AssetsStore) persistContractAssets() error {
	// Convert the ContractAssets map to JSON data
	data, err := json.MarshalIndent(s.ContractAssets, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal ContractAssets to JSON data")
	}

	// Write the JSON data to the file
	if err := os.WriteFile(assetsJSONPath, data, 0o644); err != nil {
		return errors.Wrap(err, "failed to write JSON data to file")
	}

	return nil
}
