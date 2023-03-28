package password

import (
	"github.com/ncruces/zenity"
)

type PasswordDialog struct {
	title       string
	okLabel     string
	cancelLabel string
}

func Create() (string, error) {
	onCreate := PasswordDialog{
		title:       "Create a new passsord",
		okLabel:     "Create",
		cancelLabel: "Cancel",
	}
	return onCreate.dialog()
}

func Read() (string, error) {
	onUpdate := PasswordDialog{
		title:       "Please, type your passsord",
		okLabel:     "Ok",
		cancelLabel: "Cancel",
	}
	return onUpdate.dialog()
}

func (p *PasswordDialog) dialog() (string, error) {
	_, rawPass, err := zenity.Password(
		zenity.Title(p.title),
		zenity.OKLabel(p.okLabel),
		zenity.CancelLabel(p.cancelLabel))

	if err != nil {
		return rawPass, err
	}

	return rawPass, nil
}
