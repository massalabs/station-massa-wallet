package main

import (
	"context"
	"embed"
	"fmt"
	"time"

	walletApp "github.com/massalabs/thyra-plugin-wallet/cmd/massa-wallet"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//nolint:typecheck,nolintlint
//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	fmt.Println("Greet input:" + name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func main() {

	go walletApp.StartServer()

	startChan := make(chan bool)

	go waitAndStart(startChan)

	for {
		select {
		case <-startChan:
			fmt.Println("startChan received")
			app := NewApp()
			wailApp := application.NewWithOptions(&options.App{
				Title:             "wallet-prompt",
				Width:             1024,
				Height:            768,
				StartHidden:       true,
				HideWindowOnClose: true,
				AssetServer: &assetserver.Options{
					Assets: assets,
				},
				BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
				OnStartup:        app.startup,
				Bind: []interface{}{
					app,
				},
			})

			go waitAndStop(app)
			fmt.Println("Run")
			wailApp.Run()
			return
		}
	}

}

func waitAndStart(start chan bool) {
	fmt.Println("waitAndStart")

	time.Sleep(time.Second * 3)
	start <- true
}

func waitAndStop(app *App) {

	time.Sleep(time.Second * 3)

	fmt.Println("Show!")
	runtime.Show(app.ctx)

	time.Sleep(time.Second * 3)

	fmt.Println("Hide")
	runtime.Hide(app.ctx)

	runtime.Quit(app.ctx)

}
