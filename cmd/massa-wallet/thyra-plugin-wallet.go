package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-wallet/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/pkg/privateKey"
)

func main() {
	ln, _ := net.Listen("tcp", ":")

	registerPlugin(ln)

	myApp := app.New()
	go startServer(&myApp, ln)

	myApp.Run()
}

func startServer(app *fyne.App, ln net.Listener) {
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
	port, _ := strconv.Atoi(strings.Split(ln.Addr().String(), ":")[1])
	server.Port = port

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func registerPlugin(ln net.Listener) {
	//nolint:gomnd
	if len(os.Args) < 2 {
		panic("this program must be run with correlation id argument!")
	}

	pluginID := os.Args[1]

	standaloneMode := false

	if len(os.Args) == 3 && os.Args[2] == "--standalone" {
		standaloneMode = true
	}

	if !standaloneMode {
		err := plugin.Register(pluginID, "Massa Wallet", "Massalabs", "Massa wallet for Thyra", ln.Addr())
		if err != nil {
			log.Panicln(err)
		}
	}
}
