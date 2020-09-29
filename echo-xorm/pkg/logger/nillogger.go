package logger

// NilLogger logs nothing
type NilLogger struct {
}

// NewNilLogger is a constructor
func NewNilLogger() *NilLogger {
	return new(NilLogger)
}

// Info do nothing, just match the interface
func (l *NilLogger) Info(values ...interface{}) {}

// Error do nothing, just match the interface
func (l *NilLogger) Error(values ...interface{}) {}

// Warn do nothing, just match the interface
func (l *NilLogger) Warn(values ...interface{}) {}

// Close for NilLogger do nothing
func (l *NilLogger) Close() {}
