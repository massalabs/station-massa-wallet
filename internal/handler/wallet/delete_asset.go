package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	address "github.com/massalabs/station/pkg/dnshelper"
)

func NewDeleteAsset(AssetsStore *assets.AssetsStore) operations.DeleteAssetHandler {
	return &deleteAsset{
		AssetsStore: AssetsStore,
	}
}

type deleteAsset struct {
	AssetsStore *assets.AssetsStore
}

func (h *deleteAsset) Handle(params operations.DeleteAssetParams) middleware.Responder {
	// Check if the address is valid
	if !address.IsValidAddress(params.AssetAddress) {
		// Return an error indicating the address is not valid
		errorMsg := "Invalid address format"
		return operations.NewDeleteAssetBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: errorMsg})
	}

	// Check if the asset exists in the loaded JSON
	if !h.AssetsStore.AssetExists(params.Nickname, params.AssetAddress) {
		// Return an error indicating that the asset does not exist
		errorMsg := "Asset with the provided address does not exist."
		return operations.NewDeleteAssetBadRequest().WithPayload(&models.Error{Code: errorAssetNotExists, Message: errorMsg})
	}

	// Delete Asset From the JSON file.
	if err := h.AssetsStore.DeleteAsset(params.Nickname, params.AssetAddress); err != nil {
		// Return error occurred while persisting the asset
		errorMsg := "Failed to delete the asset from the JSON file."
		return operations.NewDeleteAssetInternalServerError().WithPayload(&models.Error{Code: errorDeleteAssetJSON, Message: errorMsg})
	}

	response := operations.NewDeleteAssetCreated()
	return response
}
