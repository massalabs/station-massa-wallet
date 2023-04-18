package handler

import (
	"github.com/bluele/gcache"
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler/html"
	walletHandler "github.com/massalabs/thyra-plugin-wallet/internal/handler/wallet"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// InitializeAPI initializes the API handlers
func InitializeAPI(prompterApp *wallet.WalletPrompter, gc gcache.Cache) (*operations.MassaWalletAPI, error) {
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
	walletHandler.AppendEndpoints(api, prompterApp, gc)

	return api, nil
}
