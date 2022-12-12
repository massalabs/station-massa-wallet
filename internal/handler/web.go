package handler

import (
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	helper "github.com/massalabs/thyra-plugin-massa-wallet/pkg/openapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/web"
)

const indexHTML = "index.html"

//nolint:nolintlint,ireturn
func WebWalletHandler(params operations.WebParams) middleware.Responder {
	resourceName := params.Resource
	if params.Resource == indexHTML {
		resourceName = "wallet.html"
	}

	resourceContent, err := web.Content(resourceName)
	if err != nil {
		return operations.NewWebNotFound()
	}

	fileExtension := filepath.Ext(resourceName)

	mimeType := mime.TypeByExtension(fileExtension)

	header := map[string]string{"Content-Type": mimeType}

	return helper.NewCustomResponder(resourceContent, header, http.StatusOK)
}
