package observer

import (
	"context"
	"fmt"
)

func New(ctx context.Context, conf Configuration) *Observer {
	return &Observer{
		logger: conf.Logger,
	}
}

func (o *Observer) Error(ctx context.Context, err error, logs ...KV) {
	o.logger.Error(ctx, err, logs...)
}

func (o *Observer) Warn(ctx context.Context, msg string, logs ...KV) {
	o.logger.Warn(ctx, msg, logs...)
}

func (o *Observer) Info(ctx context.Context, msg string, logs ...KV) {
	o.logger.Info(ctx, msg, logs...)
}

func (o *Observer) Debug(ctx context.Context, msg string, logs ...KV) {
	o.logger.Debug(ctx, msg, logs...)
}

func (o *Observer) Close(ctx context.Context) error {
	if err := o.logger.Close(ctx); err != nil {
		return fmt.Errorf("logger closing: %w", o.logger.Close(ctx))
	}

	return nil
}
