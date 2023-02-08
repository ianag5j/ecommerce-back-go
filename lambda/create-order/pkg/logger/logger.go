package logger

import (
	"os"

	"github.com/rollbar/rollbar-go"
)

type (
	Logger interface {
		LogError(err error)
	}

	logger struct{}
)

func New() Logger {
	rollbar.SetToken(os.Getenv("ROLLBAR_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ENVIROMENT"))

	return logger{}
}

func (l logger) LogError(err error) {
	rollbar.Error(err)
	rollbar.Close()
}
