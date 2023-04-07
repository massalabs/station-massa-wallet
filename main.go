package main

import (
	"context"
	"embed"
	"fmt"

	fyneApp "fyne.io/fyne/v2/app"

	walletApp "github.com/massalabs/thyra-plugin-wallet/cmd/massa-wallet"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//nolint:typecheck,nolintlint
//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	fmt.Println("Greet input:" + name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func main() {

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wailds-tmpl",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	fyneApp := fyneApp.New()
	go walletApp.StartServer(&fyneApp)

	fyneApp.Run()
}
