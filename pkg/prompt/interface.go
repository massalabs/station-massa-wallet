package prompt

import walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"

// WalletPrompterInterface is used to mock the WalletPrompter struct in tests.
type WalletPrompterInterface interface {
	PromptRequest(req PromptRequest)
	EmitEvent(eventId string, data walletapp.EventData)
	App() *walletapp.WalletApp
	CtrlSink()
}
