package main

import (
	"embed"
	"fmt"

	walletServer "github.com/massalabs/station-massa-wallet/cmd/massa-wallet"
	walletApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wails"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
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

		err := wallet.MigrateWallet()
		if err != nil {
			fmt.Println("can't migrate accounts: %w", err)
		}

		err = wailApp.Run()
		if err != nil {
			panic(err)
		}
	}
}
