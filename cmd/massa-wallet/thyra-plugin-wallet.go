package app

import (
	"log"
	"os"

	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-hello-world/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/internal/handler"
	walletApp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	constants "github.com/massalabs/thyra-plugin-wallet/pkg/plugin"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
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

	// Initializes API
	massaWalletAPI, err := handler.InitializeAPI(
		promptApp,
		massaClient,
		gc,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// instantiates and configure server
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	if os.Getenv("STANDALONE") == "1" {
		server.Port = 8080
	}

	listener, err := server.HTTPListener()
	if err != nil {
		log.Fatalln(err)
	}

	plugin.RegisterPlugin(listener, plugin.Info{
		Name: constants.PluginName, Author: constants.PluginAuthor,
		Description: constants.PluginDescription, APISpec: "", Logo: "web/wallet.svg",
	})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
