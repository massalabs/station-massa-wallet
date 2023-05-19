package prompt

import (
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PASSWORD_MIN_LENGTH = 5

type PromptRequestData struct {
	Msg  string
	Data interface{}
}

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

// CtrlSink is a blocking function that waits for the cancel msg sended when the wails prompt is closed.
func (w *WalletPrompter) CtrlSink() {
	<-w.app.CtrlChan
}

func NewWalletPrompter(app *walletapp.WalletApp) *WalletPrompter {
	return &WalletPrompter{app: app}
}

// Verifies at compilation time that WalletPrompter implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &WalletPrompter{}
