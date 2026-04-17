package output

import (
	"strings"
	"testing"
)

func TestFormatOutput_JSON(t *testing.T) {
	out, err := FormatOutput(map[string]int{"a": 1}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"a": 1`) {
		t.Errorf("expected pretty JSON, got %q", out)
	}
}

func TestFormatOutput_Table(t *testing.T) {
	items := []map[string]interface{}{
		{"id": "1", "name": "a"},
		{"id": "2", "name": "b"},
	}
	out, err := FormatOutput(items, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ID") || !strings.Contains(out, "NAME") {
		t.Errorf("expected headers, got %q", out)
	}
	if !strings.Contains(out, "1") || !strings.Contains(out, "2") {
		t.Errorf("expected row values, got %q", out)
	}
}

func TestFormatOutput_CSV(t *testing.T) {
	items := []map[string]interface{}{{"id": "1", "name": "a"}}
	out, err := FormatOutput(items, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected header + row, got %q", out)
	}
	if !strings.Contains(lines[0], "id") || !strings.Contains(lines[0], "name") {
		t.Errorf("expected header line, got %q", lines[0])
	}
}

func TestFormatOutput_EmptyTable(t *testing.T) {
	out, err := FormatOutput([]map[string]interface{}{}, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results', got %q", out)
	}
}
