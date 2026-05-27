// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mkheyfets/ispro-app/internal/metrics"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/entries"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/links"
)

//go:generate swagger generate server --target ../../ispro-app-server --name IsproApp --spec ../../../api/openapi.yaml --principal string

func configureFlags(api *operations.IsproAppAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.IsproAppAPI) http.Handler {
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	if api.EntriesCreateEntryHandler == nil {
		api.EntriesCreateEntryHandler = entries.CreateEntryHandlerFunc(func(params entries.CreateEntryParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation entries.CreateEntry has not yet been implemented")
		})
	}
	if api.LinksCreateLinkHandler == nil {
		api.LinksCreateLinkHandler = links.CreateLinkHandlerFunc(func(params links.CreateLinkParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation links.CreateLink has not yet been implemented")
		})
	}
	if api.EntriesDeleteEntryHandler == nil {
		api.EntriesDeleteEntryHandler = entries.DeleteEntryHandlerFunc(func(params entries.DeleteEntryParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation entries.DeleteEntry has not yet been implemented")
		})
	}
	if api.LinksDeleteLinkHandler == nil {
		api.LinksDeleteLinkHandler = links.DeleteLinkHandlerFunc(func(params links.DeleteLinkParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation links.DeleteLink has not yet been implemented")
		})
	}
	if api.EntriesGetEntryHandler == nil {
		api.EntriesGetEntryHandler = entries.GetEntryHandlerFunc(func(params entries.GetEntryParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation entries.GetEntry has not yet been implemented")
		})
	}
	if api.LinksGetLinkHandler == nil {
		api.LinksGetLinkHandler = links.GetLinkHandlerFunc(func(params links.GetLinkParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation links.GetLink has not yet been implemented")
		})
	}
	if api.EntriesListEntriesHandler == nil {
		api.EntriesListEntriesHandler = entries.ListEntriesHandlerFunc(func(params entries.ListEntriesParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation entries.ListEntries has not yet been implemented")
		})
	}
	if api.LinksListLinksHandler == nil {
		api.LinksListLinksHandler = links.ListLinksHandlerFunc(func(params links.ListLinksParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation links.ListLinks has not yet been implemented")
		})
	}
	if api.EntriesUpdateEntryHandler == nil {
		api.EntriesUpdateEntryHandler = entries.UpdateEntryHandlerFunc(func(params entries.UpdateEntryParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation entries.UpdateEntry has not yet been implemented")
		})
	}
	if api.LinksUpdateLinkHandler == nil {
		api.LinksUpdateLinkHandler = links.UpdateLinkHandlerFunc(func(params links.UpdateLinkParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation links.UpdateLink has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	_ = tlsConfig
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(server *http.Server, scheme, addr string) {
	mux := http.NewServeMux()
	wrapped := metrics.NewMiddleware()(server.Handler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", wrapped)
	server.Handler = mux
	_ = scheme
	_ = addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return metrics.NewMiddleware()(handler)
}
