package nabu

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func New() *Logger {
	return &Logger{
		UUID: uuid.NewString(),
	}
}

func FromError(e error) *Logger {
	x := New()
	x.origin = originError
	x.Level = LevelError
	if e == nil {
		return x
	}

	x.CausedBy = e
	var ex *Logger
	if errors.As(e, &ex) {
		x.UUID = ex.UUID
	}

	return x
}

func FromMessage(msg string) *Logger {
	x := New()
	x.origin = originMessage
	x.Level = LevelInfo
	x.Msg = msg
	return x
}

func (x *Logger) WithArgs(args ...any) *Logger {
	x.Args = args
	return x
}

func (x *Logger) WithMessage(msg string) *Logger {
	x.Msg = msg
	return x
}

func (x *Logger) WithLevelDebug() *Logger {
	x.Level = LevelDebug
	return x
}

func (x *Logger) WithLevelInfo() *Logger {
	x.Level = LevelInfo
	return x
}

func (x *Logger) WithLevelWarn() *Logger {
	x.Level = LevelWarn
	return x
}

func (x *Logger) WithLevelError() *Logger {
	x.Level = LevelError
	return x
}

func (x *Logger) WithLevelFatal() *Logger {
	x.Level = LevelFatal
	return x
}

func (x *Logger) Error() string {
	if x.CausedBy != nil {
		return x.CausedBy.Error()
	}
	return ""
}

func (x *Logger) Log() error {
	if !shouldLog(x.Level) {
		return x
	}

	if x.origin == originError && x.CausedBy == nil {
		return x
	}

	fn, l := x.getFirstTrace()
	o := Output{
		UUID:     x.UUID,
		Date:     getDate(),
		Args:     x.Args,
		Msg:      x.Msg,
		Function: fn,
		Line:     l,
		Level:    x.Level,
	}
	if x.CausedBy != nil {
		o.Error = x.CausedBy.Error()
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
