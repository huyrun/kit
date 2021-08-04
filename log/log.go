package log

import "sync"

var doOnce sync.Once
var logging *log

func init() {
	doOnce.Do(func() {
		logging = new(log)
		logging.mode = false
	})
}

func LoggingMode(mode bool) {
	logging.mode = mode
}

type log struct {
	mode bool
	err  error
	info string
}

func Error(err error) LogBuilder {
	return logging.Error(err)
}

func Info(info string) LogBuilder {
	return logging.Info(info)
}

func Infof(format string, data ...interface{}) LogBuilder {
	return logging.Infof(format, data)
}
