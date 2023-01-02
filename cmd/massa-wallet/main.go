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
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/gui"
)

func main() {
	var port int
	var path string
	flag.IntVar(&port, "port", 8080, "the port to listen on")
	// The path is not actually used in the script.
	// It is included only to maintain temporary compatibility with Thyra, and it will be removed at a later time.
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

	pwdPrompter := gui.NewPasswordPrompt(app)

	localAPI.RestWalletSignOperationHandler = wallet.NewSign(pwdPrompter.Ask)
	localAPI.WebHandler = operations.WebHandlerFunc(handler.WebWalletHandler)
	localAPI.DefaultPageHandler = operations.DefaultPageHandlerFunc(handler.DefaultPageHandler)

	server.ConfigureAPI()

	// Set the port to listen on to the passed-in port
	server.Port = port

	defer (*app).Quit()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
