package walletapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WalletApp struct {
	Ctx            context.Context
	CtrlChan       chan PromptCtrl
	PasswordChan   chan string
	PrivateKeyChan chan ImportFromPKey
	WalletFileChan chan string
	Shutdown       bool
}

func (a *WalletApp) cleanExit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Shutdown = true
	if a.Ctx != nil {
		runtime.Quit(a.Ctx)
	}
}

func NewWalletApp() *WalletApp {
	app := &WalletApp{
		CtrlChan:       make(chan PromptCtrl),
		PasswordChan:   make(chan string),
		PrivateKeyChan: make(chan ImportFromPKey),
		WalletFileChan: make(chan string),
		Shutdown:       false,
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
	a.CtrlChan <- Cancel
	runtime.WindowReloadApp(a.Ctx)

	return true
}

// ApplyPassword is binded to the frontend
func (a *WalletApp) ApplyPassword(password string) {
	a.PasswordChan <- password
}

// AbortAction is bound to the frontend
// It sends a cancel message to the prompt
func (a *WalletApp) AbortAction() {
	fmt.Println("Abort action")
	a.CtrlChan <- Cancel
}

func (a *WalletApp) Show() {
	runtime.WindowShow(a.Ctx)
}

func (a *WalletApp) Hide() {
	runtime.WindowHide(a.Ctx)
}

type selectFileResult struct {
	Err      string `json:"err"`
	FilePath string `json:"filePath"`
	Nickname string `json:"nickname"`
}

func (a *WalletApp) SelectAccountFile() selectFileResult {
	filePath, err := runtime.OpenFileDialog(a.Ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return selectFileResult{Err: err.Error()}
	}
	wallet, loadErr := wallet.LoadFile(filePath)
	if loadErr != nil {
		return selectFileResult{Err: loadErr.Err.Error()}
	}

	return selectFileResult{FilePath: filePath, Nickname: wallet.Nickname}
}

func (a *WalletApp) ImportWalletFile(filePath string) {
	a.WalletFileChan <- filePath
}

type ImportFromPKey struct {
	PrivateKey string
	Password   string
	Nickname   string
}

func (a *WalletApp) ImportPrivateKey(pkey string, nickname string, password string) {
	a.PrivateKeyChan <- ImportFromPKey{PrivateKey: pkey, Nickname: nickname, Password: password}
}
