package logger

import (
	"context"
	"os"

	"github.com/rollbar/rollbar-go"
)

type (
	Logger interface {
		LogError(err error)
		LogInfo(info ...any)
	}

	logger struct{}
)

func New(ctx context.Context) Logger {
	rollbar.SetToken(os.Getenv("ROLLBAR_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ENVIROMENT"))
	rollbar.SetContext(ctx)

	return logger{}
}

func (l logger) LogError(err error) {
	rollbar.Error(err)
	rollbar.Close()
}

func (l logger) LogInfo(info ...any) {
	rollbar.Info(info)
	rollbar.Close()
}
