package logger

// Logger is an interface for debug and logging callbacks.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

// Silent is an implementation of Logger that ignores every message it gets.
// It's used as the default set of logging callbacks if none are set.
type Silent struct{}

// Debug sends debug messages
func (s Silent) Debug(args ...interface{}) {}

// Debugf sends debug messages
func (s Silent) Debugf(format string, args ...interface{}) {}

// Error sends error messages
func (s Silent) Error(args ...interface{}) {}

// Errorf sends error messages
func (s Silent) Errorf(format string, args ...interface{}) {}

// Info sends info messages
func (s Silent) Info(args ...interface{}) {}

// Infof sends info messages
func (s Silent) Infof(format string, args ...interface{}) {}
