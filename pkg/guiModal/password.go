package guiModal

// In this file is implemented the Fyne version of the guiModal.PasswordAsker interface.
// As the Fyne application must have a GUI application that is initialized from the main,
// we use the FynePrompter structure to wrap it and thus allow the Ask function to use it thanks to the closure mechanism.
// Finally, the end user input is retrieved asynchronously using a channel.

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// FynePrompter is a struct that wraps a Fyne GUI application and implements the guiModal.PasswordAsker interface.
type FynePrompter struct {
	guiApp *fyne.App
}

// Ask displays a password dialog using the given nickname.
//
// It returns the entered password and any error that may have occurred.
// Note: Reuses the Fyne wrapped application.
func (f *FynePrompter) Ask(name string) (string, error) {
	result := <-PasswordDialog(name, f.guiApp)
	return result.password, result.err
}

// NewFynePrompter creates a new password prompter with the given Fyne GUI application.
func NewFynePrompter(f *fyne.App) *FynePrompter {
	return &FynePrompter{guiApp: f}
}

// Verifies at compilation time that FynePrompter implements Asker interface.
var _ PasswordAsker = &FynePrompter{}

// passwordEntry represents a password entry, containing the password and any error that may have occurred.
// Data sent through the channel to get the password entry asynchronously.
type passwordEntry struct {
	password string
	err      error
}

// PasswordDialog displays a password dialog with the given nickname.
// It returns a channel to get what the user entered.
func PasswordDialog(nickname string, app *fyne.App) chan passwordEntry {
	// Creates the password dialog window
	window := (*app).NewWindow("Massa - Thyra plugin - Wallet")
	width := 250.0
	height := 90.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	// Creates the password widget
	passwordWidget := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Password", passwordWidget),
	}

	// Creates the result channel to listen to to get the actual entry.
	result := make(chan passwordEntry)

	// Creates a simple form with two buttons: submit and cancel.
	// Actual result are sent via the result channel.
	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			window.Hide()
			result <- passwordEntry{password: passwordWidget.Text, err: nil}
		},
		OnCancel: func() {
			window.Hide()
			result <- passwordEntry{password: "", err: errors.New("password entry: cancelled by the user")}
		},
		SubmitText: "Submit",
		CancelText: "Cancel",
	}

	// Fills the window element with created content
	window.SetContent(form)
	window.CenterOnScreen()
	window.Canvas().Focus(passwordWidget)
	window.Show()

	return result
}
