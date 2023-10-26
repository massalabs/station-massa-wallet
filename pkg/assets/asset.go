package assets

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/pkg/errors"
)

// AssetsStore encapsulates all the nicknames with their related contract assets.
type AssetsStore struct {
	Assets     map[string]Assets
	StoreMutex sync.Mutex
}

// Assets encapsulates the contract assets associated with a specific account.
type Assets struct {
	ContractAssets map[string]models.AssetInfo
}

// assetsData represents the data structure for asset information in JSON format.
type assetsData struct {
	Assets []assetData `json:"assets"`
}

// assetData defines the structure for asset information in JSON format.
type assetData struct {
	ContractAddress string `json:"contractAddress"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Decimals        int64  `json:"decimals"`
}

// NewAssetsStore creates and initializes a new instance of AssetsStore.
func NewAssetsStore(assetsJSONPath string) (*AssetsStore, error) {
	store := &AssetsStore{
		Assets: make(map[string]Assets),
	}
	if err := store.loadaccountsStore(assetsJSONPath); err != nil {
		return nil, errors.Wrap(err, "failed to create AssetsStore")
	}

	return store, nil
}

// loadaccountsStore loads the data from the assets JSON file into the AssetsStore.
func (s *AssetsStore) loadaccountsStore(assetsJSONPath string) error {
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

	// Unmarshal the JSON data into the accountsData struct
	var accountsData struct {
		Accounts map[string]struct {
			Assets []assetData `json:"assets"`
		} `json:"accounts"`
	}

	if err := json.Unmarshal(data, &accountsData); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data")
	}

	for accountName, accountData := range accountsData.Accounts {
		accountAssets := Assets{
			ContractAssets: make(map[string]models.AssetInfo),
		}

		for _, asset := range accountData.Assets {
			assetInfo := models.AssetInfo{
				Address:  asset.ContractAddress,
				Name:     asset.Name,
				Symbol:   asset.Symbol,
				Decimals: &asset.Decimals,
			}
			accountAssets.ContractAssets[asset.ContractAddress] = assetInfo
		}

		s.Assets[accountName] = accountAssets
	}

	return nil
}

// AssetExists checks if the asset information exists for a given contract address in the JSON.
func (s *AssetsStore) AssetExists(nickname, contractAddress string) bool {
	s.StoreMutex.Lock()
	defer s.StoreMutex.Unlock()

	// Check if the account exists in the accountsStore
	accountAssets, found := s.Assets[nickname]
	if !found {
		return false
	}

	// Look up the asset information in the ContractAssets map of the specific account
	_, assetFound := accountAssets.ContractAssets[contractAddress]

	return assetFound
}

// save converts the AssetsStore map to JSON format and writes it to the file.
func (s *AssetsStore) save() error {
	// Convert the AssetsStore map to the format of accountsData
	accountsData := struct {
		Accounts map[string]struct {
			Assets []assetData `json:"assets"`
		} `json:"accounts"`
	}{
		Accounts: make(map[string]struct {
			Assets []assetData `json:"assets"`
		}),
	}

	for accountName, accountAssets := range s.Assets {
		var assetsData assetsData

		for contractAddress, assetInfo := range accountAssets.ContractAssets {
			asset := assetData{
				ContractAddress: contractAddress,
				Name:            assetInfo.Name,
				Symbol:          assetInfo.Symbol,
				Decimals:        *assetInfo.Decimals,
			}
			assetsData.Assets = append(assetsData.Assets, asset)
		}
		accountsData.Accounts[accountName] = struct {
			Assets []assetData `json:"assets"`
		}{Assets: assetsData.Assets}
	}

	// Marshal the accountsData to JSON data
	data, err := json.MarshalIndent(accountsData, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal AccountsData to JSON data")
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

// AddAsset adds the asset information for a given account nickname in the JSON.
func (s *AssetsStore) AddAsset(nickname, assetAddress string, assetInfo models.AssetInfo) error {
	// Update the ContractAssets map with the new asset information
	s.AddAssetToMemory(nickname, assetAddress, assetInfo)

	// Synchronize the AssetsStore map to JSON and write to the file
	if err := s.save(); err != nil {
		return err
	}

	return nil
}

// AddAssetToMemory adds the asset information for a given account nickname to the AssetsStore.
func (s *AssetsStore) AddAssetToMemory(nickname, assetAddress string, assetInfo models.AssetInfo) {
	s.StoreMutex.Lock()
	defer s.StoreMutex.Unlock()

	// Check if the accountAssets exists in the accountsAssets map
	accountAssets, found := s.Assets[nickname]
	if !found {
		// If the accountAssets does not exist, initialize it with an empty map
		accountAssets = Assets{
			ContractAssets: make(map[string]models.AssetInfo),
		}
	}

	// Update the ContractAssets map of the specific *assets.AssetsStore with the new asset information
	accountAssets.ContractAssets[assetAddress] = assetInfo
	s.Assets[nickname] = accountAssets
}

// DeleteAssetFromMemory removes the asset information for a given account nickname and asset address from the AssetsStore.
func (s *AssetsStore) DeleteAssetFromMemory(nickname, assetAddress string) {
	s.StoreMutex.Lock()
	defer s.StoreMutex.Unlock()

	// Check if the accountAssets exists in the accountsAssets map
	accountAssets, found := s.Assets[nickname]
	if !found {
		// If the accountAssets does not exist, there's nothing to delete, so return early.
		return
	}

	// Delete the asset from the ContractAssets map of the specific *assets.AssetsStore
	delete(accountAssets.ContractAssets, assetAddress)

	// Update the asset information in the AssetsStore
	s.Assets[nickname] = accountAssets
}

// DeleteAsset deletes the asset information for a given account nickname in the JSON.
func (s *AssetsStore) DeleteAsset(nickname, assetAddress string) error {
	s.DeleteAssetFromMemory(nickname, assetAddress)

	// Synchronize the AssetsStore map to JSON and write to the file
	if err := s.save(); err != nil {
		return err
	}

	return nil
}

// createJSONFile creates an empty JSON file at the specified path.
func createJSONFile(path string) error {
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		return err
	}

	return nil
}

// GetAssetsJSONPath returns the path to the assets JSON file.
func GetAssetsJSONPath() (string, error) {
	walletPath, err := wallet.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(walletPath, "assets.json"), nil
}

func MASInfo() models.AssetInfo {
	// The hardcoded data for MAS asset
	name := "Massa"
	symbol := "MAS"
	decimals := int64(9)

	// Create the AssetInfo struct with the predefined information
	assetInfo := models.AssetInfo{
		Name:     name,
		Symbol:   symbol,
		Decimals: &decimals,
	}

	return assetInfo
}
