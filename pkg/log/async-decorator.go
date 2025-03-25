package log

import (
	"context"
	"sync"
)

type AsyncDecorator struct {
	logger Logger
	ch     chan func()
	wg     sync.WaitGroup
}

var _ Logger = (*AsyncDecorator)(nil)

func NewAsyncDecorator(logger Logger) *AsyncDecorator {
	return NewAsyncDecoratorWithBuffer(logger, 10)
}

func NewAsyncDecoratorWithBuffer(logger Logger, bufferSize int) *AsyncDecorator {
	ad := &AsyncDecorator{
		logger: logger,
		ch:     make(chan func(), bufferSize),
		wg:     sync.WaitGroup{},
	}
	ad.wg.Add(1)
	go ad.dispatchLogs()
	return ad
}

func (ad *AsyncDecorator) Info(ctx context.Context, msg string, args ...any) {
	ad.ch <- func() {
		ad.logger.Info(ctx, msg, args...)
	}
}

func (ad *AsyncDecorator) Debug(ctx context.Context, msg string, args ...any) {
	ad.ch <- func() {
		ad.logger.Debug(ctx, msg, args...)
	}
}

func (ad *AsyncDecorator) Warn(ctx context.Context, msg string, args ...any) {
	ad.ch <- func() {
		ad.logger.Warn(ctx, msg, args...)
	}
}

func (ad *AsyncDecorator) Error(ctx context.Context, msg string, err error) {
	ad.ch <- func() {
		ad.logger.Error(ctx, msg, err)
	}
}

func (ad *AsyncDecorator) SetLevel(l Level) error {
	return ad.logger.SetLevel(l)
}

func (ad *AsyncDecorator) Shutdown(context.Context) error {
	close(ad.ch)
	ad.wg.Wait()
	return nil
}

func (ad *AsyncDecorator) dispatchLogs() {
	defer ad.wg.Done()
	for fn := range ad.ch {
		fn()
	}
}
