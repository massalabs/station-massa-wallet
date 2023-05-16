package prompt

import (
	"os"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

type envPrompter struct {
	PromptLocker
}

func (e *envPrompter) PromptRequest(req PromptRequest) {
	// create a go routine to send password when requested
	go func() {
		password := os.Getenv("WALLET_PASSWORD")
		switch req.Action {
		case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.NewPassword, walletapp.Backup:
			e.PromptApp.PromptInput <- password
		}
	}()
}

func (e *envPrompter) EmitEvent(eventId string, data walletapp.EventData) {
	// unused in this implementation
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
