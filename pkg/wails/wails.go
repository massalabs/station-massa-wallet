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
		Title:       "Massa Station Wallet",
		Width:       513,
		Height:      440,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.Startup,
		OnBeforeClose: app.BeforeClose,
		Bind: []interface{}{
			app,
		},
	})
}
