package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-hello-world/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
	"github.com/massalabs/thyra-plugin-wallet/pkg/delete"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
	constants "github.com/massalabs/thyra-plugin-wallet/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/pkg/privateKey"
)

func main() {
	myApp := app.New()
	go startServer(&myApp)

	myApp.Run()
}

func startServer(app *fyne.App) {
	// mandatory to free main thread
	defer (*app).Quit()

	// Initialize cache
	gc := gcache.New(20).
		LRU().
		Build()

	var passwordPrompter password.Asker = password.NewFynePrompter(app)

	// Check if the wallet is running in test mode
	if os.Getenv("WALLET_TEST_MODE") == "1" {
		log.Println("Wallet is running in test mode")
		passwordPrompter = password.NewEnvPrompter()
	}

	// Initializes API
	massaWalletAPI, err := handler.InitializeAPI(
		passwordPrompter,
		privateKey.NewFynePrompter(app),
		delete.NewFynePrompter(app),
		gc,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// instantiates and configure server
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	if os.Getenv("STANDALONE") == "1" {
		server.Port = 8080
	}

	listener, err := server.HTTPListener()
	if err != nil {
		log.Fatalln(err)
	}

	plugin.RegisterPlugin(listener, plugin.Info{
		Name: constants.PluginName, Author: constants.PluginAuthor,
		Description: constants.PluginDescription, APISpec: "", Logo: "web/html/wallet.svg",
	})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
