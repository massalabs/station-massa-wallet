package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra-plugin-hello-world/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
	"github.com/massalabs/thyra-plugin-wallet/pkg/delete"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
	constants "github.com/massalabs/thyra-plugin-wallet/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/pkg/privateKey"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// App struct
type WApp struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *WApp {
	return &WApp{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *WApp) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *WApp) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

//go:generate cp -r ../../interface/frontend/dist ./frontend-dist
//go:embed frontend-dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	wAapp := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "interface",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        wAapp.startup,
		Bind: []interface{}{
			wAapp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	//----------

	myApp := app.New()
	go startServer(&myApp)

	myApp.Run()
}

func startServer(app *fyne.App) {
	//mandatory to free main thread
	defer (*app).Quit()

	// Initializes API
	massaWalletAPI, err := handler.InitializeAPI(
		password.NewFynePrompter(app),
		privateKey.NewFynePrompter(app),
		delete.NewFynePrompter(app),
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
		Description: constants.PluginDescription, APISpec: "", Logo: "web/html/logo.png",
	})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
