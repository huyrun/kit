package log

import (
	"fmt"
	"strings"
)

type LogBuilder interface {
	Error(err error) LogBuilder
	Info(info string) LogBuilder
	Infof(format string, data ...interface{}) LogBuilder
	build()
	reset()
}

func (l *log) Error(err error) LogBuilder {
	defer l.build()
	l.err = err
	return l
}

func (l *log) Info(info string) LogBuilder {
	defer l.build()
	l.info = info
	return l
}

func (l *log) Infof(format string, data ...interface{}) LogBuilder {
	defer l.build()
	l.info = fmt.Sprintf(format, data...)
	return l
}

func (l *log) build() {
	log := strings.Builder{}
	if l.err != nil {
		log.WriteString(fmt.Sprintf("err = %s/n", l.err.Error()))
	}
	if len(l.info) > 0 {
		log.WriteString(fmt.Sprintf("info = %s/n", l.info))
	}
	fmt.Println(log.String())
	l.reset()
}

func (l *log) reset() {
	l = new(log)
}
