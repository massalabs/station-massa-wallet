package assets

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/logger"
	"github.com/pkg/errors"
)

const (
	permissionUrwGrOr = 0o644
	assetsFilename    = "assets.json"
)

// AssetsStore encapsulates all the nicknames with their related contract assets.
type AssetsStore struct {
	Assets        map[string]Assets
	StoreMutex    sync.Mutex
	massaClient   *network.NodeFetcher
	assetsJSONDir string
}

type AssetInfoWithBalances struct {
	AssetInfo   *models.AssetInfo
	Balance     string
	MEXCSymbol  string
	DollarValue *float64
	IsDefault   bool
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
	ChainID         int64  `json:"chainID"`
}

var (
	Store            *AssetsStore
	filePathOverride string
)

// NewAssetsStore creates and initializes a new instance of AssetsStore.
// If assetsJSONDir is empty, it will use the default wallet path.
func InitAssetsStore(massaClient *network.NodeFetcher) *AssetsStore {
	Store = &AssetsStore{
		Assets:      make(map[string]Assets),
		massaClient: massaClient,
	}

	if filePathOverride != "" {
		Store.assetsJSONDir = filePathOverride
	} else {
		assetsJSONDir, err := wallet.Path()
		if err != nil {
			logger.Fatalf("Failed to create AssetsStore: Failed to get AssetsStore JSON file", err)
		}
		Store.assetsJSONDir = assetsJSONDir
	}

	if err := Store.loadAccountsStore(); err != nil {
		logger.Fatalf("Failed to create AssetsStore:", err)
	}

	err := Store.InitDefault()
	if err != nil {
		logger.Fatalf("Failed to create AssetsStore:", err)
	}

	return Store
}

// Used by unit test
func SetFileDirOverride(path string) {
	filePathOverride = path
}

// loadAccountsStore loads the data from the assets JSON file into the AssetsStore.
func (s *AssetsStore) loadAccountsStore() error {
	// Check if the file exists, and if not, create a new one with an empty object
	assetsJSONPath := getAssetJSONPath(s.assetsJSONDir)
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
			decimals := asset.Decimals // Copy the decimals value to avoid a pointer to a local variable
			assetInfo := models.AssetInfo{
				Address:  asset.ContractAddress,
				Name:     asset.Name,
				Symbol:   asset.Symbol,
				Decimals: &decimals,
				ChainID:  swag.Int64(asset.ChainID),
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
				ChainID:         *assetInfo.ChainID,
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
	if err := os.WriteFile(getAssetJSONPath(s.assetsJSONDir), data, permissionUrwGrOr); err != nil {
		return errors.Wrap(err, "failed to write JSON data to file")
	}

	return nil
}

// AddAsset adds the asset information for a given account nickname in the JSON.
func (s *AssetsStore) AddAsset(nickname string, assetInfo models.AssetInfo) error {
	// Update the ContractAssets map with the new asset information
	s.AddAssetToMemory(nickname, assetInfo)

	// Synchronize the AssetsStore map to JSON and write to the file
	if err := s.save(); err != nil {
		return err
	}

	return nil
}

// AddAssetToMemory adds the asset information for a given account nickname to the AssetsStore.
func (s *AssetsStore) AddAssetToMemory(nickname string, assetInfo models.AssetInfo) {
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
	accountAssets.ContractAssets[assetInfo.Address] = assetInfo
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

// All returns all the assets associated with a specific account nickname.
// It returns the default assets first, followed by the assets added by the user.
// If user already added the default asset, it will not be duplicated.
func (s *AssetsStore) All(nickname string, chainID int) []*AssetInfoWithBalances {
	defaultAssets, err := s.Default()
	if err != nil {
		logger.Errorf("Failed to get default assets: %s", err.Error())
	}

	assetsInfo := make([]*AssetInfoWithBalances, 0)

	// Initialize map to track addressed already added
	includedAddresses := map[string]bool{}

	for _, asset := range defaultAssets {
		if asset.ChainID != chainID {
			continue
		}
		decimals := asset.Decimals // Copy the decimals value to avoid a pointer to a local variable
		completeAsset := &AssetInfoWithBalances{
			AssetInfo: &models.AssetInfo{
				Address:  asset.Address,
				Decimals: &decimals,
				Name:     asset.Name,
				Symbol:   asset.Symbol,
				ChainID:  swag.Int64((int64(asset.ChainID))),
			},
			Balance:     "",
			MEXCSymbol:  asset.MEXCSymbol,
			DollarValue: nil,
			IsDefault:   true,
		}
		assetsInfo = append(assetsInfo, completeAsset)
		includedAddresses[asset.Address] = true
	}

	// Append default assets ensuring no duplication
	for _, asset := range s.Assets[nickname].ContractAssets {
		// skip if it's already in the list or not in the same chain
		_, exists := includedAddresses[asset.Address]
		if *asset.ChainID != int64(chainID) || exists {
			continue
		}
		completeAsset := &AssetInfoWithBalances{
			AssetInfo: &models.AssetInfo{
				Address:  asset.Address,
				Decimals: asset.Decimals,
				Name:     asset.Name,
				Symbol:   asset.Symbol,
				ChainID:  asset.ChainID,
			},
			Balance:     "",
			MEXCSymbol:  "",
			DollarValue: nil,
			IsDefault:   false,
		}
		assetsInfo = append(assetsInfo, completeAsset)
		includedAddresses[asset.Address] = true
	}

	return assetsInfo
}

func getAssetJSONPath(assetsJSONDir string) string {
	return filepath.Join(assetsJSONDir, assetsFilename)
}

// createJSONFile creates an empty JSON file at the specified path.
func createJSONFile(path string) error {
	if err := os.WriteFile(path, []byte("{}"), permissionUrwGrOr); err != nil {
		return err
	}

	return nil
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
