package prompt

import (
	"os"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

type envPrompter struct {
	app *walletapp.WalletApp
}

func (e *envPrompter) PromptRequest(req walletapp.PromptRequest, msg string, data interface{}) {
	// create a go routine to send password when requested
	go func() {
		password := os.Getenv("WALLET_PASSWORD")
		if req == walletapp.Password || req == walletapp.NewPassword || req == walletapp.Sign || req == walletapp.Transfer {
			e.app.PasswordChan <- password
		}
	}()
}

func (e *envPrompter) EmitEvent(eventId string, data walletapp.EventData) {
}

func (e *envPrompter) App() *walletapp.WalletApp {
	return e.app
}

func NewEnvPrompter(app *walletapp.WalletApp) *envPrompter {
	return &envPrompter{
		app: app,
	}
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &envPrompter{}
