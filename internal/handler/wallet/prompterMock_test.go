package wallet

import (
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
)

const WalletBackupFilepath = "walletBackup.yaml"

type walletPrompterMock struct {
	resultChannel chan walletapp.EventData
	prompt.PromptLocker
}

func (w *walletPrompterMock) PromptRequest(req prompt.PromptRequest) {
}

func (w *walletPrompterMock) EmitEvent(eventId string, data walletapp.EventData) {
	w.resultChannel <- data
}

func (w *walletPrompterMock) SelectBackupFilepath(nickname string) (string, error) {
	return WalletBackupFilepath, nil
}

// NewWalletPrompterMock creates a new password prompter with the given Fyne GUI application.
func NewWalletPrompterMock(app *walletapp.WalletApp, resultChannel chan walletapp.EventData) *walletPrompterMock {
	return &walletPrompterMock{
		resultChannel: resultChannel,
		PromptLocker: prompt.PromptLocker{
			PromptApp: app,
		},
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ prompt.WalletPrompterInterface = &walletPrompterMock{}
