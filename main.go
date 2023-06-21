package main

import (
	"embed"
	"fmt"

	walletServer "github.com/massalabs/station-massa-wallet/cmd/massa-wallet"
	walletApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wails"
)

//nolint:typecheck,nolintlint
//go:embed all:wails-frontend/dist
var wailsAssets embed.FS

func main() {
	app := walletApp.NewWalletApp()

	if walletApp.IsTestMode() {
		fmt.Println("Wallet is running in test mode")
		walletServer.StartServer(app)
	} else {
		wailApp := wails.NewWailsApp(app, wailsAssets)

		go walletServer.StartServer(app)

		err := wailApp.Run()
		if err != nil {
			panic(err)
		}
	}
}
