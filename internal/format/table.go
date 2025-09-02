// Package format provides utilities for formatting CLI output,
// including table rendering and data presentation.
package format

import (
	"fmt"
	"strings"
)

// TableWriter helps format tabular data
type TableWriter struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table with the specified headers
func NewTable(headers ...string) *TableWriter {
	return &TableWriter{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow adds a row of data to the table
func (t *TableWriter) AddRow(values ...string) {
	// Pad with empty strings if not enough values provided
	row := make([]string, len(t.headers))
	for i, value := range values {
		if i < len(row) {
			row[i] = value
		}
	}
	t.rows = append(t.rows, row)
}

// Render returns the formatted table as a string
func (t *TableWriter) Render() string {
	if len(t.headers) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(t.headers))

	// Start with header widths
	for i, header := range t.headers {
		colWidths[i] = len(header)
	}

	// Check row widths
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var result strings.Builder

	// Write headers
	for i, header := range t.headers {
		if i > 0 {
			result.WriteString("  ")
		}
		result.WriteString(fmt.Sprintf("%-*s", colWidths[i], header))
	}
	result.WriteString("\n")

	// Write separator
	for i := range t.headers {
		if i > 0 {
			result.WriteString("  ")
		}
		result.WriteString(strings.Repeat("-", colWidths[i]))
	}
	result.WriteString("\n")

	// Write rows
	for _, row := range t.rows {
		for i, cell := range row {
			if i > 0 {
				result.WriteString("  ")
			}
			if i < len(colWidths) {
				result.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}
