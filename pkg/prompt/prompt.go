package prompt

import (
	"context"
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PASSWORD_MIN_LENGTH = 5

type PromptRequest struct {
	Action walletapp.PromptRequestAction
	Msg    string
	Data   interface{}
}

// WalletPrompter is a struct that wraps a Fyne GUI application and implements the delete.Confirmer interface.
type WalletPrompter struct {
	app *walletapp.WalletApp
}

func (w *WalletPrompter) PromptRequest(req PromptRequest) {
	runtime.EventsEmit(w.app.Ctx, walletapp.PromptRequestEvent, walletapp.PromptRequestData{Action: req.Action, Msg: req.Msg, Data: req.Data})
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

func WakeUpPrompt(
	prompterApp WalletPrompterInterface,
	req PromptRequest,
	wallet *wallet.Wallet,
) (interface{}, error) {
	prompterApp.PromptRequest(req)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	for {
		select {
		case input := <-prompterApp.App().PromptInput:

			var keepListening bool
			var res interface{}
			var err error

			switch req.Action {
			case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.Export:
				res, keepListening, err = handlePasswordPrompt(prompterApp, input, wallet)
			case walletapp.NewPassword:
				res, keepListening, err = handleNewPasswordPrompt(prompterApp, input)
			case walletapp.Import:
				res, keepListening, err = handleImportPrompt(prompterApp, input)
			}

			if err != nil {
				fmt.Println(err)
				if keepListening {
					continue
				}
				return nil, err
			}
			return res, nil

		case <-prompterApp.App().CtrlChan:
			fmt.Println(ActionCanceledErr)
			return nil, fmt.Errorf(ActionCanceledErr)

		case <-ctxTimeout.Done():
			fmt.Println(TimeoutErr)
			prompterApp.EmitEvent(walletapp.PromptResultEvent,
				walletapp.EventData{Success: false, Data: TimeoutErr, Error: "timeoutError"})

			go prompterApp.CtrlSink()

			return nil, fmt.Errorf(TimeoutErr)
		}
	}
}

func InputTypeError(prompterApp WalletPrompterInterface) error {
	fmt.Println("invalid prompt input type")
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, Data: InputTypeErr, Error: "InputTypeErr"})

	go prompterApp.CtrlSink()

	return fmt.Errorf(InputTypeErr)
}
