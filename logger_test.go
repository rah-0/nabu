package nabu

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	SetLogLevel(LevelDebug)
	SetLogOutput(OutputInternal)

	code := m.Run()
	os.Exit(code)
}

func TestFromErrorSimple(t *testing.T) {
	val := "10.11"

	e := getInt64FromString(val)
	if e == nil {
		t.Error()
	}

	FromError(e).WithArgs(val).Log()

	expectedError := `strconv.ParseInt: parsing "` + val + `": invalid syntax`
	expectedLine := 28
	expectedFunction := "github.com/rah-0/nabu.TestFromErrorSimple"
	expectedArgs := []any{val}
	validateLogOutput(t, internalOutput, expectedError, expectedLine, expectedArgs, expectedFunction)
	internalOutput = ""
}

func TestFromErrorDepth(t *testing.T) {
	arg := "Init_"
	e := a(arg)
	FromError(e).WithArgs(arg).Log()

	logEntries := strings.Split(strings.TrimSpace(internalOutput), "\n")

	expectedFunctions := []string{
		"github.com/rah-0/nabu.d",
		"github.com/rah-0/nabu.c",
		"github.com/rah-0/nabu.b",
		"github.com/rah-0/nabu.a",
		"github.com/rah-0/nabu.TestFromErrorDepth",
	}
	expectedArgs := []any{
		[]any{"Init_ABCD"},
		[]any{"Init_ABC"},
		[]any{"Init_AB"},
		[]any{"Init_A"},
		[]any{"Init_"},
	}
	expectedLineNumbers := []int{91, 86, 81, 76, 41}
	expectedError := "testError"

	for i, rawLog := range logEntries {
		rawLog = strings.TrimSpace(rawLog)
		if rawLog == "" {
			continue
		}

		validateLogOutput(t, rawLog, expectedError, expectedLineNumbers[i], expectedArgs[i], expectedFunctions[i])
	}
	internalOutput = ""
}

func a(s string) error {
	s += "A"
	e := b(s)
	return FromError(e).WithArgs(s).WithMessage("A").Log()
}
func b(s string) error {
	s += "B"
	e := c(s)
	return FromError(e).WithArgs(s).WithMessage("B").Log()
}
func c(s string) error {
	s += "C"
	e := d(s)
	return FromError(e).WithArgs(s).WithMessage("C").Log()
}
func d(s string) error {
	s += "D"
	e := errors.New("testError")
	return FromError(e).WithArgs(s).WithMessage("D").Log()
}

func TestFromMessageSimple(t *testing.T) {
	msg := "This is a test message"

	FromMessage(msg).WithArgs("test_arg").Log()

	expectedError := ""
	expectedLine := 97
	expectedFunction := "github.com/rah-0/nabu.TestFromMessageSimple"
	expectedArgs := []any{"test_arg"}

	validateLogOutput(t, internalOutput, expectedError, expectedLine, expectedArgs, expectedFunction)
	internalOutput = ""
}

func getInt64FromString(input string) error {
	var e error
	_, e = strconv.ParseInt(input, 10, 64)
	return e
}

func validateLogOutput(t *testing.T, rawLog string, expectedError string, expectedLine int, expectedArgs any, expectedFunction string) {
	t.Helper()

	result := fromJson(rawLog)
	if result == nil {
		t.Error("Expected valid JSON log output, but got nil")
		return
	}

	if result.Error != expectedError {
		t.Errorf("Expected error message: %q, but got: %q", expectedError, result.Error)
	}

	if result.Line != expectedLine {
		t.Errorf("Expected line number: %d, but got: %d", expectedLine, result.Line)
	}

	if !reflect.DeepEqual(result.Args, expectedArgs) {
		t.Errorf("Expected Args: %v, but got: %v", expectedArgs, result.Args)
	}

	if result.Function != expectedFunction {
		t.Errorf("Expected function: %q, but got: %q", expectedFunction, result.Function)
	}
}

func BenchmarkFromErrorSimple(b *testing.B) {
	val := "10.11"
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		e := getInt64FromString(val)
		if e == nil {
			b.Fatal("Expected an error but got nil")
		}

		FromError(e).WithArgs(val).Log()
		internalOutput = ""
	}
}

func BenchmarkFromErrorDepth(b *testing.B) {
	arg := "Init_"
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		e := a(arg)
		if e == nil {
			b.Fatal("Expected an error but got nil")
		}

		FromError(e).WithArgs(arg).Log()
		internalOutput = ""
	}
}

func BenchmarkFromErrorSimpleNil(b *testing.B) {
	val := "1234"
	e := getInt64FromString(val)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FromError(e).WithArgs(val).Log()
		internalOutput = ""
	}
}
