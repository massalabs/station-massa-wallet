package wails

import (
	"embed"

	wApp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func NewWailsApp(app *wApp.WalletApp, assets embed.FS) *application.Application {
	// Create a new Wails application
	return application.NewWithOptions(&options.App{
		Title:       "wallet-prompt",
		Width:       1024,
		Height:      768,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnBeforeClose:    app.BeforeClose,
		Bind: []interface{}{
			app,
		},
	})
}
