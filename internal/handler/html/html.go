package html

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

// Handle a Web request.
func Handle(params operations.WebParams) middleware.Responder {
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

// DefaultRedirectHandler redirects request to "/" URL to "web/index.html"
func DefaultRedirectHandler(_ operations.DefaultPageParams) middleware.Responder {
	return helper.NewCustomResponder(nil, map[string]string{"Location": "web/index.html"}, http.StatusPermanentRedirect)
}

// AppendEndpoints appends web endpoints to the API.
func AppendEndpoints(api *operations.MassaWalletAPI) {
	api.DefaultPageHandler = operations.DefaultPageHandlerFunc(DefaultRedirectHandler)
	api.WebHandler = operations.WebHandlerFunc(Handle)
}
