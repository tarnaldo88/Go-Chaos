package observability

import "log"

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{Logger: log.Default()}
}
