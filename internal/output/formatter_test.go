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

// Wrapper list responses (items + pagination metadata) used to render as a
// single row with columns like POSTS/NEXTCURSOR/HASMORE. They should unwrap to
// the contained items so table/csv output is actually useful for list commands.
func TestFormatOutput_UnwrapsListEnvelope_Table(t *testing.T) {
	envelope := map[string]interface{}{
		"posts": []map[string]interface{}{
			{"id": "1", "title": "a"},
			{"id": "2", "title": "b"},
		},
		"nextCursor": "cur",
		"hasMore":    true,
	}
	out, err := FormatOutput(envelope, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ID") || !strings.Contains(out, "TITLE") {
		t.Errorf("expected item headers (ID/TITLE), got %q", out)
	}
	if strings.Contains(out, "NEXTCURSOR") || strings.Contains(out, "HASMORE") || strings.Contains(out, "POSTS") {
		t.Errorf("envelope metadata should not appear as columns; got %q", out)
	}
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Errorf("expected row values a,b; got %q", out)
	}
}

func TestFormatOutput_UnwrapsListEnvelope_CSV(t *testing.T) {
	envelope := map[string]interface{}{
		"labels": []map[string]interface{}{
			{"id": "1", "name": "x"},
		},
	}
	out, err := FormatOutput(envelope, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected header + 1 row, got %q", out)
	}
	if !strings.Contains(lines[0], "id") || !strings.Contains(lines[0], "name") {
		t.Errorf("expected item headers, got %q", lines[0])
	}
}

func TestFormatOutput_EmptyEnvelope(t *testing.T) {
	envelope := map[string]interface{}{
		"posts":      []map[string]interface{}{},
		"nextCursor": "",
		"hasMore":    false,
	}
	out, err := FormatOutput(envelope, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results' for empty envelope; got %q", out)
	}
}

func TestFormatOutput_SingleResourceNotUnwrapped(t *testing.T) {
	// A single-object response (e.g. GET /post/:id) has no array child and
	// must remain a one-row table with its own fields, not an error.
	single := map[string]interface{}{
		"id":    "1",
		"title": "hello",
	}
	out, err := FormatOutput(single, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ID") || !strings.Contains(out, "TITLE") {
		t.Errorf("expected single-row table, got %q", out)
	}
}
