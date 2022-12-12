package web

import (
	"embed"
)

const basePath = "html/"

//go:embed html
var content embed.FS

func Content(resource string) ([]byte, error) {
	return content.ReadFile(basePath + resource)
}
