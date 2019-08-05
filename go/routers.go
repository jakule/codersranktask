/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc paramHandler
}

type paramHandler func(c *CallParams, w http.ResponseWriter, r *http.Request)

type Routes []Route

func handlerWrapperLogger(params *CallParams, inner paramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner(params, w, r)
	})
}

func createCallParams(dbConnStr string) *CallParams {
	storage, err := NewPgStorage(dbConnStr)
	if err != nil {
		panic(err)
	}
	return &CallParams{
		ctx:     context.Background(),
		slog:    mustLogger(newProdLogger()).Sugar(),
		storage: storage,
	}
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
		callParams := createCallParams(dbConnStr)
		handler := handlerWrapperLogger(callParams, route.HandlerFunc)
		handler = Logger(callParams, handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(c *CallParams, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

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
