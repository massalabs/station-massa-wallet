package main

import (
	"embed"
	"fmt"
	"log"

	walletServer "github.com/massalabs/station-massa-wallet/cmd/massa-wallet"
	"github.com/massalabs/station-massa-wallet/internal/initialize"
	walletApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wails"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/logger"
)

//nolint:typecheck,nolintlint
//go:embed all:wails-frontend/dist
var wailsAssets embed.FS

func main() {
	app := walletApp.NewWalletApp()

	err := initialize.Logger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer logger.Close()

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
