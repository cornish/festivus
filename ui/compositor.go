package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// Compositor joins multiple columns horizontally to produce the final viewport output.
type Compositor struct {
	columns []Column
	width   int
	height  int
}

// NewCompositor creates a new compositor with the given dimensions.
func NewCompositor(width, height int) *Compositor {
	return &Compositor{
		columns: make([]Column, 0),
		width:   width,
		height:  height,
	}
}

// SetSize updates the compositor dimensions.
func (c *Compositor) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Width returns the compositor width.
func (c *Compositor) Width() int {
	return c.width
}

// Height returns the compositor height.
func (c *Compositor) Height() int {
	return c.height
}

// AddColumn adds a column to the compositor.
func (c *Compositor) AddColumn(col Column) {
	c.columns = append(c.columns, col)
}

// SetColumns replaces all columns.
func (c *Compositor) SetColumns(cols []Column) {
	c.columns = cols
}

// GetColumns returns a copy of the current columns.
func (c *Compositor) GetColumns() []Column {
	result := make([]Column, len(c.columns))
	copy(result, c.columns)
	return result
}

// EnableColumn enables or disables a column by index.
func (c *Compositor) EnableColumn(index int, enabled bool) {
	if index >= 0 && index < len(c.columns) {
		c.columns[index].Enabled = enabled
	}
}

// calculateColumnWidths determines the actual width for each enabled column.
// Fixed columns get their specified width; the flexible column gets the remainder.
func (c *Compositor) calculateColumnWidths() []int {
	widths := make([]int, len(c.columns))
	flexibleIdx := -1
	usedWidth := 0

	// First pass: assign fixed widths and find flexible column
	for i, col := range c.columns {
		if !col.Enabled {
			widths[i] = 0
			continue
		}
		if col.Flexible {
			flexibleIdx = i
		} else {
			widths[i] = col.Width
			usedWidth += col.Width
		}
	}

	// Second pass: assign remaining width to flexible column
	if flexibleIdx >= 0 {
		remaining := c.width - usedWidth
		if remaining < 1 {
			remaining = 1 // Minimum 1 character
		}
		widths[flexibleIdx] = remaining
	}

	return widths
}

// FlexibleColumnWidth returns the calculated width of the flexible column.
// This is useful for external code that needs to know the text area width.
func (c *Compositor) FlexibleColumnWidth() int {
	widths := c.calculateColumnWidths()
	for i, col := range c.columns {
		if col.Enabled && col.Flexible {
			return widths[i]
		}
	}
	return c.width // No flexible column, return full width
}

// Render renders all enabled columns and joins them horizontally.
func (c *Compositor) Render(state *RenderState) string {
	if len(c.columns) == 0 || c.height <= 0 {
		return ""
	}

	widths := c.calculateColumnWidths()

	// Render each enabled column
	columnOutputs := make([][]string, len(c.columns))
	for i, col := range c.columns {
		if !col.Enabled || widths[i] == 0 || col.Renderer == nil {
			// Disabled or zero-width: produce empty rows
			columnOutputs[i] = make([]string, c.height)
			for j := range columnOutputs[i] {
				columnOutputs[i][j] = ""
			}
			continue
		}
		columnOutputs[i] = col.Renderer.Render(widths[i], c.height, state)
		// Ensure we have exactly c.height rows
		if len(columnOutputs[i]) < c.height {
			// Pad with empty rows
			for len(columnOutputs[i]) < c.height {
				columnOutputs[i] = append(columnOutputs[i], strings.Repeat(" ", widths[i]))
			}
		} else if len(columnOutputs[i]) > c.height {
			columnOutputs[i] = columnOutputs[i][:c.height]
		}
	}

	// Join columns horizontally, row by row
	var result strings.Builder
	for row := 0; row < c.height; row++ {
		if row > 0 {
			result.WriteString("\n")
		}
		for i, col := range c.columns {
			if !col.Enabled || widths[i] == 0 {
				continue
			}
			result.WriteString(columnOutputs[i][row])
		}
	}

	return result.String()
}

// visualWidth calculates the visible width of a string, ignoring ANSI escape codes.
func visualWidth(s string) int {
	return runewidth.StringWidth(stripANSI(s))
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// padToWidth pads a string (which may contain ANSI codes) to exactly the target visual width.
// If the string is too long, it's truncated.
func padToWidth(s string, width int) string {
	vw := visualWidth(s)
	if vw == width {
		return s
	}
	if vw > width {
		// Need to truncate - this is tricky with ANSI codes
		return truncateToWidth(s, width)
	}
	// Need to pad
	return s + strings.Repeat(" ", width-vw)
}

// truncateToWidth truncates a string with ANSI codes to a visual width.
func truncateToWidth(s string, width int) string {
	var result strings.Builder
	inEscape := false
	visualPos := 0

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			result.WriteRune(r)
			continue
		}
		if inEscape {
			result.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}

		rw := runewidth.RuneWidth(r)
		if visualPos+rw > width {
			break
		}
		result.WriteRune(r)
		visualPos += rw
	}

	// Pad if we ended up short (e.g., wide character at boundary)
	if visualPos < width {
		result.WriteString(strings.Repeat(" ", width-visualPos))
	}

	return result.String()
}
