package main

import (
	"flag"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-wallet/pkg/privateKey"
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

	// Set the port to listen on to the passed-in port
	server.Port = port

	if err := server.Serve(); err != nil {

		log.Fatalln(err)
	}
}
