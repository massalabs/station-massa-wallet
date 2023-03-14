package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra-plugin-hello-world/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
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
	//mandatory to free main thread
	defer (*app).Quit()

	// Initializes API
	massaWalletAPI, err := handler.InitializeAPI(password.NewFynePrompter(app), privateKey.NewFynePrompter(app))
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
		Description: constants.PluginDescription, APISpec: "", Logo: "web/html/logo.png",
	})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
