package file

import (
	"github.com/ncruces/zenity"
)

type FileDialog struct {
	title string
}

func Open() (string, error) {
	onOpen := FileDialog{
		title: "please, select the file",
	}
	return onOpen.dialog()
}

func (p *FileDialog) dialog() (string, error) {
	const defaultPath = ``

	filePath, err := zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilters{zenity.FileFilter{
			Name:     "dat files",
			Patterns: []string{"*.dat"},
			CaseFold: true}})
	if err != nil {
		return filePath, err
	}

	return filePath, nil
}
