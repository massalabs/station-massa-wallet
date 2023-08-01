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

// AccountsStore encapsulates all the accounts with their related contract assets.
type AccountsStore struct {
	Accounts   map[string]Account
	StoreMutex sync.Mutex
}

// Account encapsulates the contract assets associated with a specific wallet.
type Account struct {
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

// NewAccountsStore creates and initializes a new instance of AccountsStore.
func NewAccountsStore() (*AccountsStore, error) {
	store := &AccountsStore{
		Accounts: make(map[string]Account),
	}
	if err := store.loadWalletsStore(); err != nil {
		return nil, errors.Wrap(err, "failed to create AccountsStore")
	}

	return store, nil
}

// loadWalletsStore loads the data from the assets JSON file into the AccountsStore.
func (s *AccountsStore) loadWalletsStore() error {
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

	// Unmarshal the JSON data into the walletsData struct
	var accountsData struct {
		Accounts map[string]struct {
			Assets []assetData `json:"assets"`
		} `json:"wallets"`
	}
	if err := json.Unmarshal(data, &accountsData); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data")
	}

	for walletName, walletData := range accountsData.Accounts {
		Account := Account{
			ContractAssets: make(map[string]models.AssetInfo),
		}

		for _, asset := range walletData.Assets {
			assetInfo := models.AssetInfo{
				Name:     asset.Name,
				Symbol:   asset.Symbol,
				Decimals: &asset.Decimals,
			}
			Account.ContractAssets[asset.ContractAddress] = assetInfo
		}

		s.Accounts[walletName] = Account
	}

	return nil
}

// AssetExists checks if the asset information exists for a given contract address in the JSON.
func (s *AccountsStore) AssetExists(walletNickname, contractAddress string) bool {
	s.StoreMutex.Lock()
	defer s.StoreMutex.Unlock()

	// Check if the wallet exists in the WalletsStore
	walletAssets, found := s.Accounts[walletNickname]
	if !found {
		return false
	}

	// Look up the asset information in the ContractAssets map of the specific wallet
	_, assetFound := walletAssets.ContractAssets[contractAddress]
	return assetFound
}

// AddAsset adds the asset information for a given wallet nickname in the JSON.
func (s *AccountsStore) AddAsset(walletNickname, assetAddress string, assetInfo models.AssetInfo) error {
	// Update the ContractAssets map with the new asset information
	s.updateAccounts(walletNickname, assetAddress, assetInfo)

	// Convert the Accounts map to the format of accountsData
	accountsData := struct {
		Accounts map[string]struct {
			Assets []assetData `json:"assets"`
		} `json:"wallets"`
	}{
		Accounts: make(map[string]struct {
			Assets []assetData `json:"assets"`
		}),
	}

	for walletName, walletAssets := range s.Accounts {
		var assetsData assetsData
		for contractAddress, assetInfo := range walletAssets.ContractAssets {
			asset := assetData{
				ContractAddress: contractAddress,
				Name:            assetInfo.Name,
				Symbol:          assetInfo.Symbol,
				Decimals:        *assetInfo.Decimals,
			}
			assetsData.Assets = append(assetsData.Assets, asset)
		}
		accountsData.Accounts[walletName] = struct {
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

// updateAccounts updates the Accounts map with the new asset information.
func (s *AccountsStore) updateAccounts(walletNickname, assetAddress string, assetInfo models.AssetInfo) {
	s.StoreMutex.Lock()
	defer s.StoreMutex.Unlock()

	// Check if the account exists in the Accounts map
	account, found := s.Accounts[walletNickname]
	if !found {
		// If the account does not exist, initialize it with an empty map
		account = Account{
			ContractAssets: make(map[string]models.AssetInfo),
		}
	}

	// Update the ContractAssets map of the specific account with the new asset information
	account.ContractAssets[assetAddress] = assetInfo
	s.Accounts[walletNickname] = account
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
	walletDir, err := wallet.GetWalletDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(walletDir, "assets.json"), nil
}
