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

// Posts and comments routinely carry multiline text (the Ordinal OpenAPI
// examples themselves include embedded newlines). Writing those cells
// verbatim into a table splits one row across several terminal lines and
// destroys column alignment. Table output must collapse \n/\r/\t to spaces
// so each logical row stays on a single line — and CSV must leave them
// alone so csv.Writer-quoted content round-trips unchanged.
func TestFormatOutput_TableNormalizesMultilineCells(t *testing.T) {
	items := []map[string]interface{}{
		{"id": "1", "text": "line one\nline two\r\nline three"},
		{"id": "2", "text": "tabbed\tvalue"},
	}
	out, _, err := FormatOutput(items, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect exactly header + separator + 2 data rows = 4 lines. Embedded
	// newlines would have inflated this; that's the regression.
	if got := strings.Count(out, "\n"); got != 3 {
		t.Errorf("expected 3 newlines (1 header + 1 separator + 2 rows - 1 trimmed trailing); got %d in %q", got, out)
	}
	if strings.Contains(out, "line one\n") || strings.Contains(out, "tabbed\t") {
		t.Errorf("table cells must not contain raw newlines/tabs; got %q", out)
	}
	// All segments of the multiline value should still be present on a
	// single line so readers can see the full content. Each control rune
	// becomes one space, so \r\n collapses to two spaces — assert presence
	// of each segment rather than an exact joiner.
	for _, want := range []string{"line one", "line two", "line three", "tabbed value"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output; got %q", want, out)
		}
	}
}

func TestFormatOutput_CSVPreservesMultilineCells(t *testing.T) {
	items := []map[string]interface{}{
		{"id": "1", "text": "line one\nline two"},
	}
	out, _, err := FormatOutput(items, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(strings.NewReader(out))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("parsing csv: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected header + 1 row; got %d", len(records))
	}
	textIdx := -1
	for i, h := range records[0] {
		if h == "text" {
			textIdx = i
		}
	}
	if textIdx < 0 {
		t.Fatalf("text column missing: %v", records[0])
	}
	if got := records[1][textIdx]; got != "line one\nline two" {
		t.Errorf("CSV must preserve embedded newlines for round-tripping; got %q", got)
	}
}

// Emoji and CJK runes occupy two terminal columns in a monospace font. Byte
// length undercounts them, so a table with emoji in one cell and a wider
// plain-ASCII cell beneath it would mis-align — either the emoji row gets
// over-padded (byte-padded past the column) or the column under-reserves
// space. Width accounting must be display-based.
func TestFormatOutput_TableAlignsEmojiAndCJK(t *testing.T) {
	items := []map[string]interface{}{
		{"id": "1", "name": "🚀"},
		{"id": "2", "name": "中文"},
		{"id": "3", "name": "abcd"},
	}
	out, _, err := FormatOutput(items, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(out, "\n")
	// header, separator, 3 rows.
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines; got %d: %q", len(lines), out)
	}
	// Every line must occupy the same number of display columns — that is
	// the whole point of padding. If widths were computed byte-wise, the
	// rows containing multi-byte emoji/CJK would come up short on display
	// columns and the row containing "abcd" would be wider.
	refCols := displayWidth(lines[0])
	for i, l := range lines {
		if w := displayWidth(l); w != refCols {
			t.Errorf("line %d has %d display cols, expected %d: %q", i, w, refCols, l)
		}
	}
}

func TestDisplayWidth(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abc", 3},
		{"🚀", 2},
		{"中文", 4},
		{"café", 4},         // 'é' precomposed: narrow
		{"cafe\u0301", 4},   // 'e' + combining acute: combining mark is 0-width
		{"a\u200Db", 2},     // ZWJ is 0-width
		{"\u270C\uFE0F", 1}, // victory hand U+270C (Neutral kind here) + VS-16 — VS is 0; base is narrow per EAW
		{"hi\tthere", 7},    // control char (\t) treated as 0; we sanitize before measuring anyway
	}
	for _, c := range cases {
		if got := displayWidth(c.in); got != c.want {
			t.Errorf("displayWidth(%q) = %d; want %d", c.in, got, c.want)
		}
	}
}

// A 204-style empty-body response from a read endpoint can manifest as an
// empty JSON object after a round-trip through json.Unmarshal. The formatter
// must not error for table/csv rendering; it should simply say "No results"
// so the CLI can pass {} through unchanged for reads without a format-
// dependent crash.
func TestFormatOutput_EmptyObjectRendersAsNoResults(t *testing.T) {
	empty := map[string]interface{}{}
	out, footer, err := FormatOutput(empty, FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results'; got %q", out)
	}
	if footer != "" {
		t.Errorf("expected no footer; got %q", footer)
	}

	out, footer, err = FormatOutput(empty, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty csv body; got %q", out)
	}
	if footer != "" {
		t.Errorf("expected no footer; got %q", footer)
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
