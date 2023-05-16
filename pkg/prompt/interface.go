package prompt

import walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"

// WalletPrompterInterface is used to mock the WalletPrompter struct in tests.
type WalletPrompterInterface interface {
	PromptRequest(req walletapp.PromptRequest, msg string, data interface{})
	EmitEvent(eventId string, data walletapp.EventData)
	App() *walletapp.WalletApp
}
