package server

import (
	"github.com/jakule/codersranktask/internal/storage"
	"go.uber.org/zap"
)

type CallParams struct {
	slog    *zap.SugaredLogger
	storage storage.Storage
}

func NewCallParams(logger *zap.Logger,
	storage storage.Storage) *CallParams {

	slog := logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	return &CallParams{slog: slog, storage: storage}
}

func (c *CallParams) Infof(template string, args ...interface{}) {
	c.slog.Infof(template, args...)
}

func (c *CallParams) Errorf(template string, args ...interface{}) {
	c.slog.Errorf(template, args...)
}

func (c *CallParams) Storage() storage.Storage {
	return c.storage
}
