package xlog

import (
	"github.com/op/go-logging"
	"os"
)

var simpleFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} %{level:.4s} %{message}",
)

var colorFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{level:.4s}%{color:reset} %{message}",
)

var f *os.File

func Open(fileName string) (*logging.Logger, error) {
	Close()

	var err error
	f, err = os.OpenFile(fileName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(f, "", 0)

	backend1Formatter := logging.NewBackendFormatter(backend1, colorFormat)
	backend2Formatter := logging.NewBackendFormatter(backend2, simpleFormat)

	logging.SetBackend(backend1Formatter, backend2Formatter)
	return logging.MustGetLogger(fileName), nil
}

func Close() {
	if f != nil {
		f.Close()
	}
}
