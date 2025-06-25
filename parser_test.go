package nabu

import (
	"testing"
	"time"
)

func TestParserIncrementalParsing(t *testing.T) {
	cases := [][]string{
		{
			`{"Date":"2025-06-24 18:16:08.835372","Msg":"Sample [A]","Level":1}`,
		},
		{
			`{"Date":"2025-06-24 18:16:08.835372","Msg":"Sample [A]","Level":1}`,
			`{"Date":"2025-06-24 18:16:10.056499","Msg":"Sample [B]","Level":2}`,
		},
		{
			`{"Date":"2025-06-24 18:16:08.835372","Msg":"Sample [A]","Level":1}`,
			`{"Date":"2025-06-24 18:16:10.056499","Msg":"Sample [B]","Level":2}`,
			`{"Date":"2025-06-25 01:26:01.447444","Msg":"Sample [C]","Level":1}`,
		},
		{
			`{"UUID":"a3a17cc1-8481-4454-bff4-1fcec8398d6c","Date":"2025-06-25 01:26:02.408736","Error":"EOF","Function":"github.com/rah-0/nabu/node.(*Node).handleConnection","Line":276,"Level":3}`,
			`{"UUID":"a3a17cc1-8481-4454-bff4-1fcec8398d6c","Date":"2025-06-25 01:26:02.408897","Error":"EOF","Function":"github.com/rah-0/nabu/node.(*Node).handleErrors.func1","Line":396,"Level":4}`,
		},
		{
			`{"Date":"2025-06-24 18:16:08.835372","Msg":"Sample [A]","Level":1}`,
			`{"Date":"2025-06-24 18:16:10.056499","Msg":"Sample [B]","Level":2}`,
			`{"Date":"2025-06-25 01:26:01.447444","Msg":"Sample [C]","Level":1}`,
			`{"UUID":"b29b369e-56b1-49ee-892f-b6f43464e3f1","Date":"2025-06-25 01:26:02.380090","Error":"EOF","Function":"github.com/rah-0/nabu/node.(*Node).handleConnection","Line":276,"Level":3}`,
			`{"UUID":"a3a17cc1-8481-4454-bff4-1fcec8398d6c","Date":"2025-06-25 01:26:02.408736","Error":"EOF","Function":"github.com/rah-0/nabu/node.(*Node).handleConnection","Line":276,"Level":3}`,
			`{"UUID":"a3a17cc1-8481-4454-bff4-1fcec8398d6c","Date":"2025-06-25 01:26:02.408897","Error":"EOF","Function":"github.com/rah-0/nabu/node.(*Node).handleErrors.func1","Line":396,"Level":4}`,
		},
	}

	expectedEntries := []int{1, 2, 3, 0, 3}
	expectedTraces := []int{0, 0, 0, 1, 2}

	for i, lines := range cases {
		parser := NewParser().FromLines(lines)
		parsed := parser.Parse()

		if len(parsed.Entries) != expectedEntries[i] {
			t.Errorf("case %d: expected %d entries, got %d", i, expectedEntries[i], len(parsed.Entries))
		}
		if len(parsed.Traces) != expectedTraces[i] {
			t.Errorf("case %d: expected %d traces, got %d", i, expectedTraces[i], len(parsed.Traces))
		}
	}
}

func TestParserWithAfterDate(t *testing.T) {
	logs := []string{
		`{"Date":"2025-06-24 18:16:08.835372","Msg":"Sample","Level":1}`,
		`{"Date":"2025-06-25 01:26:01.447444","Msg":"Connection open","Level":1}`,
		`{"UUID":"x","Date":"2025-06-25 01:26:02.408736","Error":"EOF","Function":"F","Line":1,"Level":3}`,
		`{"UUID":"x","Date":"2025-06-25 01:26:02.408837","Error":"EOF","Function":"G","Line":2,"Level":4}`,
	}

	cutoff, _ := time.Parse(TimeLayout, "2025-06-25 00:00:00.000000")
	parser := NewParser().FromLines(logs).AfterDate(cutoff)
	parsed := parser.Parse()

	if len(parsed.Entries) != 1 {
		t.Errorf("afterDate: expected 1 entry, got %d", len(parsed.Entries))
	}
	if len(parsed.Traces) != 1 {
		t.Errorf("afterDate: expected 1 trace, got %d", len(parsed.Traces))
	}
	if len(parsed.Traces[0].Frames) != 2 {
		t.Errorf("afterDate: expected 2 frames, got %d", len(parsed.Traces[0].Frames))
	}
}

func TestParserCases(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) {
		parser := NewParser().FromLines([]string{})
		parsed := parser.Parse()
		if len(parsed.Entries) != 0 || len(parsed.Traces) != 0 {
			t.Error("expected no entries or traces for empty input")
		}
	})

	t.Run("Invalid JSON skipped", func(t *testing.T) {
		logs := []string{
			`{"Date":"2025-06-25 01:00:00.000000","Msg":"Valid","Level":1}`,
			`invalid json line`,
			`{"Date":"2025-06-25 01:01:00.000000","Msg":"Also valid","Level":1}`,
		}
		parsed := NewParser().FromLines(logs).Parse()
		if len(parsed.Entries) != 2 {
			t.Errorf("expected 2 valid entries, got %d", len(parsed.Entries))
		}
	})

	t.Run("Duplicate UUID entries form single trace", func(t *testing.T) {
		logs := []string{
			`{"UUID":"z","Date":"2025-06-25 01:01:00.000000","Error":"Some error","Function":"A","Line":1,"Level":3}`,
			`{"UUID":"z","Date":"2025-06-25 01:01:01.000000","Function":"B","Line":2,"Level":3}`,
			`{"UUID":"z","Date":"2025-06-25 01:01:02.000000","Function":"C","Line":3,"Level":3}`,
		}
		parsed := NewParser().FromLines(logs).Parse()
		if len(parsed.Traces) != 1 {
			t.Errorf("expected 1 trace, got %d", len(parsed.Traces))
		}
		if parsed.Traces[0].Error != "Some error" {
			t.Errorf("expected error to be 'Some error', got '%s'", parsed.Traces[0].Error)
		}
		if len(parsed.Traces[0].Frames) != 3 {
			t.Errorf("expected 3 frames, got %d", len(parsed.Traces[0].Frames))
		}
	})

	t.Run("Only invalid entries", func(t *testing.T) {
		logs := []string{
			`not json`,
			`also not json`,
		}
		parsed := NewParser().FromLines(logs).Parse()
		if len(parsed.Entries) != 0 && len(parsed.Traces) != 0 {
			t.Error("expected no entries or traces from invalid input")
		}
	})
}
