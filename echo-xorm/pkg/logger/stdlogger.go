package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// StdLogger logs to Stdout
type StdLogger struct {
	logger *log.Logger
	id     string
	tag    string
}

// NewStdLogger is StdLogger constructor
func NewStdLogger(id, tag string) *StdLogger {
	l := StdLogger{
		logger: log.New(os.Stdout, "", 0),
		id:     id,
		tag:    tag,
	}

	return &l
}

// Info logs "info" messages. First value should be Event string
func (l *StdLogger) Info(values ...interface{}) {
	l.logTagged("info", values...)
}

// Error logs "error" messages. First value should be Event string
func (l *StdLogger) Error(values ...interface{}) {
	l.logTagged("error", values...)
}

// Warn logs "warning" messages. First value should be Event string
func (l *StdLogger) Warn(values ...interface{}) {
	l.logTagged("warning", values...)
}

// Close for Stdlogger does nothing
func (l *StdLogger) Close() {}

// main log implementation
func (l *StdLogger) logTagged(category string, values ...interface{}) {
	if len(values) <= 1 {
		l.logUntagged(values...)
		return
	}
	logData := map[string]interface{}{
		"App":   l.tag,
		"ID":    l.id,
		"Cat":   category,
		"Event": fmt.Sprint(values[0]),
		"Msg":   fmt.Sprintf("%v", values[1:]),
	}

	logString, err := json.Marshal(logData)
	if err != nil { // on json-error log it as string
		log.Printf(`{"App":%s, "ID":%s, "Cat":"error", "Event":"logger", "Msg":"json fail on %s"}`, l.tag, l.id, fmt.Sprintf("%+v", logData))
	} else {
		log.Printf("%s", string(logString))
	}
}

// untagged log implementation
func (l *StdLogger) logUntagged(values ...interface{}) {
	logData := map[string]interface{}{
		"App":   l.tag,
		"ID":    l.id,
		"Cat":   "unknown",
		"Event": "unknown",
		"Msg":   fmt.Sprintf("%v", values),
	}
	logString, err := json.Marshal(logData)
	if err != nil { // on json-error log it as string
		log.Printf(`{"App":%s, "ID":%s, "Cat":"error", "Event":"logger", "Msg":"json fail on %s"}`, l.tag, l.id, fmt.Sprintf("%+v", logData))
	} else {
		log.Printf("%s", string(logString))
	}
}
