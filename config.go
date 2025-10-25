package nabu

import "sync"

var (
	// configMutex protects access to global configuration variables
	configMutex sync.RWMutex

	// logLevel is the minimum severity level that will be logged
	// Default is LevelDebug (all logs will be displayed)
	logLevel = LevelDebug

	// logOutput determines where logs will be written
	// Default is OutputStderr (standard error)
	logOutput = OutputStderr

	// internalOutput is a buffer used to capture logs for testing
	// when OutputInternal is selected
	internalOutput string
)

// SetLogLevel configures the minimum log level that will be processed.
// Logs with a level lower than this will be ignored.
// Default is LevelDebug (all logs will be displayed).
func SetLogLevel(l LogLevel) {
	configMutex.Lock()
	defer configMutex.Unlock()
	logLevel = l
}

// SetLogOutput configures where logs will be written.
// Options are OutputStderr (default), OutputStdout, or OutputInternal.
func SetLogOutput(o LogOutput) {
	configMutex.Lock()
	defer configMutex.Unlock()
	logOutput = o
}

// shouldLog determines if a log with the given level should be processed
// based on the current configured log level.
func shouldLog(l LogLevel) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return l >= logLevel
}
