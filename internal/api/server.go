package api

import (
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-core/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-core/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-core/internal/api/wallet"
)

func StartServer(app *fyne.App) {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewMassaCoreAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map // to be replaced by channel

	localAPI.RestWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.RestWalletDeleteHandler = wallet.NewDelete(&walletStorage)
	localAPI.RestWalletImportHandler = wallet.NewImport(&walletStorage)
	localAPI.RestWalletListHandler = wallet.NewList(&walletStorage)
	// localAPI.RestWalletGetHandler doesn't exist in Thyra
	localAPI.RestWalletSignOperationHandler = wallet.NewSign(&walletStorage, app)
	localAPI.WebWalletHandler = operations.WebWalletHandlerFunc(WebWalletHandler)

	server.ConfigureAPI()

	defer (*app).Quit()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
