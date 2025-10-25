package nabu

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFromErrorSimple(t *testing.T) {
	resetTestState()
	val := "10.11"

	e := getInt64FromString(val)
	if e == nil {
		t.Error()
	}

	FromError(e).WithArgs(val).Log()

	expectedError := `strconv.ParseInt: parsing "` + val + `": invalid syntax`
	expectedLine := 20
	expectedFunction := "github.com/rah-0/nabu.TestFromErrorSimple"
	expectedArgs := []any{val}
	validateLogOutput(t, getInternalOutput(), expectedError, expectedLine, expectedArgs, expectedFunction)
}

func TestFromErrorDepth(t *testing.T) {
	resetTestState()
	arg := "Init_"
	e := a(arg)
	FromError(e).WithArgs(arg).Log()

	logEntries := strings.Split(strings.TrimSpace(getInternalOutput()), "\n")

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
	expectedLineNumbers := []int{84, 79, 74, 69, 33}
	// Only the first log (where error occurs) shows the error
	// Subsequent logs in the chain only show their Msg field
	expectedErrors := []string{"testError", "", "", "", ""}

	for i, rawLog := range logEntries {
		rawLog = strings.TrimSpace(rawLog)
		if rawLog == "" {
			continue
		}

		validateLogOutput(t, rawLog, expectedErrors[i], expectedLineNumbers[i], expectedArgs[i], expectedFunctions[i])
	}
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
	resetTestState()
	msg := "This is a test message"

	FromMessage(msg).WithArgs("test_arg").EnableStackTrace().Log()

	expectedError := ""
	expectedLine := 91
	expectedFunction := "github.com/rah-0/nabu.TestFromMessageSimple"
	expectedArgs := []any{"test_arg"}

	validateLogOutput(t, getInternalOutput(), expectedError, expectedLine, expectedArgs, expectedFunction)
}

func TestFromError_Is(t *testing.T) {
	err := makeWrapped()
	if !errors.Is(err, ErrSentinel) {
		t.Fatalf("expected errors.Is to match ErrSentinel, got %v", err)
	}
}

func TestWithUuid(t *testing.T) {
	resetTestState()
	customUUID := "custom-uuid-12345"

	err := errors.New("test error")
	FromError(err).WithUuid(customUUID).WithMessage("test message").Log()

	output := getInternalOutput()
	if !strings.Contains(output, customUUID) {
		t.Errorf("Expected output to contain custom UUID %s, but got: %s", customUUID, output)
	}

	// Verify the UUID is exactly what we set
	logEntry := fromJson(output)
	if logEntry == nil {
		t.Fatal("Failed to parse log output")
	}
	if logEntry.UUID != customUUID {
		t.Errorf("Expected UUID to be %s, but got: %s", customUUID, logEntry.UUID)
	}
}

func TestMixedStack(t *testing.T) {
	resetTestState()

	// Test 1: Msg -> Error -> Msg
	t.Run("Msg_Error_Msg", func(t *testing.T) {
		resetTestState()
		_ = FromMessage("starting operation").Log()
		err := errors.New("network timeout")
		log1 := FromError(err).WithMessage("failed to connect").Log()
		_ = FromError(log1).WithMessage("retrying operation").Log()

		output := getInternalOutput()
		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 3 {
			t.Fatalf("Expected 3 log entries, got %d", len(lines))
		}

		// Parse each log
		entry1 := fromJson(lines[0])
		entry2 := fromJson(lines[1])
		entry3 := fromJson(lines[2])

		// Check first log (Msg only, has UUID for correlation)
		if entry1.Msg != "starting operation" {
			t.Errorf("Entry 1: Expected Msg='starting operation', got '%s'", entry1.Msg)
		}
		if entry1.Error != "" {
			t.Errorf("Entry 1: Expected no Error, got '%s'", entry1.Error)
		}
		if entry1.UUID == "" {
			t.Error("Entry 1: Expected UUID to be generated for message")
		}

		// Check second log (Error occurs here, different UUID)
		if entry2.Error != "network timeout" {
			t.Errorf("Entry 2: Expected Error='network timeout', got '%s'", entry2.Error)
		}
		if entry2.Msg != "failed to connect" {
			t.Errorf("Entry 2: Expected Msg='failed to connect', got '%s'", entry2.Msg)
		}
		if entry2.UUID == "" {
			t.Error("Entry 2: Expected UUID to be generated")
		}
		if entry2.UUID == entry1.UUID {
			t.Error("Entry 2: Expected different UUID from entry 1 (different chain)")
		}

		// Check third log (wraps Logger, should have same UUID as entry2, no Error duplication)
		if entry3.Error != "" {
			t.Errorf("Entry 3: Expected no Error (already logged), got '%s'", entry3.Error)
		}
		if entry3.Msg != "retrying operation" {
			t.Errorf("Entry 3: Expected Msg='retrying operation', got '%s'", entry3.Msg)
		}
		if entry3.UUID != entry2.UUID {
			t.Errorf("Entry 3: Expected same UUID as entry 2 (%s), got '%s'", entry2.UUID, entry3.UUID)
		}
	})

	// Test 2: Error -> Msg -> Error
	t.Run("Error_Msg_Error", func(t *testing.T) {
		resetTestState()
		err1 := errors.New("database error")
		log1 := FromError(err1).WithMessage("query failed").Log()
		_ = FromError(log1).WithMessage("processing data").Log()
		err2 := errors.New("validation error")
		_ = FromError(err2).WithMessage("invalid input").Log()

		output := getInternalOutput()
		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 3 {
			t.Fatalf("Expected 3 log entries, got %d", len(lines))
		}

		entry1 := fromJson(lines[0])
		entry2 := fromJson(lines[1])
		entry3 := fromJson(lines[2])

		// First error chain
		if entry1.Error != "database error" {
			t.Errorf("Entry 1: Expected Error='database error', got '%s'", entry1.Error)
		}
		if entry1.UUID == "" {
			t.Error("Entry 1: Expected UUID to be generated")
		}

		// Wrapping first error
		if entry2.Error != "" {
			t.Errorf("Entry 2: Expected no Error, got '%s'", entry2.Error)
		}
		if entry2.UUID != entry1.UUID {
			t.Errorf("Entry 2: Expected same UUID as entry 1, got different UUIDs")
		}

		// New error (should have different UUID)
		if entry3.Error != "validation error" {
			t.Errorf("Entry 3: Expected Error='validation error', got '%s'", entry3.Error)
		}
		if entry3.UUID == "" {
			t.Error("Entry 3: Expected UUID to be generated")
		}
		if entry3.UUID == entry1.UUID {
			t.Error("Entry 3: Expected different UUID from entry 1 (new error chain)")
		}
	})

	// Test 3: Msg -> Error -> Error (chaining two different errors)
	t.Run("Msg_Error_Error", func(t *testing.T) {
		resetTestState()
		_ = FromMessage("initialization").Log()
		err1 := errors.New("config error")
		log1 := FromError(err1).WithMessage("failed to load config").Log()
		_ = FromError(log1).WithMessage("startup failed").Log()

		output := getInternalOutput()
		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 3 {
			t.Fatalf("Expected 3 log entries, got %d", len(lines))
		}

		entry1 := fromJson(lines[0])
		entry2 := fromJson(lines[1])
		entry3 := fromJson(lines[2])

		// Msg only (has UUID)
		if entry1.UUID == "" {
			t.Error("Entry 1: Expected UUID for message log")
		}

		// First error (different UUID from message)
		if entry2.Error != "config error" {
			t.Errorf("Entry 2: Expected Error='config error', got '%s'", entry2.Error)
		}
		if entry2.UUID == "" {
			t.Error("Entry 2: Expected UUID to be generated")
		}
		if entry2.UUID == entry1.UUID {
			t.Error("Entry 2: Expected different UUID from entry 1 (different chain)")
		}

		// Wrapping the error (same UUID as error)
		if entry3.Error != "" {
			t.Errorf("Entry 3: Expected no Error, got '%s'", entry3.Error)
		}
		if entry3.UUID != entry2.UUID {
			t.Errorf("Entry 3: Expected same UUID as entry 2")
		}
	})

	// Test 4: Complex mixed stack
	t.Run("Complex_Mixed", func(t *testing.T) {
		resetTestState()
		_ = FromMessage("step 1").Log()
		err1 := errors.New("error A")
		log1 := FromError(err1).WithMessage("handling error A").Log()
		_ = FromError(log1).WithMessage("step 2").Log()
		_ = FromMessage("step 3").Log()
		err2 := errors.New("error B")
		_ = FromError(err2).WithMessage("handling error B").Log()

		output := getInternalOutput()
		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 5 {
			t.Fatalf("Expected 5 log entries, got %d", len(lines))
		}

		entries := make([]*Output, 5)
		for i := 0; i < 5; i++ {
			entries[i] = fromJson(lines[i])
		}

		// Entry 0: Msg only, has UUID
		uuid0 := entries[0].UUID
		if uuid0 == "" {
			t.Error("Entry 0: Expected UUID for message")
		}

		// Entry 1: First error, UUID generated (different from message)
		if entries[1].Error != "error A" {
			t.Errorf("Entry 1: Expected Error='error A', got '%s'", entries[1].Error)
		}
		uuid1 := entries[1].UUID
		if uuid1 == "" {
			t.Error("Entry 1: Expected UUID")
		}
		if uuid1 == uuid0 {
			t.Error("Entry 1: Expected different UUID from entry 0 (different chain)")
		}

		// Entry 2: Wraps error, same UUID as entry 1, no Error
		if entries[2].Error != "" {
			t.Errorf("Entry 2: Expected no Error, got '%s'", entries[2].Error)
		}
		if entries[2].UUID != uuid1 {
			t.Error("Entry 2: Expected same UUID as entry 1")
		}

		// Entry 3: New message, has UUID (different from previous chains)
		uuid3 := entries[3].UUID
		if uuid3 == "" {
			t.Error("Entry 3: Expected UUID for message")
		}
		if uuid3 == uuid0 || uuid3 == uuid1 {
			t.Error("Entry 3: Expected different UUID from previous entries (new chain)")
		}

		// Entry 4: New error, new UUID (different from all previous)
		if entries[4].Error != "error B" {
			t.Errorf("Entry 4: Expected Error='error B', got '%s'", entries[4].Error)
		}
		uuid2 := entries[4].UUID
		if uuid2 == "" {
			t.Error("Entry 4: Expected UUID")
		}
		if uuid2 == uuid0 || uuid2 == uuid1 || uuid2 == uuid3 {
			t.Error("Entry 4: Expected different UUID from all previous entries (new error chain)")
		}
	})
}

var (
	ErrSentinel = errors.New("sentinel error")
)

func makeWrapped() error {
	// Wrap a sentinel error so we can test errors.Is
	return FromError(ErrSentinel).WithMessage("wrapped sentinel").Log()
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
	resetTestState()
	val := "10.11"
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		e := getInt64FromString(val)
		if e == nil {
			b.Fatal("Expected an error but got nil")
		}

		FromError(e).WithArgs(val).Log()
		resetTestState()
	}
}

func BenchmarkFromErrorDepth(b *testing.B) {
	resetTestState()
	arg := "Init_"
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		e := a(arg)
		if e == nil {
			b.Fatal("Expected an error but got nil")
		}

		FromError(e).WithArgs(arg).Log()
		resetTestState()
	}
}

func BenchmarkFromErrorSimpleNil(b *testing.B) {
	resetTestState()
	val := "1234"
	e := getInt64FromString(val)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FromError(e).WithArgs(val).Log()
		resetTestState()
	}
}


