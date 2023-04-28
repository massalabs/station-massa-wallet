// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name MassaWallet --spec ../../walletApi-V0.yml --principal interface{} --exclude-main

func configureFlags(api *operations.MassaWalletAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MassaWalletAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinProducer = runtime.ByteStreamProducer()
	api.CSSProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("css producer has not yet been implemented")
	})
	api.HTMLProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("html producer has not yet been implemented")
	})
	api.JsProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("js producer has not yet been implemented")
	})
	api.JSONProducer = runtime.JSONProducer()
	api.TextWebpProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("textWebp producer has not yet been implemented")
	})

	if api.DefaultPageHandler == nil {
		api.DefaultPageHandler = operations.DefaultPageHandlerFunc(func(params operations.DefaultPageParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.DefaultPage has not yet been implemented")
		})
	}
	if api.RestAccountDeleteHandler == nil {
		api.RestAccountDeleteHandler = operations.RestAccountDeleteHandlerFunc(func(params operations.RestAccountDeleteParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RestAccountDelete has not yet been implemented")
		})
	}
	if api.RestAccountGetHandler == nil {
		api.RestAccountGetHandler = operations.RestAccountGetHandlerFunc(func(params operations.RestAccountGetParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RestAccountGet has not yet been implemented")
		})
	}
	if api.RestAccountImportHandler == nil {
		api.RestAccountImportHandler = operations.RestAccountImportHandlerFunc(func(params operations.RestAccountImportParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RestAccountImport has not yet been implemented")
		})
	}
	if api.RestAccountListHandler == nil {
		api.RestAccountListHandler = operations.RestAccountListHandlerFunc(func(params operations.RestAccountListParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RestAccountList has not yet been implemented")
		})
	}
	if api.RestAccountSignOperationHandler == nil {
		api.RestAccountSignOperationHandler = operations.RestAccountSignOperationHandlerFunc(func(params operations.RestAccountSignOperationParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RestAccountSignOperation has not yet been implemented")
		})
	}
	if api.WebHandler == nil {
		api.WebHandler = operations.WebHandlerFunc(func(params operations.WebParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Web has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodDelete},
	}).Handler

	return handleCORS(handler)
}
