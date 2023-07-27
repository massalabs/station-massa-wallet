package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
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
	// fetch the assets information for the list of contract addresses
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
