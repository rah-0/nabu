package nabu

import (
	"strings"
	"testing"
)

func TestToJsonError(t *testing.T) {
	type UnmarshalableStruct struct {
		Ch chan int
	}

	unmarshalable := UnmarshalableStruct{
		Ch: make(chan int),
	}

	result := toJson(unmarshalable)

	if !strings.Contains(result, "json: unsupported type: chan int") {
		t.Errorf("Expected error message, got: %s", result)
	}
}
