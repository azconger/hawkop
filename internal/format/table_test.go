package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TableTestSuite struct {
	suite.Suite
}

func (suite *TableTestSuite) TestNewTable() {
	table := NewTable("ID", "NAME", "STATUS")

	assert.Len(suite.T(), table.headers, 3)
	assert.Equal(suite.T(), "ID", table.headers[0])
	assert.Equal(suite.T(), "NAME", table.headers[1])
	assert.Equal(suite.T(), "STATUS", table.headers[2])
	assert.Len(suite.T(), table.rows, 0)
}

func (suite *TableTestSuite) TestAddRow() {
	table := NewTable("ID", "NAME", "STATUS")

	// Add complete row
	table.AddRow("123", "Test App", "ACTIVE")
	assert.Len(suite.T(), table.rows, 1)
	assert.Equal(suite.T(), []string{"123", "Test App", "ACTIVE"}, table.rows[0])

	// Add incomplete row (should pad with empty strings)
	table.AddRow("456", "Another App")
	assert.Len(suite.T(), table.rows, 2)
	assert.Equal(suite.T(), []string{"456", "Another App", ""}, table.rows[1])

	// Add row with extra values (should truncate)
	table.AddRow("789", "Third App", "INACTIVE", "EXTRA", "MORE")
	assert.Len(suite.T(), table.rows, 3)
	assert.Equal(suite.T(), []string{"789", "Third App", "INACTIVE"}, table.rows[2])
}

func (suite *TableTestSuite) TestRender_EmptyTable() {
	table := NewTable()
	result := table.Render()
	assert.Equal(suite.T(), "", result)
}

func (suite *TableTestSuite) TestRender_HeadersOnly() {
	table := NewTable("ID", "NAME", "STATUS")
	result := table.Render()

	expected := "ID  NAME  STATUS\n--  ----  ------\n"
	assert.Equal(suite.T(), expected, result)
}

func (suite *TableTestSuite) TestRender_WithData() {
	table := NewTable("ID", "NAME", "STATUS")
	table.AddRow("123", "Test App", "ACTIVE")
	table.AddRow("456789", "Another App With Long Name", "INACTIVE")

	result := table.Render()

	// Verify structure
	lines := strings.Split(result, "\n")
	assert.Len(suite.T(), lines, 5) // header + separator + 2 rows + final newline creates empty line

	// Check header line
	headerLine := lines[0]
	assert.Contains(suite.T(), headerLine, "ID")
	assert.Contains(suite.T(), headerLine, "NAME")
	assert.Contains(suite.T(), headerLine, "STATUS")

	// Check separator line
	separatorLine := lines[1]
	assert.Contains(suite.T(), separatorLine, "---")

	// Check data rows
	firstRow := lines[2]
	assert.Contains(suite.T(), firstRow, "123")
	assert.Contains(suite.T(), firstRow, "Test App")
	assert.Contains(suite.T(), firstRow, "ACTIVE")

	secondRow := lines[3]
	assert.Contains(suite.T(), secondRow, "456789")
	assert.Contains(suite.T(), secondRow, "Another App With Long Name")
	assert.Contains(suite.T(), secondRow, "INACTIVE")
}

func (suite *TableTestSuite) TestRender_ColumnWidthCalculation() {
	table := NewTable("ID", "APPLICATION")
	table.AddRow("1", "Short")
	table.AddRow("123456", "Very Long Application Name That Exceeds Header Width")

	result := table.Render()

	// The APPLICATION column should be sized to fit the longest content
	lines := strings.Split(result, "\n")
	headerLine := lines[0]

	// Find the APPLICATION column position and verify it's wide enough
	appHeaderStart := strings.Index(headerLine, "APPLICATION")
	assert.NotEqual(suite.T(), -1, appHeaderStart)

	// Check that the long application name fits
	dataLine := lines[3] // Second data row
	assert.Contains(suite.T(), dataLine, "Very Long Application Name That Exceeds Header Width")
}

func (suite *TableTestSuite) TestRender_EmptyValues() {
	table := NewTable("ID", "NAME", "STATUS")
	table.AddRow("123", "", "ACTIVE")
	table.AddRow("", "Test App", "")

	result := table.Render()

	// Should handle empty values gracefully
	lines := strings.Split(result, "\n")
	assert.Len(suite.T(), lines, 5) // header + separator + 2 rows + empty

	// Verify empty values are handled
	firstRow := lines[2]
	assert.Contains(suite.T(), firstRow, "123")
	assert.Contains(suite.T(), firstRow, "ACTIVE")

	secondRow := lines[3]
	assert.Contains(suite.T(), secondRow, "Test App")
}

func (suite *TableTestSuite) TestRender_ColumnAlignment() {
	table := NewTable("SHORT", "VERY_LONG_HEADER")
	table.AddRow("A", "B")

	result := table.Render()
	lines := strings.Split(result, "\n")

	headerLine := lines[0]
	separatorLine := lines[1]
	dataLine := lines[2]

	// Verify alignment - SHORT should be left-padded to match separator
	shortHeaderPos := strings.Index(headerLine, "SHORT")
	shortSepStart := strings.Index(separatorLine, "-----")
	assert.Equal(suite.T(), shortHeaderPos, shortSepStart)

	// Data should align with headers
	dataAPos := strings.Index(dataLine, "A")
	assert.Equal(suite.T(), shortHeaderPos, dataAPos)
}

func TestTableTestSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}
