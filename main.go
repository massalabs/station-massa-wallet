package main

import (
	"embed"
	"log"

	walletServer "github.com/massalabs/station-massa-wallet/cmd/massa-wallet"
	"github.com/massalabs/station-massa-wallet/internal/initialize"
	walletApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wails"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
	"github.com/massalabs/station/pkg/logger"
)

//nolint:typecheck,nolintlint
//go:embed all:wails-frontend/dist
var wailsAssets embed.FS

func main() {
	err := initialize.Logger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer logger.Close()

	wallet, err := walletmanager.New("")
	if err != nil {
		logger.Fatalf("Failed to initialize wallet: %v", err)
	}

	app := walletApp.NewWalletApp(wallet)

	if walletApp.IsTestMode() {
		logger.Info("Wallet is running in test mode")
		walletServer.StartServer(app)
	} else {
		wailApp := wails.NewWailsApp(app, wailsAssets)

		go walletServer.StartServer(app)

		err = wailApp.Run()
		if err != nil {
			logger.Fatalf("Failed to run Wails app: %v", err)
		}
	}
}
