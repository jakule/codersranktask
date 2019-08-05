package internal

import (
	"context"

	"github.com/jakule/codersranktask/internal/storage"
	"go.uber.org/zap"
)

type CallParams struct {
	ctx     context.Context
	slog    *zap.SugaredLogger
	storage storage.Storage
}

func NewCallParams(ctx context.Context, slog *zap.SugaredLogger,
	storage storage.Storage) *CallParams {

	return &CallParams{ctx: ctx, slog: slog, storage: storage}
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
