package output

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestFormatOutput_UnknownFormatErrors(t *testing.T) {
	_, _, err := FormatOutput(map[string]int{"a": 1}, Format("yaml"))
	if err == nil {
		t.Fatalf("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "invalid output format") {
		t.Errorf("expected invalid output format message, got %v", err)
	}
}

func TestIsValidFormat(t *testing.T) {
	for _, f := range []Format{FormatJSON, FormatTable, FormatCSV} {
		if !IsValidFormat(f) {
			t.Errorf("expected %q to be valid", f)
		}
	}
	for _, f := range []Format{"", "yaml", "JSON", "jsonl"} {
		if IsValidFormat(f) {
			t.Errorf("expected %q to be invalid", f)
		}
	}
}

func TestFormatOutput_JSON(t *testing.T) {
	out, footer, err := FormatOutput(map[string]int{"a": 1}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if footer != "" {
		t.Errorf("expected no footer for JSON; got %q", footer)
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
	out, _, err := FormatOutput(items, FormatTable)
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
	out, _, err := FormatOutput(items, FormatCSV)
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
	out, _, err := FormatOutput([]map[string]interface{}{}, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results', got %q", out)
	}
}

// Wrapper list responses (items + pagination metadata) used to render as a
// single row with columns like POSTS/NEXTCURSOR/HASMORE. They should unwrap to
// the contained items so table/csv output is actually useful for list
// commands, and their pagination metadata should surface as a footer so
// callers can paginate without switching to JSON.
func TestFormatOutput_UnwrapsListEnvelope_Table(t *testing.T) {
	envelope := map[string]interface{}{
		"posts": []map[string]interface{}{
			{"id": "1", "title": "a"},
			{"id": "2", "title": "b"},
		},
		"nextCursor": "cur",
		"hasMore":    true,
	}
	out, footer, err := FormatOutput(envelope, FormatTable)
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
	if !strings.Contains(footer, "nextCursor: cur") {
		t.Errorf("expected footer to carry nextCursor; got %q", footer)
	}
	if !strings.Contains(footer, "hasMore: true") {
		t.Errorf("expected footer to carry hasMore; got %q", footer)
	}
}

func TestFormatOutput_UnwrapsListEnvelope_CSV(t *testing.T) {
	envelope := map[string]interface{}{
		"labels": []map[string]interface{}{
			{"id": "1", "name": "x"},
		},
	}
	out, footer, err := FormatOutput(envelope, FormatCSV)
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
	// A pagination-less envelope produces no footer.
	if footer != "" {
		t.Errorf("expected no footer for pagination-less envelope; got %q", footer)
	}
}

// CSV output for paginated list responses must stay parseable (pure rows on
// stdout) while the pagination cursor is still surfaced via the footer.
func TestFormatOutput_UnwrapsListEnvelope_CSV_Paginated(t *testing.T) {
	envelope := map[string]interface{}{
		"posts": []map[string]interface{}{
			{"id": "1", "title": "a"},
		},
		"nextCursor": "cur",
		"hasMore":    true,
	}
	out, footer, err := FormatOutput(envelope, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "nextCursor") || strings.Contains(out, "hasMore") {
		t.Errorf("CSV body must not include pagination metadata; got %q", out)
	}
	if !strings.Contains(footer, "nextCursor: cur") || !strings.Contains(footer, "hasMore: true") {
		t.Errorf("expected pagination footer; got %q", footer)
	}
}

func TestFormatOutput_EmptyEnvelope(t *testing.T) {
	envelope := map[string]interface{}{
		"posts":      []map[string]interface{}{},
		"nextCursor": "",
		"hasMore":    false,
	}
	out, footer, err := FormatOutput(envelope, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results' for empty envelope; got %q", out)
	}
	// Zero-value pagination fields (final page) should not produce noise.
	if footer != "" {
		t.Errorf("expected empty footer for terminal page; got %q", footer)
	}
}

// Nested composite cell values (maps, slices) must serialize as compact JSON
// so table/CSV output stays well-formed and stable. fmt.Sprintf("%v", v) on a
// map produces Go-style literals like "map[k:v]" with unstable key ordering,
// which breaks diffing and downstream parsing for fields like webhook headers
// or idea labels.
func TestFormatOutput_NestedMapRendersAsJSON(t *testing.T) {
	items := []map[string]interface{}{
		{
			"id":      "1",
			"headers": map[string]interface{}{"X-Key": "value", "X-Other": "v2"},
		},
	}
	out, _, err := FormatOutput(items, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "map[") {
		t.Errorf("nested map should not render as Go literal; got %q", out)
	}
	if !strings.Contains(out, `{"X-Key":"value","X-Other":"v2"}`) {
		t.Errorf("expected compact JSON with sorted keys; got %q", out)
	}
}

func TestFormatOutput_NestedSliceRendersAsJSON_CSV(t *testing.T) {
	items := []map[string]interface{}{
		{
			"id": "1",
			"labels": []interface{}{
				map[string]interface{}{"id": "lbl_1", "name": "priority"},
				map[string]interface{}{"id": "lbl_2", "name": "growth"},
			},
		},
	}
	out, _, err := FormatOutput(items, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "map[") {
		t.Errorf("nested slice should not render as Go literal; got %q", out)
	}
	// The CSV writer double-quotes cells containing quotes, so compare the
	// parsed cell value (quotes unescaped) rather than raw CSV bytes.
	r := csv.NewReader(strings.NewReader(out))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("parsing csv: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected header + 1 row; got %d", len(records))
	}
	labelsIdx := -1
	for i, h := range records[0] {
		if h == "labels" {
			labelsIdx = i
		}
	}
	if labelsIdx < 0 {
		t.Fatalf("labels header missing: %v", records[0])
	}
	want := `[{"id":"lbl_1","name":"priority"},{"id":"lbl_2","name":"growth"}]`
	if got := records[1][labelsIdx]; got != want {
		t.Errorf("expected labels cell %q; got %q", want, got)
	}
}

func TestFormatOutput_SingleResourceNotUnwrapped(t *testing.T) {
	// A single-object response (e.g. GET /post/:id) has no array child and
	// must remain a one-row table with its own fields, not an error.
	single := map[string]interface{}{
		"id":    "1",
		"title": "hello",
	}
	out, _, err := FormatOutput(single, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ID") || !strings.Contains(out, "TITLE") {
		t.Errorf("expected single-row table, got %q", out)
	}
}

// A single resource like an Idea can have scalar fields (id, title, status)
// plus an array-of-objects field (labels). It must not be mistaken for a list
// envelope of labels — the row should show the idea itself, not its labels.
func TestFormatOutput_SingleResourceWithArrayFieldNotUnwrapped(t *testing.T) {
	idea := map[string]interface{}{
		"id":     "idea_1",
		"title":  "My idea",
		"status": "draft",
		"labels": []map[string]interface{}{
			{"id": "lbl_1", "name": "priority"},
			{"id": "lbl_2", "name": "growth"},
		},
	}
	out, footer, err := FormatOutput(idea, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "idea_1") || !strings.Contains(out, "My idea") {
		t.Errorf("expected the idea's own fields to render; got %q", out)
	}
	// Must be a single-row table (headers + separator + 1 data row = 3 lines),
	// not one row per label.
	if lines := strings.Split(out, "\n"); len(lines) != 3 {
		t.Errorf("expected 3 lines (headers, separator, one row); got %d: %q", len(lines), out)
	}
	if footer != "" {
		t.Errorf("single resource should produce no pagination footer; got %q", footer)
	}
}
