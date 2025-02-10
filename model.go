package nabu

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type LogOutput int

const (
	OutputStderr LogOutput = iota
	OutputStdout
	OutputInternal
)

const (
	originError = iota
	originMessage
)

type Output struct {
	UUID     string   `json:",omitempty"`
	Date     string   `json:",omitempty"`
	Error    string   `json:",omitempty"`
	Args     any      `json:",omitempty"`
	Msg      string   `json:",omitempty"`
	Function string   `json:",omitempty"`
	Line     int      `json:",omitempty"`
	Level    LogLevel `json:",omitempty"`
}

type Logger struct {
	CausedBy error
	UUID     string
	Msg      string
	Args     any
	Level    LogLevel

	origin           int
	enableStackTrace bool
}
