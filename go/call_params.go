package swagger

import (
	"context"

	"go.uber.org/zap"
)

type CallParams struct {
	ctx  context.Context
	slog *zap.SugaredLogger
}

func (c *CallParams) Infof(template string, args ...interface{}) {
	c.slog.Infof(template, args)
}

func (c *CallParams) Errorf(template string, args ...interface{}) {
	c.slog.Errorf(template, args)
}
