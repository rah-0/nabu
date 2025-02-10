package nabu

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getDate() string {
	ct := time.Now().UTC()
	return ct.Format("2006-01-02 15:04:05.000000")
}

func fromJson(jsonString string) *Output {
	jsonString = strings.TrimSpace(jsonString)

	var o Output
	err := json.Unmarshal([]byte(jsonString), &o)
	if err != nil {
		return nil
	}
	return &o
}

func toJson(targetObject any) (output string) {
	b, e := json.Marshal(targetObject)
	if e != nil {
		output = `{"UUID":"` + uuid.NewString() + `","Date":"` + getDate() + `","Error":"` + e.Error() + `","Level":4}`
		return
	}
	output = string(b)
	return
}

func (x *Logger) getFirstTrace() (string, int) {
	programCounters := make([]uintptr, 1)
	callersCount := runtime.Callers(3, programCounters) // Ignore: runtime.Callers, getFirstTrace, Log

	if callersCount == 0 {
		return "", 0
	}

	frames := runtime.CallersFrames(programCounters)
	frame, _ := frames.Next()
	return frame.Function, frame.Line
}
