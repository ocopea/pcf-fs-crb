// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"crb/restapi/impl"
	"crb/restapi/operations"
	"crb/restapi/operations/crb_web"
	"crb/utils"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name crb --spec ../swagger.yaml

func configureFlags(api *operations.CrbAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CrbAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.BinConsumer = runtime.ByteStreamConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.BinProducer = runtime.ByteStreamProducer()

	api.CrbWebCreateCopyHandler = crb_web.CreateCopyHandlerFunc(func(params crb_web.CreateCopyParams) middleware.Responder {
		return impl.CreateCopyHandler(params, mysqldb, sftpClient)
	})
	api.CrbWebDeleteCopyHandler = crb_web.DeleteCopyHandlerFunc(func(params crb_web.DeleteCopyParams) middleware.Responder {
		return impl.DeleteCopyHandler(params, mysqldb, sftpClient)
	})
	api.CrbWebGetCopyMetaDataHandler = crb_web.GetCopyMetaDataHandlerFunc(func(params crb_web.GetCopyMetaDataParams) middleware.Responder {
		return impl.GetCopyMetaData(params)
	})
	api.CrbWebGetCopyinstancesHandler = crb_web.GetCopyinstancesHandlerFunc(func(params crb_web.GetCopyinstancesParams) middleware.Responder {
		return impl.GetCopyInstances(params)
	})
	api.CrbWebGetInfoHandler = crb_web.GetInfoHandlerFunc(func(params crb_web.GetInfoParams) middleware.Responder {
		return crb_web.NewGetInfoOK().WithPayload(impl.CrbInfoResponse())
	})
	api.CrbWebGetRepositoryInfoHandler = crb_web.GetRepositoryInfoHandlerFunc(func(params crb_web.GetRepositoryInfoParams) middleware.Responder {
		return (impl.GetRepositoryInfoHandler(params))
	})
	api.CrbWebRetrieveCopyHandler = crb_web.RetrieveCopyHandlerFunc(func(params crb_web.RetrieveCopyParams) middleware.Responder {
		return impl.RetrieveCopyHandler(params, mysqldb, sftpClient)
	})
	api.CrbWebStoreRepositoryInfoHandler = crb_web.StoreRepositoryInfoHandlerFunc(func(params crb_web.StoreRepositoryInfoParams) middleware.Responder {
		return impl.StoreRepositoryInfoHandler(params)
	})

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
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

var mysqldb *utils.Database
var sftpClient *utils.SftpClient

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.EscapedPath(), "/data") ||
			strings.Contains(r.Method, "DELETE") {
			sftpClient = impl.GetSftpConnection(mysqldb, sftpClient)
			if sftpClient != nil {
				defer utils.CloseSftpConnection(sftpClient)
			}
		}
		handler.ServeHTTP(w, r)
	})
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
