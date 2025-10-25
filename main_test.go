package nabu

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetLogLevel(LevelDebug)
	SetLogOutput(OutputInternal)

	code := m.Run()
	os.Exit(code)
}

// resetTestState resets the internal output buffer before each test
// to ensure test isolation. This function is thread-safe.
func resetTestState() {
	configMutex.Lock()
	defer configMutex.Unlock()
	internalOutput = ""
}

// getInternalOutput safely reads the internal output buffer
func getInternalOutput() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return internalOutput
}
