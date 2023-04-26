package wallet

import (
	"context"
	"fmt"
	"strings"
	"time"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PASSWORD_MIN_LENGTH = 5

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

type PromptRequestData struct {
	Msg  string
	Data interface{}
}

func (w *Wallet) PromptPassword(
	prompterApp WalletPrompterInterface,
	action walletapp.PromptRequest,
	data *PromptRequestData,
) (string, error) {
	prompterApp.PromptRequest(action, data.Msg, data.Data)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case password := <-prompterApp.App().PasswordChan:
			err := w.Unprotect(password)
			if err != nil {
				errStr := "error unprotecting wallet:" + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			return password, nil
		case <-prompterApp.App().CtrlChan:
			msg := "Action canceled by user"
			fmt.Println(msg)
			return "", fmt.Errorf(msg)
		case <-ctxTimeout.Done():
			errStr := "Password prompt reached timeout"
			fmt.Println(errStr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: errStr, Error: "timeoutError"})
			return "", fmt.Errorf(errStr)
		}
	}
}

func PromptCreatePassword(
	prompterApp WalletPrompterInterface,
	nickname string,
) (string, error) {
	data := &PromptRequestData{
		Msg:  fmt.Sprintf("Creating new password for account %s", nickname),
		Data: nil,
	}
	prompterApp.PromptRequest(walletapp.NewPassword, data.Msg, data.Data)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case password := <-prompterApp.App().PasswordChan:
			trimmedPassword := strings.TrimSpace(password)
			if len(trimmedPassword) < PASSWORD_MIN_LENGTH {
				// TODO implement password strength check
				errStr := fmt.Sprintf("password must %d minimum length", PASSWORD_MIN_LENGTH)
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			return trimmedPassword, nil

		case <-prompterApp.App().CtrlChan:
			msg := "Action canceled by user"
			fmt.Println(msg)
			return "", fmt.Errorf(msg)
		case <-ctxTimeout.Done():
			errStr := "Password prompt reached timeout"
			fmt.Println(errStr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: errStr, Error: "timeoutError"})
			return "", fmt.Errorf(errStr)
		}
	}
}

// Verifies at compilation time that WalletPrompter implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &WalletPrompter{}
