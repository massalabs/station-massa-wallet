package wallet

import (
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

type walletPrompterMock struct {
	app           *walletapp.WalletApp
	resultChannel chan walletapp.EventData
}

func (w *walletPrompterMock) PromptRequest(req prompt.PromptRequest) {
}

func (w *walletPrompterMock) EmitEvent(eventId string, data walletapp.EventData) {
	w.resultChannel <- data
}

func (w *walletPrompterMock) App() *walletapp.WalletApp {
	return w.app
}

func (w *walletPrompterMock) Lock() {
	w.app.IsListening = true
}

func (w *walletPrompterMock) Unlock() {
	w.app.IsListening = false
}

func (w *walletPrompterMock) IsListening() bool {
	return w.app.IsListening
}

// NewWalletPrompter creates a new password prompter with the given Fyne GUI application.
func NewWalletPrompterMock(app *walletapp.WalletApp, resultChannel chan walletapp.EventData) *walletPrompterMock {
	return &walletPrompterMock{
		app:           app,
		resultChannel: resultChannel,
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ prompt.WalletPrompterInterface = &walletPrompterMock{}
