/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	swagger "github.com/jakule/codersranktask/internal"
	"github.com/jakule/codersranktask/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc paramHandler
}

type paramHandler func(c *swagger.CallParams, w http.ResponseWriter, r *http.Request)

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

	Route{
		"AddSecret",
		strings.ToUpper("Post"),
		"/v1/secret",
		AddSecret,
	},

	Route{
		"GetSecretByHash",
		strings.ToUpper("Get"),
		"/v1/secret/{hash}",
		GetSecretByHash,
	},
}

func handlerWrapperLogger(params *swagger.CallParams, inner paramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner(params, w, r)
	})
}

func newCallParams(dbConnStr string) *swagger.CallParams {
	s, err := storage.NewPgStorage(dbConnStr)
	if err != nil {
		panic(err)
	}
	return swagger.NewCallParams(context.Background(),
		mustLogger(newProdLogger()).Sugar(), s)
}

func newProdLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func mustLogger(logger *zap.Logger, err error) *zap.Logger {
	if err != nil {
		panic(err)
	}
	return logger
}

func NewRouter(dbConnStr string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		callParams := newCallParams(dbConnStr)
		handler := handlerWrapperLogger(callParams, route.HandlerFunc)
		handler = swagger.Logger(callParams, handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.Name("metrics").Handler(promhttp.Handler())

	return router
}