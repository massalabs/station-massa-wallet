package handler

import (
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/internal/handler/html"
	"github.com/massalabs/thyra-plugin-massa-wallet/internal/handler/wallet"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/guiModal"
)

// InitializeAPI initializes the API handlers
func InitializeAPI(passwordPrompter guiModal.PasswordAsker, walletInfoModal guiModal.WalletInfoAsker) (*operations.MassaWalletAPI, error) {
	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	// Create a new MassaWalletAPI instance
	api := operations.NewMassaWalletAPI(swaggerSpec)

	// Set web endpoints
	html.AppendEndpoints(api)

	// Set wallet API endpoints
	wallet.AppendEndpoints(api, walletInfoModal, passwordPrompter)

	return api, nil
}
