package prompt

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"io"

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
	CorrelationId string
}

// WalletPrompter is a struct that wraps a Fyne GUI application and implements the delete.Confirmer interface.
type WalletPrompter struct {
	PromptLocker
}

const CORRELATION_ID_SIZE = 32

func generateCorrelationID() (string, error) {
	rand := cryptorand.Reader

	correlationId := make([]byte, CORRELATION_ID_SIZE)
	if _, err := io.ReadFull(rand, correlationId); err != nil {
		return "", err
	}

	return string(correlationId), nil
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

	correlationId, err := generateCorrelationID()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.CorrelationId = correlationId
	prompterApp.PromptRequest(req)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var output interface{} = nil

	for {
		select {
		case input := <-prompterApp.App().PromptInput:
			receivedCorrelationId := input.GetCorrelationID()
			// TODO: in test environnement we don't provide the correlation id for now
			// so here we accept to have an empty correlation id
			if receivedCorrelationId != "" && receivedCorrelationId != correlationId {
				logger.Warnf("received input with wrong correlation id")
				return nil, WrongCorrelationIdError(prompterApp)
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

				// continue
				if !keepListening {
					return nil, err
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
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: utils.ErrInvalidInputType.Error()})

	return utils.ErrInvalidInputType
}

func WrongCorrelationIdError(prompterApp WalletPrompterInterface) error {
	logger.Error(utils.ErrWrongPromptCorrelationId.Error())
	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: utils.ErrWrongPromptCorrelationId.Error()})

	return utils.ErrWrongPromptCorrelationId
}
