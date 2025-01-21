package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
)

func NewDeleteAsset() operations.DeleteAssetHandler {
	return &deleteAsset{}
}

type deleteAsset struct{}

func (d *deleteAsset) Handle(params operations.DeleteAssetParams) middleware.Responder {
	// Check if the address is valid
	if !utils.IsValidAddress(params.AssetAddress) {
		// Return an error indicating the address is not valid
		errorMsg := "Invalid address format"
		return operations.NewDeleteAssetBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: errorMsg})
	}

	// Check if the asset exists in the loaded JSON
	if !assets.Store.AssetExists(params.Nickname, params.AssetAddress) {
		// Return an error indicating that the asset does not exist
		errorMsg := "Asset with the provided address does not exist."
		return operations.NewDeleteAssetBadRequest().WithPayload(&models.Error{Code: errorAssetNotExists, Message: errorMsg})
	}

	// Delete Asset From the JSON file.
	if err := assets.Store.DeleteAsset(params.Nickname, params.AssetAddress); err != nil {
		// Return error occurred while persisting the asset
		errorMsg := "Failed to delete the asset from the JSON file."
		return operations.NewDeleteAssetInternalServerError().WithPayload(&models.Error{Code: errorDeleteAssetJSON, Message: errorMsg})
	}

	response := operations.NewDeleteAssetNoContent()

	return response
}
