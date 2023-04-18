package main

import (
	"embed"

	walletServer "github.com/massalabs/thyra-plugin-wallet/cmd/massa-wallet"
	walletApp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wails"
)

//nolint:typecheck,nolintlint
//go:embed all:wails-frontend/dist
var wailsAssets embed.FS

func main() {
	walletApp := walletApp.NewWalletApp()

	wailApp := wails.NewWailsApp(walletApp, wailsAssets)
	go walletServer.StartServer(walletApp)

	err := wailApp.Run()
	if err != nil {
		panic(err)
	}
}
