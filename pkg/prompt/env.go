package prompt

import (
	"os"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

type envPrompter struct {
	app *walletapp.WalletApp
}

func (e *envPrompter) PromptRequest(req PromptRequest) {
	// create a go routine to send password when requested
	go func() {
		password := os.Getenv("WALLET_PASSWORD")
		switch req.Action {
		case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.NewPassword, walletapp.Export:
			e.app.PromptInput <- password
		}
	}()
}

func (e *envPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	// unused in this implementation
}

func (e *envPrompter) App() *walletapp.WalletApp {
	return e.app
}

func (w *envPrompter) Lock() {
	w.app.IsListening = true
}

func (w *envPrompter) Unlock() {
	w.app.IsListening = false
}

func (w *envPrompter) IsListening() bool {
	return w.app.IsListening
}

func NewEnvPrompter(app *walletapp.WalletApp) *envPrompter {
	return &envPrompter{
		app: app,
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &envPrompter{}
