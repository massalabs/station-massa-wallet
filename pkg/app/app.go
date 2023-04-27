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

// App struct
type WalletApp struct {
	Ctx            context.Context
	CtrlChan       chan PromptCtrl
	PasswordChan   chan string
	PrivateKeyChan chan string
	WalletFileChan chan string
	Shutdown       bool
}

func (a *WalletApp) cleanExit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Shutdown = true
	runtime.Quit(a.Ctx)
}

func NewWalletApp() *WalletApp {
	app := &WalletApp{
		CtrlChan:       make(chan PromptCtrl),
		PasswordChan:   make(chan string),
		PrivateKeyChan: make(chan string),
		WalletFileChan: make(chan string),
		Shutdown:       false,
	}
	go app.cleanExit()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *WalletApp) Startup(ctx context.Context) {
	a.Ctx = ctx
	// Hide the window on startup. This is a workaround for a Wails bug on Linux
	// see: https://github.com/wailsapp/wails/issues/2605
	a.Hide()
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
	fmt.Println("Received password input!")
	a.PasswordChan <- password
}

// AbortAction is binded to the frontend
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

// App struct
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
	wallet, err := wallet.LoadFile(filePath)
	if err != nil {
		return selectFileResult{Err: err.Error()}
	}

	return selectFileResult{FilePath: filePath, Nickname: wallet.Nickname}
}

func (a *WalletApp) ImportWalletFile(filePath string) {
	a.WalletFileChan <- filePath
}

func (a *WalletApp) ImportPrivateKey(pkey string) {
	fmt.Println("Received password input!")
	a.PrivateKeyChan <- pkey
}
