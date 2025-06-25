package nabu

import (
	"time"
)

// LogLevel defines the severity of a log entry.
type LogLevel int

const (
	// LevelDebug is used for verbose debugging information
	LevelDebug LogLevel = iota
	// LevelInfo is used for general information about application progress
	LevelInfo
	// LevelWarn is used for non-critical issues that might require attention
	LevelWarn
	// LevelError is used for errors that affect normal operation but don't cause termination
	LevelError
	// LevelFatal is used for critical errors that may cause application termination
	LevelFatal
)

// LogOutput defines where the log entries will be written.
type LogOutput int

const (
	// OutputStderr sends logs to standard error stream
	OutputStderr LogOutput = iota
	// OutputStdout sends logs to standard output stream
	OutputStdout
	// OutputInternal stores logs in an internal buffer for testing
	OutputInternal
)

const (
	// originError indicates the log entry originated from an error
	originError = iota
	// originMessage indicates the log entry originated from a message
	originMessage
)

const TimeLayout = "2006-01-02 15:04:05.000000"

// Output represents the JSON structure of a log entry.
type Output struct {
	UUID     string   `json:",omitempty"` // Unique identifier for tracking related log entries
	Date     string   `json:",omitempty"` // Timestamp when log was created
	Error    string   `json:",omitempty"` // Error message if this is an error log
	Args     any      `json:",omitempty"` // Additional structured data for the log entry
	Msg      string   `json:",omitempty"` // Main log message
	Function string   `json:",omitempty"` // Function where the log was generated
	Line     int      `json:",omitempty"` // Line number where the log was generated
	Level    LogLevel `json:",omitempty"` // Severity level of the log
}

// Logger is the main logging object that holds log details before they're written.
type Logger struct {
	CausedBy error    // Original error that caused this log entry
	UUID     string   // Unique identifier for related log entries
	Msg      string   // Log message
	Args     any      // Additional structured data
	Level    LogLevel // Severity level

	origin           int  // Whether the log originated from an error or message
	enableStackTrace bool // Whether to include stack trace information
}

type ParsedErrorTrace struct {
	UUID   string
	Error  string
	Frames []Output // Ordered oldest to newest
}

type ParsedLogs struct {
	Entries []Output
	Traces  []ParsedErrorTrace
}

type Parser struct {
	lines     []string
	afterDate *time.Time
}
