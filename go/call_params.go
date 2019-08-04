package swagger

import (
	"context"

	"go.uber.org/zap"
)

type CallParams struct {
	ctx     context.Context
	slog    *zap.SugaredLogger
	storage Storage
}

func (c *CallParams) Infof(template string, args ...interface{}) {
	c.slog.Infof(template, args...)
}

func (c *CallParams) Errorf(template string, args ...interface{}) {
	c.slog.Errorf(template, args...)
}

func (c *CallParams) Storage() Storage {
	return c.storage
}
