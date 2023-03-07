package main

import (
	"log"
	"net"
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

	ln, err := server.HTTPListener()
	if err != nil {
		log.Fatalln(err)
	}
	registerPlugin(ln, "", "logo.svg")

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func registerPlugin(ln net.Listener, spec string, logoURL string) {
	if os.Getenv("STANDALONE") == "1" {
		return
	}

	if len(os.Args) >= 2 {
		err := plugin.Register(os.Args[1], constants.PluginName, constants.PluginAuthor, constants.PluginDescription,
			ln.Addr(), spec, logoURL)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		panic("Usage: program must be started with an ID command line argument")
	}
}
