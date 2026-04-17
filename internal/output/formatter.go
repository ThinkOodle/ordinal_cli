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

// FormatOutput formats data according to the given format.
func FormatOutput(data interface{}, format Format) (string, error) {
	switch format {
	case FormatTable:
		return formatTable(data)
	case FormatCSV:
		return formatCSV(data)
	case FormatJSON:
		return formatJSON(data)
	default:
		return formatJSON(data)
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

// formatTable formats data as an aligned text table.
func formatTable(data interface{}) (string, error) {
	rows, headers, err := extractRows(data)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 && len(headers) == 0 {
		return "No results", nil
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

	return strings.TrimRight(buf.String(), "\n"), nil
}

// formatCSV formats data as CSV.
func formatCSV(data interface{}) (string, error) {
	rows, headers, err := extractRows(data)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 && len(headers) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return "", fmt.Errorf("writing csv header: %w", err)
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("writing csv row: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("flushing csv: %w", err)
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}

// extractRows converts data into string rows and headers for table/CSV output.
// Supports slices of maps, slices of structs, single maps, and single structs.
//
// For list-style wrapper responses with exactly one array-valued field
// alongside scalar pagination siblings (e.g. {"posts":[...], "nextCursor":"...",
// "hasMore":true}), the array is unwrapped and its elements become the rows.
// This keeps table/csv output useful for the many list endpoints that return
// an envelope rather than a bare array.
func extractRows(data interface{}) ([][]string, []string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling data: %w", err)
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(b, &items); err == nil {
		if len(items) == 0 {
			return nil, nil, nil
		}
		return mapSliceToRows(items)
	}

	var item map[string]interface{}
	if err := json.Unmarshal(b, &item); err == nil && len(item) > 0 {
		if unwrapped, ok := unwrapListEnvelope(item); ok {
			if len(unwrapped) == 0 {
				return nil, nil, nil
			}
			return mapSliceToRows(unwrapped)
		}
		return mapSliceToRows([]map[string]interface{}{item})
	}

	return nil, nil, fmt.Errorf("unsupported data type for table/csv output: %s", reflect.TypeOf(data))
}

// unwrapListEnvelope detects a list-response envelope — a JSON object with
// exactly one field whose value is an array of objects, alongside only scalar
// siblings (typical pagination metadata like nextCursor, hasMore, total). When
// matched, returns the array contents so table/csv output can render the items
// as rows instead of turning the envelope itself into a single row.
func unwrapListEnvelope(item map[string]interface{}) ([]map[string]interface{}, bool) {
	var (
		arrayKey   string
		arrayValue []interface{}
		arrayCount int
	)
	for k, v := range item {
		switch vv := v.(type) {
		case []interface{}:
			arrayKey = k
			arrayValue = vv
			arrayCount++
		case map[string]interface{}:
			// Nested objects aren't pagination metadata; don't treat this
			// as a list envelope.
			return nil, false
		}
	}
	if arrayCount != 1 || arrayKey == "" {
		return nil, false
	}

	items := make([]map[string]interface{}, 0, len(arrayValue))
	for _, entry := range arrayValue {
		m, ok := entry.(map[string]interface{})
		if !ok {
			return nil, false
		}
		items = append(items, m)
	}
	return items, true
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
