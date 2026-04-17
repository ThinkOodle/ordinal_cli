// Package output handles formatting CLI output as JSON, table, or CSV.
package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Format represents an output format.
type Format string

const (
	// FormatJSON outputs data as pretty-printed JSON.
	FormatJSON Format = "json"

	// FormatTable outputs data as an aligned text table.
	FormatTable Format = "table"

	// FormatCSV outputs data as comma-separated values.
	FormatCSV Format = "csv"
)

// IsValidFormat reports whether f is one of the supported output formats.
// Callers should validate user-supplied format strings at input time so a
// typo in --output, ORDINAL_OUTPUT_FORMAT, or the config file fails fast
// instead of silently degrading to a default.
func IsValidFormat(f Format) bool {
	switch f {
	case FormatJSON, FormatTable, FormatCSV:
		return true
	}
	return false
}

// FormatOutput formats data according to the given format. Returns the main
// content plus an optional footer string. When a list-envelope response is
// unwrapped for table/csv rendering, the footer carries the pagination
// metadata (e.g. "nextCursor: abc, hasMore: true") so callers can surface it
// separately from the parseable body (typically to stderr for csv, below the
// table for human-readable formats).
func FormatOutput(data interface{}, format Format) (string, string, error) {
	switch format {
	case FormatTable:
		return formatTable(data)
	case FormatCSV:
		return formatCSV(data)
	case FormatJSON:
		out, err := formatJSON(data)
		return out, "", err
	default:
		return "", "", fmt.Errorf("invalid output format %q: must be one of json, table, csv", format)
	}
}

// formatJSON returns pretty-printed JSON.
func formatJSON(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling json: %w", err)
	}
	return string(b), nil
}

// formatTable formats data as an aligned text table. The second return value
// is a non-empty footer string when list-envelope pagination metadata was
// unwrapped alongside the rows.
func formatTable(data interface{}) (string, string, error) {
	rows, headers, meta, err := extractRows(data)
	if err != nil {
		return "", "", err
	}
	if len(rows) == 0 && len(headers) == 0 {
		return "No results", formatFooter(meta), nil
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var buf bytes.Buffer

	for i, h := range headers {
		if i > 0 {
			buf.WriteString("  ")
		}
		buf.WriteString(padRight(strings.ToUpper(h), widths[i]))
	}
	buf.WriteString("\n")

	for i, w := range widths {
		if i > 0 {
			buf.WriteString("  ")
		}
		buf.WriteString(strings.Repeat("-", w))
	}
	buf.WriteString("\n")

	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				buf.WriteString("  ")
			}
			buf.WriteString(padRight(cell, widths[i]))
		}
		buf.WriteString("\n")
	}

	return strings.TrimRight(buf.String(), "\n"), formatFooter(meta), nil
}

// formatCSV formats data as CSV. The second return value is a non-empty
// footer string when list-envelope pagination metadata was unwrapped; it is
// kept out of the CSV body so the body stays strictly parseable.
func formatCSV(data interface{}) (string, string, error) {
	rows, headers, meta, err := extractRows(data)
	if err != nil {
		return "", "", err
	}
	if len(rows) == 0 && len(headers) == 0 {
		return "", formatFooter(meta), nil
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return "", "", fmt.Errorf("writing csv header: %w", err)
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return "", "", fmt.Errorf("writing csv row: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", "", fmt.Errorf("flushing csv: %w", err)
	}

	return strings.TrimRight(buf.String(), "\n"), formatFooter(meta), nil
}

// paginationSiblingKeys are the only scalar field names tolerated alongside
// the items array when unwrapping a list envelope. Any other scalar sibling
// signals that the object is a single resource (which can have scalar fields
// like id/title/status plus an array field like labels), not a list envelope.
var paginationSiblingKeys = map[string]bool{
	"nextCursor": true,
	"hasMore":    true,
	"total":      true,
	"totalCount": true,
	"count":      true,
}

// extractRows converts data into string rows and headers for table/CSV output.
// Supports slices of maps, slices of structs, single maps, and single structs.
//
// For list-style wrapper responses — a JSON object whose only fields are a
// single array-of-objects plus (optionally) a fixed set of pagination keys
// like nextCursor/hasMore/total (e.g. {"posts":[...], "nextCursor":"...",
// "hasMore":true} or the bare {"labels":[...]}) — the array is unwrapped so
// its elements become the rows and the pagination siblings are returned as
// footer metadata. Objects with other scalar siblings (id, title, etc.) are
// treated as single resources and rendered as a one-row table.
func extractRows(data interface{}) ([][]string, []string, map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("marshaling data: %w", err)
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(b, &items); err == nil {
		if len(items) == 0 {
			return nil, nil, nil, nil
		}
		rows, headers, err := mapSliceToRows(items)
		return rows, headers, nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(b, &item); err == nil && len(item) > 0 {
		if unwrapped, meta, ok := unwrapListEnvelope(item); ok {
			if len(unwrapped) == 0 {
				return nil, nil, meta, nil
			}
			rows, headers, err := mapSliceToRows(unwrapped)
			return rows, headers, meta, err
		}
		rows, headers, err := mapSliceToRows([]map[string]interface{}{item})
		return rows, headers, nil, err
	}

	return nil, nil, nil, fmt.Errorf("unsupported data type for table/csv output: %s", reflect.TypeOf(data))
}

// unwrapListEnvelope detects a list-response envelope — a JSON object whose
// only fields are exactly one array-of-objects plus, at most, recognized
// pagination siblings (nextCursor, hasMore, total, totalCount, count). When
// matched, returns the array contents as rows and the pagination scalars as
// metadata so callers can preserve pagination info alongside the rows.
//
// Any unrecognized scalar sibling (id, title, status, …) or nested object
// means this is a single resource, not a list envelope, and the function
// returns ok=false so the caller renders it as a single row.
func unwrapListEnvelope(item map[string]interface{}) ([]map[string]interface{}, map[string]interface{}, bool) {
	var (
		arrayKey   string
		arrayValue []interface{}
		arrayCount int
	)
	meta := make(map[string]interface{})
	for k, v := range item {
		switch vv := v.(type) {
		case []interface{}:
			arrayKey = k
			arrayValue = vv
			arrayCount++
		case map[string]interface{}:
			return nil, nil, false
		default:
			if !paginationSiblingKeys[k] {
				return nil, nil, false
			}
			meta[k] = v
		}
	}
	if arrayCount != 1 || arrayKey == "" {
		return nil, nil, false
	}

	items := make([]map[string]interface{}, 0, len(arrayValue))
	for _, entry := range arrayValue {
		m, ok := entry.(map[string]interface{})
		if !ok {
			return nil, nil, false
		}
		items = append(items, m)
	}
	return items, meta, true
}

// formatFooter renders pagination metadata as a stable "key: value" string
// (keys sorted) for display beneath a table or on stderr for CSV. Returns an
// empty string when there is no metadata to show, or when all metadata values
// are falsy zero-values (empty string, false, 0) — those carry no information
// for the user.
func formatFooter(meta map[string]interface{}) string {
	if len(meta) == 0 {
		return ""
	}
	keys := make([]string, 0, len(meta))
	for k, v := range meta {
		if isZeroValue(v) {
			continue
		}
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s: %v", k, meta[k])
	}
	return strings.Join(parts, "  ")
}

// isZeroValue reports whether v is a zero-value scalar (empty string, false,
// or numeric 0). Used to suppress meaningless footer entries like
// hasMore: false or nextCursor: "" on the final page of a listing.
func isZeroValue(v interface{}) bool {
	switch vv := v.(type) {
	case nil:
		return true
	case string:
		return vv == ""
	case bool:
		return !vv
	case float64:
		return vv == 0
	case int:
		return vv == 0
	}
	return false
}

// mapSliceToRows converts a slice of maps to rows with sorted headers.
func mapSliceToRows(items []map[string]interface{}) ([][]string, []string, error) {
	if len(items) == 0 {
		return nil, nil, nil
	}

	headerSet := make(map[string]bool)
	for _, item := range items {
		for k := range item {
			headerSet[k] = true
		}
	}

	headers := make([]string, 0, len(headerSet))
	for k := range headerSet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	rows := make([][]string, len(items))
	for i, item := range items {
		row := make([]string, len(headers))
		for j, h := range headers {
			v, ok := item[h]
			if !ok || v == nil {
				row[j] = ""
			} else {
				row[j] = fmt.Sprintf("%v", v)
			}
		}
		rows[i] = row
	}

	return rows, headers, nil
}

// padRight pads a string with spaces to the given width.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
