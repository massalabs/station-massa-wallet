package delete

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
)

// FynePrompter is a struct that wraps a Fyne GUI application and implements the delete.Confirmer interface.
type FynePrompter struct {
	guiApp *fyne.App
}

func (f *FynePrompter) Confirm(walletName string) (string, error) {
	result := <-PasswordDeleteDialog(walletName, f.guiApp)
	return result.Password, result.Err
}

// NewFynePrompter creates a new password prompter with the given Fyne GUI application.
func NewFynePrompter(f *fyne.App) *FynePrompter {
	return &FynePrompter{guiApp: f}
}

// Verifies at compilation time that FynePrompter implements Asker interface.
var _ Confirmer = &FynePrompter{}

// PasswordDeleteDialog displays a password dialog with the given nickname and ask confirmation to delete the account.
// It returns a channel to get what the user entered.
func PasswordDeleteDialog(nickname string, app *fyne.App) chan password.PasswordEntry {
	// Creates the result channel to listen in order to get the actual entry.
	result := make(chan password.PasswordEntry)

	window := (*app).NewWindow("Massa - Thyra")

	width := 250.0
	height := 80.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	passwordWidget := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Password", passwordWidget),
	}

	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			window.Hide()
			result <- password.PasswordEntry{Password: passwordWidget.Text, Err: nil}
		},
		OnCancel: func() {
			result <- password.PasswordEntry{Password: "", Err: errors.New("Confirm delete: cancelled by the user")}
			window.Hide()
		},
		SubmitText: "Delete",
		CancelText: "Cancel",
	}
	spacer := layout.NewSpacer()
	text1 := widget.NewLabel(`Delete "` + nickname + `" Wallet ?`)
	title := container.New(layout.NewHBoxLayout(), spacer, text1, spacer)
	text2 := widget.NewLabel("If you delete a wallet, you will lose your MAS associated to it and ")
	text3 := widget.NewLabel("won't be able to edit websites linked to this wallet anymore ")
	content := container.New(layout.NewVBoxLayout(), text2, text3, spacer)
	centeredForm := container.New(layout.NewVBoxLayout(), spacer, form, spacer)
	window.SetContent(container.New(layout.NewVBoxLayout(), title, spacer, content, spacer, centeredForm, spacer))
	window.CenterOnScreen()
	window.Canvas().Focus(passwordWidget)
	window.Show()

	return result
}
