package wallet

import (
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WalletPrompter is a struct that wraps a Fyne GUI application and implements the delete.Confirmer interface.
type WalletPrompter struct {
	app *walletapp.WalletApp
}

func (w *WalletPrompter) PromptRequest(req walletapp.PromptRequest, msg string, data interface{}) {
	runtime.EventsEmit(w.app.Ctx, walletapp.PromptRequesEvent, walletapp.PromptRequestData{Action: req, Msg: msg, Data: data})
	w.app.Show()
}

func (w *WalletPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	runtime.EventsEmit(w.app.Ctx, eventId, data)
}

func (w *WalletPrompter) App() *walletapp.WalletApp {
	return w.app
}

// NewWalletPrompter creates a new password prompter with the given Fyne GUI application.
func NewWalletPrompter(app *walletapp.WalletApp) *WalletPrompter {
	return &WalletPrompter{app: app}
}

// Confirmer is the delete account confirm interface
type WalletPrompterInterface interface {
	PromptRequest(req walletapp.PromptRequest, msg string, data interface{})
	EmitEvent(eventId string, data walletapp.EventData)
	App() *walletapp.WalletApp
}

// Verifies at compilation time that WalletPrompter implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &WalletPrompter{}
