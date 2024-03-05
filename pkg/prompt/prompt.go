package prompt

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PromptRequest struct {
	Action        walletapp.PromptRequestAction
	Msg           string
	Data          interface{}
	CodeMessage   string
	CorrelationID string
}

// WalletPrompter is a struct that wraps a Wails GUI application and implements the WalletPrompterInterface interface.
type WalletPrompter struct {
	PromptLocker
}

func NewWalletPrompter(app *walletapp.WalletApp) *WalletPrompter {
	return &WalletPrompter{
		PromptLocker: PromptLocker{
			PromptApp: app,
		},
	}
}

func (w *WalletPrompter) PromptRequest(req PromptRequest) {
	runtime.EventsEmit(w.PromptApp.Ctx, walletapp.PromptRequestEvent, req)
	w.PromptApp.Show()
}

func (w *WalletPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	runtime.EventsEmit(w.PromptApp.Ctx, eventId, data)
}

func (w *WalletPrompter) SelectBackupFilepath(nickname string) (string, error) {
	return runtime.SaveFileDialog(w.PromptApp.Ctx, runtime.SaveDialogOptions{
		Title:           "Backup Account File",
		DefaultFilename: wallet.Filename(nickname),
		Filters:         []runtime.FileFilter{{DisplayName: "Account File (*.yaml)", Pattern: "*.yaml"}},
	})
}

// Verifies at compilation time that WalletPrompter implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &WalletPrompter{}

func WakeUpPrompt(
	prompterApp WalletPrompterInterface,
	req PromptRequest,
	acc *account.Account,
) (interface{}, error) {
	if prompterApp.IsListening() {
		logger.Warn(AlreadyListeningErr)
		return nil, fmt.Errorf(AlreadyListeningErr)
	}

	prompterApp.Lock()
	defer prompterApp.Unlock()

	correlationId := strconv.FormatUint(rand.Uint64(), 10)
	req.CorrelationID = correlationId
	prompterApp.PromptRequest(req)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var output interface{} = nil

	for {
		select {
		case input := <-prompterApp.App().PromptInput:
			receivedCorrelationId := input.GetCorrelationID()
			// In test environnement we can't provide the correlation id,
			// so here we continue only if the correlation id is not the same as the one we sent,
			// and if the correlation id is not 1 (which is the test value for correlation id).
			if receivedCorrelationId != "1" && receivedCorrelationId != correlationId {
				logger.Warn("Received a prompt input with a different correlation id")
				continue
			}

			var keepListening bool
			var err error

			switch req.Action {
			case walletapp.Delete, walletapp.Unprotect:
				output, keepListening, err = handlePasswordPrompt(prompterApp, input, acc)
			case walletapp.Sign:
				output, keepListening, err = handleSignPrompt(prompterApp, input, acc)
			case walletapp.NewPassword:
				output, keepListening, err = handleNewPasswordPrompt(prompterApp, input)
			case walletapp.Import:
				output, keepListening, err = handleImportPrompt(prompterApp, input)
			case walletapp.Backup:
				// If output is nil, it means that the user has not yet chosen a backup method.
				if output == nil {
					output, keepListening, err = handleBackupMethod(prompterApp, input)
				} else {
					output, keepListening, err = handlePasswordPrompt(prompterApp, input, acc)
				}
			}

			if err != nil {
				logger.Error(err)

				if !keepListening {
					return nil, err
				} else {
					logger.Warn("Handling the user prompt input failed, keep listening for another input...")
					output = nil
				}
			}

			if keepListening {
				continue
			}

			return output, nil

		case <-prompterApp.App().CtrlChan:
			logger.Warn(utils.ErrActionCanceled.Error())

			return nil, utils.ErrActionCanceled

		case <-ctxTimeout.Done():
			logger.Warn(utils.ErrTimeout.Error())
			prompterApp.EmitEvent(walletapp.PromptResultEvent,
				walletapp.EventData{Success: false, CodeMessage: utils.ErrTimeoutMsg})

			return nil, utils.ErrTimeout
		}
	}
}

func InputTypeError(prompterApp WalletPrompterInterface) error {
	logger.Error(utils.ErrInvalidInputType.Error())
	return utils.ErrInvalidInputType
}
