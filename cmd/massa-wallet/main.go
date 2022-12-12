package main

import (
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/internal/handler"
	"github.com/massalabs/thyra-plugin-massa-wallet/internal/handler/wallet"
)

func main() {
	myApp := app.New()

	go StartServer(&myApp)

	myApp.Run()
}

func StartServer(app *fyne.App) {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map // to be replaced by channel

	localAPI.RestWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.RestWalletDeleteHandler = wallet.NewDelete(&walletStorage, app)
	localAPI.RestWalletImportHandler = wallet.NewImport(&walletStorage, app)
	localAPI.RestWalletListHandler = wallet.NewList(&walletStorage)

	localAPI.RestWalletSignOperationHandler = wallet.NewSign(&walletStorage, app)
	localAPI.WebHandler = operations.WebHandlerFunc(handler.WebWalletHandler)

	server.ConfigureAPI()

	defer (*app).Quit()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
