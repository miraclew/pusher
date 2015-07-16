package log

import (
	"io"
	logp "log"
)

type Logger struct {
	*logp.Logger
	Debug bool
}

func New(out io.Writer) *Logger {
	return NewNamed(out, "")
}

func NewNamed(out io.Writer, prefix string) *Logger {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = "[" + prefix + "] "
	}

	logger := logp.New(out, prefixStr, logp.LstdFlags)

	return &Logger{logger, false}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		l.Logger.Printf("DEBUG: "+format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.Debug {
		l.Logger.Printf("ERROR: "+format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Debug {
		l.Logger.Printf("INFO: "+format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.Debug {
		l.Logger.Printf("WARN: "+format, args...)
	}
}
