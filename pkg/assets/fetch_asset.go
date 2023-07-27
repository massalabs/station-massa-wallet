package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/labstack/gommon/log"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/pkg/errors"
)

const assetsJSONPath = "assets.json"

type assetData struct {
	Assets []struct {
		ContractAddress string `json:"contractAddress"`
		Name            string `json:"name"`
		Symbol          string `json:"symbol"`
		Decimals        int64  `json:"decimals"`
	} `json:"assets"`
}

// AssetsStore encapsulates the contract assets and mutex.
type AssetsStore struct {
	ContractAssets      map[string]models.AssetInfo
	contractAssetsMutex sync.Mutex
}

// RetrieveAssetsInfo retrieves asset information for a given list of contract addresses concurrently.
func (s *AssetsStore) RetrieveAssetsInfo(contractAddresses []string) ([]models.AssetInfo, error) {
	var wg sync.WaitGroup
	assetInfoCh := make(chan *models.AssetInfo, len(contractAddresses))
	errCh := make(chan error, len(contractAddresses))

	for _, address := range contractAddresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()

			assetInfo, err := s.RetrieveAssetInfo(address)
			if err != nil {
				errCh <- errors.Wrapf(err, "failed to retrieve asset information for contract address %s", address)
				return
			}

			assetInfoCh <- assetInfo
		}(address)
	}

	wg.Wait()
	close(assetInfoCh)
	close(errCh)

	// Collect asset information and errors
	var assetInfos []models.AssetInfo
	var allErrors []error

	for assetInfo := range assetInfoCh {
		assetInfos = append(assetInfos, *assetInfo)
	}

	for err := range errCh {
		allErrors = append(allErrors, err)
	}

	// If there were errors, return them along with the retrieved asset information
	if len(allErrors) > 0 {
		return assetInfos, fmt.Errorf("failed to retrieve some assets: %w", gatherErrors(allErrors))
	}

	return assetInfos, nil
}

// RetrieveAssetInfo retrieves the asset information for a given contract address.
func (s *AssetsStore) RetrieveAssetInfo(contractAddress string) (*models.AssetInfo, error) {
	// Check if the asset information exists in the JSON
	assetInfo, err := s.AssetInfo(contractAddress)
	if err == nil {
		// Asset information found in the JSON
		return assetInfo, nil
	} else if errors.Is(err, ErrAssetInfoNotFound) {
		// If asset information is not found in the JSON, call the smart contract to retrieve it
		assetInfoFromSC, err := s.getAssetInfoFromSC(contractAddress)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve asset information from smart contract")
		}

		// Update the ContractAssets map with the new asset information
		s.contractAssetsMutex.Lock()
		defer s.contractAssetsMutex.Unlock()

		s.ContractAssets[contractAddress] = *assetInfoFromSC

		// Persist the updated ContractAssets map back to the JSON file
		if err := s.persistAssets(); err != nil {
			log.Errorf("failed to persist contract assets to JSON file: %v", err)
		}

		return assetInfoFromSC, nil
	}

	return nil, errors.Wrap(err, "failed to retrieve asset information")
}

// gatherErrors creates a new MultiError from a slice of errors.
func gatherErrors(errorsSlice []error) error {
	var errMsg string
	for _, err := range errorsSlice {
		errMsg += err.Error() + "; "
	}
	return errors.New(errMsg)
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

	// Unmarshal the JSON data into the assetData struct
	var assetsData assetData
	if err := json.Unmarshal(data, &assetsData); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data")
	}

	// Update the ContractAssets map with the new asset information
	s.contractAssetsMutex.Lock()
	defer s.contractAssetsMutex.Unlock()

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

// AssetInfo retrieves the asset information for a given contract address from the JSON.
func (s *AssetsStore) AssetInfo(contractAddress string) (*models.AssetInfo, error) {
	// Look up the asset information in the ContractAssets map
	assetInfo, found := s.ContractAssets[contractAddress]
	if !found {
		return nil, errors.Wrap(ErrAssetInfoNotFound, fmt.Sprintf("asset information not found for contract address: %s", contractAddress))
	}

	return &assetInfo, nil
}

// ErrAssetInfoNotFound is an error indicating that the asset information is not found in the JSON.
var ErrAssetInfoNotFound = errors.New("asset information not found")

// persistAssets saves the ContractAssets map to the specified JSON file.
func (s *AssetsStore) persistAssets() error {
	// Convert the ContractAssets map to the format of assetData
	var assetsData assetData
	for contractAddress, assetInfo := range s.ContractAssets {
		asset := struct {
			ContractAddress string `json:"contractAddress"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
			Decimals        int64  `json:"decimals"`
		}{
			ContractAddress: contractAddress,
			Name:            assetInfo.Name,
			Symbol:          assetInfo.Symbol,
			Decimals:        *assetInfo.Decimals,
		}
		assetsData.Assets = append(assetsData.Assets, asset)
	}

	// Marshal the assetData to JSON data
	data, err := json.MarshalIndent(assetsData, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal ContractAssets to JSON data")
	}

	// Write the JSON data to the file
	if err := os.WriteFile(assetsJSONPath, data, 0o644); err != nil {
		return errors.Wrap(err, "failed to write JSON data to file")
	}

	return nil
}

// getAssetInfoFromSC retrieves the asset information for a given contract address by making a smart contract call.
func (s *AssetsStore) getAssetInfoFromSC(contractAddress string) (*models.AssetInfo, error) {
	// mocking for now implementation to call the smart contract and retrieve asset information goes here.
	// todo : Replace the following dummy assetInfo with the actual asset information retrieved from the smart contract.

	assetInfoFromSC := &models.AssetInfo{
		Name:     "Token C",
		Symbol:   "TKC",
		Decimals: swag.Int64(6),
	}

	return assetInfoFromSC, nil
}
