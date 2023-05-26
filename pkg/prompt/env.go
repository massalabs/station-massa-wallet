package prompt

import (
	"os"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type envPrompter struct {
	PromptLocker
}

func (e *envPrompter) PromptRequest(req PromptRequest) {
	// create a go routine to send password when requested
	go func() {
		password := os.Getenv("WALLET_PASSWORD")
		switch req.Action {
		case walletapp.Delete, walletapp.Transfer, walletapp.Sign, walletapp.TradeRolls, walletapp.NewPassword, walletapp.Unprotect:
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

func (w *WalletPrompter) SaveAccountFile(nickname string) (string, error) {
	return runtime.SaveFileDialog(w.PromptApp.Ctx, runtime.SaveDialogOptions{
		Title:           "Backup Account File",
		DefaultFilename: wallet.Filename(nickname),
		Filters:         []runtime.FileFilter{{DisplayName: "Account File (*.yml)", Pattern: "*.yml"}},
	})
}

// Verifies at compilation time that walletPrompterMock implements WalletPrompterInterface interface.
var _ WalletPrompterInterface = &envPrompter{}
