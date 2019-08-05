package swagger

import (
	"context"

	"github.com/jakule/codersranktask/go/storage"
	"go.uber.org/zap"
)

type CallParams struct {
	ctx     context.Context
	slog    *zap.SugaredLogger
	storage storage.Storage
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
