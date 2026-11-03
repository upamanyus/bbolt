package bbolt

// See https://github.com/etcd-io/raft/blob/main/logger.go
import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

func getDiscardLogger() Logger {
	return discardLogger
}

var (
	discardLogger Logger = &nopLogger{}
)

// nopLogger is bbolt's logger when Options.Logger is nil. It discards every
// message without formatting it, as a DefaultLogger writing to io.Discard
// does after formatting; Fatal and Panic still exit and panic.
type nopLogger struct{}

func (*nopLogger) Debug(v ...interface{})                   {}
func (*nopLogger) Debugf(format string, v ...interface{})   {}
func (*nopLogger) Error(v ...interface{})                   {}
func (*nopLogger) Errorf(format string, v ...interface{})   {}
func (*nopLogger) Info(v ...interface{})                    {}
func (*nopLogger) Infof(format string, v ...interface{})    {}
func (*nopLogger) Warning(v ...interface{})                 {}
func (*nopLogger) Warningf(format string, v ...interface{}) {}
func (*nopLogger) Fatal(v ...interface{})                   { os.Exit(1) }
func (*nopLogger) Fatalf(format string, v ...interface{})   { os.Exit(1) }
func (*nopLogger) Panic(v ...interface{})                   { panic(fmt.Sprint(v...)) }
func (*nopLogger) Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}

const (
	calldepth = 2
)

// DefaultLogger is a default implementation of the Logger interface.
type DefaultLogger struct {
	*log.Logger
	debug bool
}

func (l *DefaultLogger) EnableTimestamps() {
	l.SetFlags(l.Flags() | log.Ldate | log.Ltime)
}

func (l *DefaultLogger) EnableDebug() {
	l.debug = true
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	if l.debug {
		_ = l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
	}
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	if l.debug {
		_ = l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
	}
}

func (l *DefaultLogger) Info(v ...interface{}) {
	_ = l.Output(calldepth, header("INFO", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("INFO", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Error(v ...interface{}) {
	_ = l.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Warning(v ...interface{}) {
	_ = l.Output(calldepth, header("WARN", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Warningf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("WARN", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Fatal(v ...interface{}) {
	_ = l.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *DefaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}
