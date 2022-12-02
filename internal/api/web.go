package api

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-core/api/server/restapi/operations"
)

const indexHTML = "index.html"

const basePath = "html/front/"

//go:embed html/front
var content embed.FS

//nolint:nolintlint,ireturn
func WebWalletHandler(params operations.WebWalletParams) middleware.Responder {
	file := params.Resource
	if params.Resource == indexHTML {
		file = "wallet.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewWebWalletNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}
