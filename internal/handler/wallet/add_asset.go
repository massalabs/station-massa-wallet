package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	address "github.com/massalabs/station/pkg/dnshelper"
)

type AssetInfoListResponse struct {
	Assets []models.AssetInfo `json:"assets"`
}

func NewAddAsset(assetsStore *assets.AssetsStore) operations.AddAssetHandler {
	return &addAsset{
		assetsStore: assetsStore,
	}
}

type addAsset struct {
	assetsStore *assets.AssetsStore
}

func (h *addAsset) Handle(params operations.AddAssetParams) middleware.Responder {
	// Check if the address is valid
	if !address.IsValidAddress(params.AssetAddress) {
		// Return an error indicating the address is not valid
		errorMsg := "Invalid address format"
		return operations.NewAddAssetBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: errorMsg})
	}

	// Check if the address exists in the loaded JSON
	if h.assetsStore.AssetExists(params.AssetAddress) {
		// Return that the asset already exists
		errorMsg := "Asset with the provided address already exists."
		return operations.NewAddAssetBadRequest().WithPayload(&models.Error{Code: errorAssetExists, Message: errorMsg})
	}

	// Fetch the asset information from the SC
	assetInfoFromSC, err := h.assetsStore.AssetInfo(params.AssetAddress)
	if err != nil {
		// Return error occurred during SC fetch
		errorMsg := "Failed to fetch asset information from the smart contract."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorFetchAssetSC, Message: errorMsg})
	}

	// Update the ContractAssets map with the new asset information
	h.assetsStore.ContractAssetsMutex.Lock()
	defer h.assetsStore.ContractAssetsMutex.Unlock()
	h.assetsStore.ContractAssets[params.AssetAddress] = *assetInfoFromSC

	// Add the assets to the persisted JSON file.
	if err := h.assetsStore.AddAsset(); err != nil {
		// Return error occurred while persisting the asset
		errorMsg := "Failed to add the asset to the JSON file."
		return operations.NewAddAssetInternalServerError().WithPayload(&models.Error{Code: errorAddAssetJSON, Message: errorMsg})
	}

	// Return success response with the retrieved asset information
	return operations.NewAddAssetCreated().WithPayload(assetInfoFromSC)
}