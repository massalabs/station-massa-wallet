package assets

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/pkg/errors"
)

// AssetsStore encapsulates the contract assets and mutex.
type AssetsStore struct {
	ContractAssets      map[string]models.AssetInfo
	ContractAssetsMutex sync.Mutex
}

type assetsData struct {
	Assets []assetData `json:"assets"`
}

// Define the struct for asset information
type assetData struct {
	ContractAddress string `json:"contractAddress"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Decimals        int64  `json:"decimals"`
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

func createJSONFile(path string) error {
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		return err
	}
	return nil
}

func (s *AssetsStore) loadContractAssets() error {
	assetsJSONPath, err := GetAssetsJSONPath()
	if err != nil {
		return errors.Wrap(err, "error getting assets JSON file")
	}

	// Check if the file exists, and if not, create a new one with an empty object
	if _, err := os.Stat(assetsJSONPath); os.IsNotExist(err) {
		if err := createJSONFile(assetsJSONPath); err != nil {
			return errors.Wrap(err, "failed to create assets JSON file")
		}
	}

	file, err := os.Open(assetsJSONPath)
	if err != nil {
		return errors.Wrap(err, "failed to open assets JSON file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read assets JSON data")
	}

	// Unmarshal the JSON data into the assetData struct
	var assetsData assetsData
	if err := json.Unmarshal(data, &assetsData); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data")
	}

	for _, asset := range assetsData.Assets {
		assetInfo := models.AssetInfo{
			Name:     asset.Name,
			Symbol:   asset.Symbol,
			Decimals: &asset.Decimals,
		}

		s.ContractAssets[asset.ContractAddress] = assetInfo
	}

	return nil
}

// AssetExists checks if the asset information exists for a given contract address in the JSON.
func (s *AssetsStore) AssetExists(contractAddress string) bool {
	// Look up the asset information in the ContractAssets map
	_, found := s.ContractAssets[contractAddress]
	return found
}

func (s *AssetsStore) AddAsset(assetAddress string, assetInfo models.AssetInfo) error {
	// Update the ContractAssets map with the new asset information
	s.updateAssets(assetAddress, assetInfo)

	// Convert the ContractAssets map to the format of assetsData
	var assets assetsData
	for contractAddress, assetInfo := range s.ContractAssets {
		asset := assetData{
			ContractAddress: contractAddress,
			Name:            assetInfo.Name,
			Symbol:          assetInfo.Symbol,
			Decimals:        *assetInfo.Decimals,
		}
		assets.Assets = append(assets.Assets, asset)
	}

	// Marshal the assetsData to JSON data
	data, err := json.MarshalIndent(assets, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal ContractAssets to JSON data")
	}

	// Write the JSON data to the file
	assetsJSONPath, err := GetAssetsJSONPath()
	if err != nil {
		return errors.Wrap(err, "error getting assets JSON file")
	}

	if err := os.WriteFile(assetsJSONPath, data, 0o644); err != nil {
		return errors.Wrap(err, "failed to write JSON data to file")
	}

	return nil
}

// updateAssets updates the Assets map with the new asset information.
func (s *AssetsStore) updateAssets(assetAddress string, assetInfo models.AssetInfo) {
	s.ContractAssetsMutex.Lock()
	defer s.ContractAssetsMutex.Unlock()
	s.ContractAssets[assetAddress] = assetInfo
}

// GetAssetsJSONPath returns the path to the assets JSON file.
func GetAssetsJSONPath() (string, error) {
	walletDir, err := wallet.GetWalletDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(walletDir, "assets.json"), nil
}

// AssetInfo retrieves the asset information for a given contract address by making a smart contract call.
func (s *AssetsStore) AssetInfo(contractAddress string) (*models.AssetInfo, error) {
	// mocking for now implementation to call the smart contract and retrieve asset information goes here.
	// todo : Replace the following dummy assetInfo with the actual asset information retrieved from the smart contract.

	assetInfoFromSC := &models.AssetInfo{
		Name:     "Token SC",
		Symbol:   "TKSC",
		Decimals: swag.Int64(6),
	}

	return assetInfoFromSC, nil
}
