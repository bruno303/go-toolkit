package log

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAsyncDecorator_Info(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)

	logger.EXPECT().Info(gomock.Any(), "test message", gomock.Any()).Times(1)

	ad := NewAsyncDecorator(logger)

	ad.Info(context.Background(), "test message")
	ad.Shutdown(ctx)
}

func TestAsyncDecorator_Debug(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)

	logger.EXPECT().Debug(gomock.Any(), "test message", gomock.Any()).Times(1)
	ad := NewAsyncDecorator(logger)

	ad.Debug(context.Background(), "test message")
	ad.Shutdown(ctx)
}

func TestAsyncDecorator_Warn(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)

	logger.EXPECT().Warn(gomock.Any(), "test message", gomock.Any()).Times(1)
	ad := NewAsyncDecorator(logger)

	ad.Warn(context.Background(), "test message")
	ad.Shutdown(ctx)
}

func TestAsyncDecorator_Error(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)

	expectedErr := errors.New("test error")

	logger.EXPECT().Error(gomock.Any(), "test message", gomock.Eq(expectedErr)).Times(1)
	ad := NewAsyncDecorator(logger)

	ad.Error(context.Background(), "test message", expectedErr)
	ad.Shutdown(ctx)
}
