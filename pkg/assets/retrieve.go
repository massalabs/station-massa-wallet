package assets

import (
	"fmt"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/gommon/log"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/pkg/errors"
)

type AssetInfoListResponse struct {
	Assets []models.AssetInfo `json:"assets"`
}

func NewRetrieveAssetInfoList(assetsStore *AssetsStore) operations.AddAssetHandler{
	return &retrieveAssetInfoList{
		assetsStore: assetsStore,
	}
}

type retrieveAssetInfoList struct {
	assetsStore *AssetsStore
}

func (h *retrieveAssetInfoList) Handle(params operations.AddAssetParams) middleware.Responder {
	// Retrieve the asset information for the list of contract addresses
	assetInfoList, retrievalErr := h.assetsStore.RetrieveAssetsInfo(params.AssetAddresses)

	if retrievalErr != nil {
		// If some asset information was retrieved successfully
		if len(assetInfoList) > 0 {
			// Create the response containing the partial asset information list
			response := operations.AddAssetPartialContentBody{
				Assets: assetInfoList,
			}

			// Return the response with the partial list using 206 Partial Content
			return operations.NewAddAssetPartialContent().WithPayload(&response)
		}

		// If there is no retrieved asset information, return the internal server error
		return operations.NewAddAssetInternalServerError().WithPayload(
			&models.Error{
				Code:    "internal_server_error",
				Message: "Failed to retrieve asset information",
			},
		)
	}

	// All asset information retrieved successfully

	// Create the response containing the asset information list
	response := operations.AddAssetCreatedBody{
		Assets: assetInfoList,
	}

	// Return the response with the full list of retrieved asset information using 200 OK
	return operations.NewAddAssetCreated().WithPayload(&response)
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
	assetInfo, err := s.getAssetInfoFromJSON(contractAddress)
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
