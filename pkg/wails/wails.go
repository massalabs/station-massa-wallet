package wails

import (
	"embed"

	wApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func NewWailsApp(app *wApp.WalletApp, assets embed.FS) *application.Application {
	// Create a new Wails application
	return application.NewWithOptions(&options.App{
		Title:       "MassaStation Wallet",
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
		EnumBind: []interface{}{
			wApp.PromptRequest,
			wApp.EventTypes,
			config.RuleTypes,
		},
	})
}
