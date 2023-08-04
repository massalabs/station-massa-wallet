package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	address "github.com/massalabs/station/pkg/dnshelper"
)

type AssetInfoListResponse struct {
	Assets []models.AssetInfo `json:"assets"`
}

func NewAddAsset(AssetsStore *assets.AssetsStore, massaClient network.NodeFetcherInterface) operations.AddAssetHandler {
	return &addAsset{
		AssetsStore: AssetsStore,
		massaClient: massaClient,
	}
}

type addAsset struct {
	AssetsStore *assets.AssetsStore
	massaClient network.NodeFetcherInterface
}

func (h *addAsset) Handle(params operations.AddAssetParams) middleware.Responder {
	// Check if the address is valid
	if !address.IsValidAddress(params.AssetAddress) {
		// Return an error indicating the address is not valid
		errorMsg := "Invalid address format"
		return operations.NewAddAssetUnprocessableEntity().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: errorMsg})
	}

	// Check if the address exists in the loaded JSON
	if h.AssetsStore.AssetExists(params.Nickname, params.AssetAddress) {
		// Return that the asset already exists
		errorMsg := "Asset with the provided address already exists."
		return operations.NewAddAssetBadRequest().WithPayload(&models.Error{Code: errorAssetExists, Message: errorMsg})
	}

	// Fetch the asset information from the SC
	assetInfoFromSC, err := assets.AssetInfo(params.AssetAddress, h.massaClient)
	if err != nil {
		// Return error occurred during SC fetch
		errorMsg := "Failed to fetch asset information from the smart contract."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorFetchAssetSC, Message: errorMsg})
	}

	// Add Asset and persist in JSON file.
	if err := h.AssetsStore.AddAsset(params.Nickname, params.AssetAddress, *assetInfoFromSC); err != nil {
		// Return error occurred while persisting the asset
		errorMsg := "Failed to add the asset to the JSON file."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorAddAssetJSON, Message: errorMsg})
	}

	// Return success response with the retrieved asset information
	return operations.NewAddAssetCreated().WithPayload(assetInfoFromSC)
}
