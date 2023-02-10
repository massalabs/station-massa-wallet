package privateKey

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// FynePrompter is a struct that wraps a Fyne GUI application and implements the password.PasswordAsker interface.
type FynePrompter struct {
	guiApp *fyne.App
}

type PrivateKeyEntry struct {
	PrivateKey string
	Err        error
}

func (f *FynePrompter) Ask() (string, error) {
	result := <-PrivateKeyDialog(f.guiApp)

	return result.PrivateKey, result.Err
}

// NewFynePrompter creates a new privateKey prompter with the given Fyne GUI application.
func NewFynePrompter(f *fyne.App) *FynePrompter {
	return &FynePrompter{guiApp: f}
}

func PrivateKeyDialog(app *fyne.App) chan PrivateKeyEntry {
	result := make(chan PrivateKeyEntry)

	window := (*app).NewWindow("Massa - Thyra")

	width := 700.0
	height := 100.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	privateKey := widget.NewPasswordEntry()
	items := []*widget.FormItem{

		widget.NewFormItem("Private key", privateKey),
	}

	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			var err error
			if privateKey.Text == "" {
				err = fmt.Errorf("Private key is required")
			}
			if err != nil {
				dialog.ShowError(err, window)
			} else {

				window.Hide()
				result <- PrivateKeyEntry{

					PrivateKey: privateKey.Text,
					Err:        nil,
				}
			}
		},
		OnCancel: func() {
			result <- PrivateKeyEntry{
				PrivateKey: "",
				Err:        errors.New("wallet loading cancelled by the user"),
			}
			window.Hide()
		},
		SubmitText: "Load",
		CancelText: "Cancel",
	}
	spacer := layout.NewSpacer()
	text := widget.NewLabel(`Load a Wallet`)
	title := container.New(layout.NewHBoxLayout(), spacer, text, spacer)
	centeredForm := container.New(layout.NewVBoxLayout(), spacer, form, spacer)
	window.SetContent(container.New(layout.NewVBoxLayout(), title, spacer, centeredForm, spacer))
	window.CenterOnScreen()
	window.Canvas().Focus(privateKey)
	window.Show()

	return result
}
