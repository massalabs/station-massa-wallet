package wallet

import (
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

type walletPrompterMock struct {
	app           *walletapp.WalletApp
	resultChannel chan walletapp.EventData
}

func (w *walletPrompterMock) PromptRequest(req walletapp.PromptRequest, msg string, data interface{}) {
}

func (w *walletPrompterMock) EmitEvent(eventId string, data walletapp.EventData) {
	w.resultChannel <- data
}

func (w *walletPrompterMock) App() *walletapp.WalletApp {
	return w.app
}

// NewWalletPrompter creates a new password prompter with the given Fyne GUI application.
func NewWalletPrompterMock(app *walletapp.WalletApp, resultChannel chan walletapp.EventData) *walletPrompterMock {
	return &walletPrompterMock{
		app:           app,
		resultChannel: resultChannel,
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ wallet.WalletPrompterInterface = &walletPrompterMock{}
