package handler

import (
	"github.com/go-openapi/loads"
	"github.com/massalabs/station-massa-wallet/api/server/restapi"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/internal/handler/html"
	walletHandler "github.com/massalabs/station-massa-wallet/internal/handler/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
)

// InitializeAPI initializes the API handlers
func InitializeAPI(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) (*operations.MassaWalletAPI, error) {
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
	walletHandler.AppendEndpoints(api, prompterApp, massaClient)

	return api, nil
}
