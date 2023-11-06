package prompt

import (
	"os"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
)

type envPrompter struct {
	PromptLocker
}

func (e *envPrompter) PromptRequest(req PromptRequest) {
	// create a go routine to send password when requested
	go func() {
		password := os.Getenv("WALLET_PASSWORD")

		switch req.Action {
		case walletapp.Sign, walletapp.TradeRolls:
			e.PromptApp.PromptInput <- walletapp.SignPromptInput{Password: password, Fees: "500"}
		case walletapp.Delete, walletapp.NewPassword, walletapp.Unprotect:
			e.PromptApp.PromptInput <- password
		}
	}()
}

func (e *envPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	// unused in this implementation
}

func (w *envPrompter) SelectBackupFilepath(nickname string) (string, error) {
	return os.Getenv("WALLET_BACKUP_FILEPATH"), nil
}

func NewEnvPrompter(app *walletapp.WalletApp) *envPrompter {
	return &envPrompter{
		PromptLocker: PromptLocker{
			PromptApp: app,
		},
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &envPrompter{}
