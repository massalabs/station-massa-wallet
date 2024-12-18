package walletapp

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WalletApp struct {
	Ctx         context.Context
	CtrlChan    chan PromptCtrl
	PromptInput chan EventInterface
	Wallet      *wallet.Wallet
	Shutdown    bool
	IsListening bool
}

func (a *WalletApp) cleanExit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if a.Ctx != nil {
		a.Shutdown = true
		runtime.Quit(a.Ctx)
	}
}

func NewWalletApp(wallet *wallet.Wallet) *WalletApp {
	app := &WalletApp{
		CtrlChan:    make(chan PromptCtrl),
		PromptInput: make(chan EventInterface),
		Wallet:      wallet,
		Shutdown:    false,
		IsListening: false,
	}

	go app.cleanExit()

	return app
}

func IsTestMode() bool {
	return os.Getenv("WALLET_PASSWORD") != ""
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *WalletApp) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *WalletApp) BeforeClose(ctx context.Context) bool {
	a.Hide()

	if a.Shutdown {
		return false
	}

	// Send a cancel message to the prompt and do NOT shutdown
	if a.IsListening {
		logger.Warn("canceling prompt before closing")
		a.CtrlChan <- Cancel
	}

	runtime.WindowReloadApp(a.Ctx)

	return true
}

// Functions to send user input to the backend from the wails frontend

func (a *WalletApp) SendPromptInput(input string) {
	if !a.IsListening {
		logger.Warn("Not listening (in SendPromptInput)")
		return
	}

	a.PromptInput <- &StringPromptInput{BaseMessage: BaseMessage{}, Message: input}
}

func (a *WalletApp) SendSignPromptInput(password string, fees string) {
	if !a.IsListening {
		logger.Warn("Not listening (in SendSignPromptInput)")
		return
	}

	a.PromptInput <- &SignPromptInput{BaseMessage: BaseMessage{}, Password: password, Fees: fees}
}

func (a *WalletApp) SendPKeyPromptInput(privateKeyText string, nickname string, password string) {
	if !a.IsListening {
		logger.Warn("Not listening (in SendPKeyPromptInput)")
		return
	}

	guardedPrivateKey := memguard.NewBufferFromBytes([]byte(privateKeyText))

	a.PromptInput <- &ImportPKeyPromptInput{
		BaseMessage: BaseMessage{},
		PrivateKey:  guardedPrivateKey,
		Nickname:    nickname,
		Password:    memguard.NewBufferFromBytes([]byte(password)),
	}
}

// Function called by the wails frontend

// AbortAction sends a cancel message to the prompt
func (a *WalletApp) AbortAction() {
	if a.IsListening {
		logger.Warn("Abort action")
		a.CtrlChan <- Cancel
	}
}

func (a *WalletApp) Show() {
	runtime.WindowShow(a.Ctx)
}

func (a *WalletApp) Hide() {
	runtime.WindowHide(a.Ctx)
}

type selectFileResult struct {
	Err         string `json:"err"`
	CodeMessage string `json:"codeMessage"`
	FilePath    string `json:"filePath"`
	Nickname    string `json:"nickname"`
}

func (a *WalletApp) SelectAccountFile() selectFileResult {
	filePath, err := runtime.OpenFileDialog(a.Ctx, runtime.OpenDialogOptions{
		Title:   "Import account file",
		Filters: []runtime.FileFilter{{DisplayName: "Account File (*.yaml;*.yml)", Pattern: "*.yaml;*.yml"}},
	})
	if err != nil {
		logger.Errorf("error while selecting a file: %v", err)
		return selectFileResult{Err: err.Error(), CodeMessage: utils.ErrAccountFile}
	}

	if filePath == "" {
		return selectFileResult{Err: utils.ErrNoFile, CodeMessage: utils.ErrNoFile}
	}

	acc, err := a.Wallet.Load(filePath)
	if err != nil {
		logger.Errorf("error while loading file: %v", err)
		return selectFileResult{Err: err.Error(), CodeMessage: utils.WailsErrorCode(err)}
	}

	return selectFileResult{FilePath: filePath, Nickname: acc.Nickname}
}

func (a *WalletApp) IsNicknameUnique(nickname string) bool {
	return a.Wallet.NicknameIsUnique(nickname) != nil
}

func (a *WalletApp) IsNicknameValid(nickname string) bool {
	return account.NicknameIsValid(nickname)
}
