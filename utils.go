package nabu

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

// getDate returns the current UTC time formatted as a string.
// Format: YYYY-MM-DD HH:MM:SS.microseconds
func getDate() string {
	ct := time.Now().UTC()
	return ct.Format(TimeLayout)
}

// fromJson parses a JSON string into an Output struct.
// Returns nil if the JSON is invalid or cannot be parsed.
func fromJson(jsonString string) *Output {
	jsonString = strings.TrimSpace(jsonString)

	var o Output
	err := json.Unmarshal([]byte(jsonString), &o)
	if err != nil {
		return nil
	}
	return &o
}

// toJson converts any object to a JSON string.
// If marshaling fails, it returns a fallback JSON error message with a new UUID.
func toJson(targetObject any) (output string) {
	b, e := json.Marshal(targetObject)
	if e != nil {
		// Fallback to a manually constructed error JSON if marshaling fails
		output = `{"UUID":"` + uuid.NewString() + `","Date":"` + getDate() + `","Error":"` + e.Error() + `","Level":4}`
		return
	}
	output = string(b)
	return
}

// getFirstTrace captures the first relevant stack frame information.
// Returns the function name and line number of the caller that generated the log.
// Skips internal frames related to the logging machinery itself.
func (x *Logger) getFirstTrace() (string, int) {
	f := make([]uintptr, 1)
	callersCount := runtime.Callers(3, f) // Ignore: runtime.Callers, getFirstTrace, Log

	if callersCount == 0 {
		return "", 0
	}

	frames := runtime.CallersFrames(f)
	frame, _ := frames.Next()
	return frame.Function, frame.Line
}
