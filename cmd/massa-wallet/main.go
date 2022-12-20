package main

import (
	"flag"
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
	var port int
	var path string
	flag.IntVar(&port, "port", 8080, "the port to listen on")
	flag.StringVar(&path, "path", "", "the path to listen on")
	flag.Parse()

	myApp := app.New()
	go StartServer(&myApp, port)

	myApp.Run()
}

func StartServer(app *fyne.App, port int) {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map // to be replaced by channel

	localAPI.RestWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.RestWalletDeleteHandler = wallet.NewDelete(&walletStorage)
	localAPI.RestWalletImportHandler = wallet.NewImport(&walletStorage)
	localAPI.RestWalletListHandler = wallet.NewList(&walletStorage)

	localAPI.RestWalletSignOperationHandler = wallet.NewSign(&walletStorage, app)
	localAPI.WebHandler = operations.WebHandlerFunc(handler.WebWalletHandler)

	server.ConfigureAPI()

	// Set the port to listen on to the passed-in port
	server.Port = port

	defer (*app).Quit()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
