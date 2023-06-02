package html

import (
	"embed"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/openapi"
	"github.com/massalabs/thyra-plugin-wallet/web"
)

const (
	indexHTML      = "index.html"
	basePathWebApp = "dist/"
)

//nolint:typecheck,nolintlint
//go:embed dist
var contentWebApp embed.FS

// Handle a Web request.
// For legacy frontend
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

	return openapi.NewCustomResponder(resourceContent, header, http.StatusOK)
}

// Handle a Web request.
func HandleWebApp(params operations.WebAppParams) middleware.Responder {
	resourceName := params.Resource

	resourceContent, err := contentWebApp.ReadFile(basePathWebApp + resourceName)
	if err != nil {
		resourceName = "index.html"
		resourceContent, err = contentWebApp.ReadFile(basePathWebApp + resourceName)
		if err != nil {
			return operations.NewWebAppNotFound()
		}
	}

	fileExtension := filepath.Ext(resourceName)

	mimeType := mime.TypeByExtension(fileExtension)

	header := map[string]string{"Content-Type": mimeType}

	return openapi.NewCustomResponder(resourceContent, header, http.StatusOK)
}

// DefaultRedirectHandler redirects request to "/" URL to "web-app/index"
func DefaultRedirectHandler(_ operations.DefaultPageParams) middleware.Responder {
	return openapi.NewCustomResponder(nil, map[string]string{"Location": "web-app/index"}, http.StatusPermanentRedirect)
}

// AppendEndpoints appends web endpoints to the API.
func AppendEndpoints(api *operations.MassaWalletAPI) {
	api.DefaultPageHandler = operations.DefaultPageHandlerFunc(DefaultRedirectHandler)
	api.WebHandler = operations.WebHandlerFunc(Handle)
}
