package nabu

import (
	"testing"
)

func TestShouldLog(t *testing.T) {
	levels := []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal}

	// Test with default LevelDebug (should allow all levels)
	SetLogLevel(LevelDebug)
	for _, level := range levels {
		if !shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return true at LevelDebug, but got false", level)
		}
	}

	// Test with LevelInfo (should block Debug)
	SetLogLevel(LevelInfo)
	if shouldLog(LevelDebug) {
		t.Errorf("Expected shouldLog(LevelDebug) to return false at LevelInfo, but got true")
	}
	for _, level := range levels[1:] { // Info and above should be allowed
		if !shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return true at LevelInfo, but got false", level)
		}
	}

	// Test with LevelWarn (should block Debug & Info)
	SetLogLevel(LevelWarn)
	for _, level := range []LogLevel{LevelDebug, LevelInfo} {
		if shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return false at LevelWarn, but got true", level)
		}
	}
	for _, level := range levels[2:] { // Warn and above should be allowed
		if !shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return true at LevelWarn, but got false", level)
		}
	}

	// Test with LevelError (should block Debug, Info, and Warn)
	SetLogLevel(LevelError)
	for _, level := range []LogLevel{LevelDebug, LevelInfo, LevelWarn} {
		if shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return false at LevelError, but got true", level)
		}
	}
	for _, level := range levels[3:] { // Error and Fatal should be allowed
		if !shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return true at LevelError, but got false", level)
		}
	}

	// Test with LevelFatal (should only allow Fatal)
	SetLogLevel(LevelFatal)
	for _, level := range levels[:4] { // All except Fatal should be blocked
		if shouldLog(level) {
			t.Errorf("Expected shouldLog(%v) to return false at LevelFatal, but got true", level)
		}
	}
	if !shouldLog(LevelFatal) {
		t.Errorf("Expected shouldLog(LevelFatal) to return true at LevelFatal, but got false")
	}
}
