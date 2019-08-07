package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jakule/codersranktask/internal/server"
	"github.com/jakule/codersranktask/internal/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	router *mux.Router
	db     storage.Storage
	logger *zap.Logger
	slog   *zap.SugaredLogger
}

func (a *App) Init(dbConnStr string) {
	var err error
	a.logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	a.slog = a.logger.Sugar()

	a.db, err = storage.NewPgStorage(dbConnStr)
	if err != nil {
		a.slog.Fatal(err)
	}

	if err := MigrateDB(dbConnStr); err != nil {
		a.slog.Fatal(err)
	}

	c := server.NewCallParams(a.logger, a.db)
	a.router = server.NewRouter(c)
}

func (a *App) Run(port string) {
	addr, err := getAddr(port)
	if err != nil {
		a.slog.Fatal(err)
	}
	a.slog.Info("starting server")

	srv := &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		a.slog.Fatal(err)
	}
}

func getAddr(port string) (string, error) {
	if port == "" {
		return "", errors.New("PORT is not set")
	}
	return fmt.Sprintf(":%s", port), nil
}
