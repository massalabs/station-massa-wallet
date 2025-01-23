package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
)

type AssetInfoListResponse struct {
	Assets []models.AssetInfo `json:"assets"`
}

func NewAddAsset(massaClient network.NodeFetcherInterface) operations.AddAssetHandler {
	return &addAsset{
		massaClient: massaClient,
	}
}

type addAsset struct {
	massaClient network.NodeFetcherInterface
}

func (a *addAsset) Handle(params operations.AddAssetParams) middleware.Responder {
	if !utils.IsValidContract(params.AssetAddress) {
		errorMsg := "Invalid address format"
		return operations.NewAddAssetUnprocessableEntity().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: errorMsg})
	}

	// First, check if the asset exists in the network
	if !a.massaClient.AssetExistInNetwork(params.AssetAddress) {
		// If the asset does not exist in the network, return a 404 response
		errorMsg := "Asset with the provided address not found in the network."
		return operations.NewAddAssetNotFound().WithPayload(&models.Error{Code: errorAssetNotFound, Message: errorMsg})
	}

	// Check if the address exists in the loaded JSON
	if assets.Store.AssetExists(params.Nickname, params.AssetAddress) {
		// Return that the asset already exists
		errorMsg := "Asset with the provided address already exists."
		return operations.NewAddAssetBadRequest().WithPayload(&models.Error{Code: errorAssetExists, Message: errorMsg})
	}

	// Fetch the asset information from the SC
	assetInfoFromSC, err := assets.AssetInfo(params.AssetAddress, a.massaClient)
	if err != nil {
		// Return error occurred during SC fetch
		errorMsg := "Failed to fetch asset information from the smart contract."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorFetchAssetSC, Message: errorMsg})
	}

	// Add Asset and persist in JSON file.
	if err := assets.Store.AddAsset(params.Nickname, *assetInfoFromSC); err != nil {
		// Return error occurred while persisting the asset
		errorMsg := "Failed to add the asset to the JSON file."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorAddAssetJSON, Message: errorMsg})
	}

	// Return success response with the retrieved asset information
	return operations.NewAddAssetCreated().WithPayload(assetInfoFromSC)
}
