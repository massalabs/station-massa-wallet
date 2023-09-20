package app

import (
	"log"
	"os"

	"github.com/bluele/gcache"
	"github.com/massalabs/station-massa-hello-world/pkg/plugin"
	"github.com/massalabs/station-massa-wallet/api/server/restapi"
	"github.com/massalabs/station-massa-wallet/internal/handler"
	"github.com/massalabs/station-massa-wallet/internal/initialize"
	walletApp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station/pkg/logger"
)

func StartServer(app *walletApp.WalletApp) {
	// Initialize cache
	gc := gcache.New(20).
		LRU().
		Build()

	massaClient := network.NewNodeFetcher()

	var promptApp prompt.WalletPrompterInterface = prompt.NewWalletPrompter(app)
	if walletApp.IsTestMode() {
		promptApp = prompt.NewEnvPrompter(app)
	}

	err := initialize.Logger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	AssetsStore, err := assets.NewAssetsStore()
	if err != nil {
		logger.Fatalf("Failed to create AssetsStore: %v", err)
	}

	// Initializes API
	massaWalletAPI, err := handler.InitializeAPI(
		promptApp,
		massaClient,
		AssetsStore,
		gc,
	)
	if err != nil {
		logger.Fatalf("Failed to initialize API: %v", err)
	}

	// instantiates and configure server
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	if os.Getenv("STANDALONE") == "1" {
		server.Port = 8080
	}

	listener, err := server.HTTPListener()
	if err != nil {
		logger.Fatalf("Failed to create HTTP listener: %v", err)
	}

	plugin.RegisterPlugin(listener)

	if err := server.Serve(); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
