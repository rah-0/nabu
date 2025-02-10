package nabu

var (
	logLevel       = LevelDebug
	logOutput      = OutputStderr
	internalOutput string
)

func SetLogLevel(l LogLevel) {
	logLevel = l
}

func SetLogOutput(o LogOutput) {
	logOutput = o
}

func shouldLog(l LogLevel) bool {
	return l >= logLevel
}
