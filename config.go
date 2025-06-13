package nabu

var (
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
	logLevel = l
}

// SetLogOutput configures where logs will be written.
// Options are OutputStderr (default), OutputStdout, or OutputInternal.
func SetLogOutput(o LogOutput) {
	logOutput = o
}

// shouldLog determines if a log with the given level should be processed
// based on the current configured log level.
func shouldLog(l LogLevel) bool {
	return l >= logLevel
}
