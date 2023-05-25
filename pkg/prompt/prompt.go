package prompt

import (
	"context"
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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
	PromptLocker
}

func (w *WalletPrompter) PromptRequest(req PromptRequest) {
	runtime.EventsEmit(w.PromptApp.Ctx, walletapp.PromptRequestEvent, walletapp.PromptRequestData{Action: req.Action, Msg: req.Msg, Data: req.Data})
	w.PromptApp.Show()
}

func (w *WalletPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	runtime.EventsEmit(w.PromptApp.Ctx, eventId, data)
}

func NewWalletPrompter(app *walletapp.WalletApp) *WalletPrompter {
	return &WalletPrompter{
		PromptLocker: PromptLocker{
			PromptApp: app,
		},
	}
}

// Verifies at compilation time that WalletPrompter implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &WalletPrompter{}

func WakeUpPrompt(
	prompterApp WalletPrompterInterface,
	req PromptRequest,
	wallet *wallet.Wallet,
) (interface{}, error) {
	if prompterApp.IsListening() {
		fmt.Println(AlreadyListeningErr)
		return nil, fmt.Errorf(AlreadyListeningErr)
	}
	prompterApp.Lock()
	defer prompterApp.Unlock()

	prompterApp.PromptRequest(req)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var isPrivateKeyBackup *bool = nil

	for {
		select {
		case input := <-prompterApp.App().PromptInput:

			var keepListening bool
			var res interface{} = nil
			var err error

			switch req.Action {
			case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.Uncipher:
				res, keepListening, err = handlePasswordPrompt(prompterApp, input, wallet)
			case walletapp.NewPassword:
				res, keepListening, err = handleNewPasswordPrompt(prompterApp, input)
			case walletapp.Import:
				res, keepListening, err = handleImportPrompt(prompterApp, input)
			case walletapp.Backup:
				if isPrivateKeyBackup == nil {
					keepListening, err = handleBackupMethod(prompterApp, input)
					isPrivateKeyBackup = &keepListening
				} else {
					res, keepListening, err = handlePasswordPrompt(prompterApp, input, wallet)
				}
			}

			if err != nil {
				fmt.Println(err)
				if !keepListening {
					return nil, err
				}
			}
			if keepListening {
				continue
			}
			return res, nil

		case <-prompterApp.App().CtrlChan:
			fmt.Println(ActionCanceledErr)
			return nil, fmt.Errorf(ActionCanceledErr)

		case <-ctxTimeout.Done():
			fmt.Println(TimeoutErr)
			prompterApp.EmitEvent(walletapp.PromptResultEvent,
				walletapp.EventData{Success: false, CodeMessage: utils.ErrTimeout})

			return nil, fmt.Errorf(TimeoutErr)
		}
	}
}

func InputTypeError(prompterApp WalletPrompterInterface) error {
	fmt.Println("invalid prompt input type")
	// TODO: upgrade CodeMessage
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: InputTypeErr})

	return fmt.Errorf(InputTypeErr)
}

func UserChoiceError(prompterApp WalletPrompterInterface) error {
	fmt.Println("invalid user choice input")
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: UserChoiceErr})

	return fmt.Errorf(UserChoiceErr)
}
