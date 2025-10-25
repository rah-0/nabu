package nabu

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// New creates a new empty Logger instance.
func New() *Logger {
	return &Logger{}
}

// FromError creates a Logger instance from an error.
// If the error is nil, an empty Logger is returned.
// If the error is a *Logger, its UUID is preserved to maintain the error chain.
// Otherwise, a new UUID is generated for tracking related logs.
func FromError(e error) *Logger {
	x := New()
	x.origin = originError
	x.Level = LevelError
	if e == nil {
		return x
	}

	x.CausedBy = e
	x.enableStackTrace = true

	var ex *Logger
	if errors.As(e, &ex) {
		x.UUID = ex.UUID
	} else {
		x.UUID = uuid.NewString()
	}

	return x
}

// FromMessage creates a Logger instance from a message string.
// The default log level is set to LevelInfo.
func FromMessage(msg string) *Logger {
	x := New()
	x.origin = originMessage
	x.Level = LevelInfo
	x.Msg = msg
	return x
}

// WithArgs attaches structured data to the log entry.
// This data will be serialized as part of the JSON output.
func (x *Logger) WithArgs(args ...any) *Logger {
	x.Args = args
	return x
}

// WithMessage sets or updates the log message.
func (x *Logger) WithMessage(msg string) *Logger {
	x.Msg = msg
	return x
}

// WithUuid sets a custom UUID for the log entry.
// This is useful for correlating logs across different services or from external sources.
func (x *Logger) WithUuid(uuid string) *Logger {
	x.UUID = uuid
	return x
}

// WithLevelDebug sets the log level to Debug.
func (x *Logger) WithLevelDebug() *Logger {
	x.Level = LevelDebug
	return x
}

// WithLevelInfo sets the log level to Info.
func (x *Logger) WithLevelInfo() *Logger {
	x.Level = LevelInfo
	return x
}

// WithLevelWarn sets the log level to Warning.
func (x *Logger) WithLevelWarn() *Logger {
	x.Level = LevelWarn
	return x
}

// WithLevelError sets the log level to Error.
func (x *Logger) WithLevelError() *Logger {
	x.Level = LevelError
	return x
}

// WithLevelFatal sets the log level to Fatal.
func (x *Logger) WithLevelFatal() *Logger {
	x.Level = LevelFatal
	return x
}

// EnableStackTrace forces inclusion of stack trace information (function name and line number).
// By default, stack traces are only enabled for error logs.
func (x *Logger) EnableStackTrace() *Logger {
	x.enableStackTrace = true
	return x
}

// Error implements the error interface to allow using Logger as an error.
// Returns the underlying error message or empty string if no error is present.
func (x *Logger) Error() string {
	if x.CausedBy != nil {
		return x.CausedBy.Error()
	}
	return ""
}

// Unwrap implements the unwrappable interface to allow error chain inspection.
// Returns the underlying error that was logged.
func (x *Logger) Unwrap() error {
	return x.CausedBy
}

// Log outputs the log entry and returns itself as an error.
// This method checks if the log level is enabled before writing the log.
// If the log originates from an error but no error is set, nothing is logged.
// The log entry includes timestamp, UUID, message/error, arguments and stack trace if enabled.
// The output is directed to the configured output destination (stderr by default).
func (x *Logger) Log() error {
	if !shouldLog(x.Level) {
		return x
	}

	if x.origin == originError && x.CausedBy == nil {
		return x
	}

	o := Output{
		UUID:  x.UUID,
		Date:  getDate(),
		Args:  x.Args,
		Msg:   x.Msg,
		Level: x.Level,
	}
	if x.CausedBy != nil {
		// Only show the immediate error, not the full chain
		// This prevents error duplication across the stack
		var loggerErr *Logger
		if errors.As(x.CausedBy, &loggerErr) {
			// CausedBy is a Logger - don't show its error in this log
			// The error was already logged when that Logger was created
		} else {
			// CausedBy is a regular error - show it
			o.Error = x.CausedBy.Error()
		}
	}
	if x.enableStackTrace {
		o.Function, o.Line = x.getFirstTrace()
	}

	log := toJson(o)

	switch logOutput {
	case OutputInternal:
		internalOutput += strings.TrimSpace(log) + "\n"
	case OutputStdout:
		fmt.Fprintln(os.Stdout, log)
	case OutputStderr:
		fmt.Fprintln(os.Stderr, log)
	}

	return x
}
