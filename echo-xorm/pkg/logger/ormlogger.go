package logger

import (
	"fmt"

	"github.com/go-xorm/core"
)

// OrmLogger is implementation of xorm core.ILogger with proxy logger
type OrmLogger struct {
	proxy Logger
}

// NewOrmLogger is logger constructor
func NewOrmLogger(proxyLogger Logger) *OrmLogger {
	l := new(OrmLogger)
	l.proxy = proxyLogger
	return l
}

// Error implement core.ILogger
func (l *OrmLogger) Error(v ...interface{}) {
	l.proxy.Error("SQL", fmt.Sprint(v...))
}

// Errorf implement core.ILogger
func (l *OrmLogger) Errorf(format string, v ...interface{}) {
	l.proxy.Error("SQL", fmt.Sprintf(format, v...))
}

// Debug implement core.ILogger
func (l *OrmLogger) Debug(v ...interface{}) {
	l.Error(v...)
}

// Debugf implement core.ILogger
func (l *OrmLogger) Debugf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

// Info implement core.ILogger
func (l *OrmLogger) Info(v ...interface{}) {
	l.proxy.Info("SQL", fmt.Sprint(v...))
}

// Infof implement core.ILogger
func (l *OrmLogger) Infof(format string, v ...interface{}) {
	l.proxy.Info("SQL", fmt.Sprintf(format, v...))
}

// Warn implement core.ILogger
func (l *OrmLogger) Warn(v ...interface{}) {
	l.proxy.Warn("SQL", fmt.Sprint(v...))
}

// Warnf implement core.ILogger
func (l *OrmLogger) Warnf(format string, v ...interface{}) {
	l.proxy.Warn("SQL", fmt.Sprintf(format, v...))
}

// Level implement core.ILogger
func (l *OrmLogger) Level() core.LogLevel {
	return core.LOG_UNKNOWN
}

// SetLevel implement core.ILogger
func (l *OrmLogger) SetLevel(core.LogLevel) {}

// ShowSQL implement core.ILogger
func (l *OrmLogger) ShowSQL(...bool) {}

// IsShowSQL implement core.ILogger
func (l *OrmLogger) IsShowSQL() bool {
	return true
}
