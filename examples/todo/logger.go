package todo

import "log"

type (
	Logger interface {
		Log(...any)
	}

	logger struct{}
)

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) Log(args ...any) {
	log.Println(args...)
}
