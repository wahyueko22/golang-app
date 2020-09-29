package logger

// Logger is an interface for logging
type Logger interface {
	Info(values ...interface{})  // used to log "info" messages
	Error(values ...interface{}) // used to log "error" messages
	Warn(values ...interface{})  // used to log "warning" messages
	Close()
}
