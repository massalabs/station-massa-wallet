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
}

// WalletPrompter is a struct that wraps a Fyne GUI application and implements the delete.Confirmer interface.
type WalletPrompter struct {
	PromptLocker
}

func (w *WalletPrompter) PromptRequest(req PromptRequest) {
	runtime.EventsEmit(w.PromptApp.Ctx, walletapp.PromptRequestEvent, walletapp.PromptRequestData{Action: req.Action, Msg: req.Msg})
	w.PromptApp.Show()
}

func (w *WalletPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	runtime.EventsEmit(w.PromptApp.Ctx, eventId, data)
}

func (w *WalletPrompter) SelectBackupFilepath(nickname string) (string, error) {
	return runtime.SaveFileDialog(w.PromptApp.Ctx, runtime.SaveDialogOptions{
		Title:           "Backup Account File",
		DefaultFilename: wallet.Filename(nickname),
		Filters:         []runtime.FileFilter{{DisplayName: "Account File (*.yml)", Pattern: "*.yml"}},
	})
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

	var output interface{} = nil

	for {
		select {
		case input := <-prompterApp.App().PromptInput:

			var keepListening bool
			var err error

			switch req.Action {
			case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.Unprotect:
				output, keepListening, err = handlePasswordPrompt(prompterApp, input, wallet)
			case walletapp.NewPassword:
				output, keepListening, err = handleNewPasswordPrompt(prompterApp, input)
			case walletapp.Import:
				output, keepListening, err = handleImportPrompt(prompterApp, input)
			case walletapp.Backup:
				// If output is nil, it means that the user has not yet chosen a backup method.
				if output == nil {
					output, keepListening, err = handleBackupMethod(prompterApp, input)
				} else {
					output, keepListening, err = handlePasswordPrompt(prompterApp, input, wallet)
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
			return output, nil

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
	fmt.Println(InputTypeErr)
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: utils.ErrPromptInputType})

	return fmt.Errorf(InputTypeErr)
}
