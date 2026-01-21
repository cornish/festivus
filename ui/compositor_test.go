package ui

import (
	"strings"
	"testing"
)

// mockRenderer is a simple test renderer that produces fixed content.
type mockRenderer struct {
	char string // Character to fill with
}

func (m *mockRenderer) Render(width, height int, state *RenderState) []string {
	rows := make([]string, height)
	for i := 0; i < height; i++ {
		rows[i] = strings.Repeat(m.char, width)
	}
	return rows
}

// mockColorRenderer produces content with ANSI color codes.
type mockColorRenderer struct {
	char  string
	color string // ANSI color code
}

func (m *mockColorRenderer) Render(width, height int, state *RenderState) []string {
	rows := make([]string, height)
	for i := 0; i < height; i++ {
		rows[i] = m.color + strings.Repeat(m.char, width) + "\033[0m"
	}
	return rows
}

func TestCompositorBasic(t *testing.T) {
	c := NewCompositor(20, 3)

	c.AddColumn(Column{
		Width:    5,
		Flexible: false,
		Enabled:  true,
		Renderer: &mockRenderer{char: "L"},
	})
	c.AddColumn(Column{
		Width:    0,
		Flexible: true,
		Enabled:  true,
		Renderer: &mockRenderer{char: "T"},
	})
	c.AddColumn(Column{
		Width:    1,
		Flexible: false,
		Enabled:  true,
		Renderer: &mockRenderer{char: "S"},
	})

	result := c.Render(nil)
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Line numbers (5) + Text (14) + Scrollbar (1) = 20
	expected := "LLLLL" + strings.Repeat("T", 14) + "S"
	for i, line := range lines {
		if line != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, line)
		}
	}
}

func TestCompositorDisabledColumn(t *testing.T) {
	c := NewCompositor(20, 2)

	c.AddColumn(Column{
		Width:    5,
		Flexible: false,
		Enabled:  false, // Disabled
		Renderer: &mockRenderer{char: "L"},
	})
	c.AddColumn(Column{
		Width:    0,
		Flexible: true,
		Enabled:  true,
		Renderer: &mockRenderer{char: "T"},
	})

	result := c.Render(nil)
	lines := strings.Split(result, "\n")

	// With line numbers disabled, text should fill all 20 chars
	expected := strings.Repeat("T", 20)
	for i, line := range lines {
		if line != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, line)
		}
	}
}

func TestCompositorFlexibleWidth(t *testing.T) {
	c := NewCompositor(80, 1)

	c.AddColumn(Column{Width: 5, Enabled: true, Renderer: &mockRenderer{char: "1"}})
	c.AddColumn(Column{Flexible: true, Enabled: true, Renderer: &mockRenderer{char: "2"}})
	c.AddColumn(Column{Width: 8, Enabled: true, Renderer: &mockRenderer{char: "3"}})
	c.AddColumn(Column{Width: 1, Enabled: true, Renderer: &mockRenderer{char: "4"}})

	flexWidth := c.FlexibleColumnWidth()
	expected := 80 - 5 - 8 - 1 // 66
	if flexWidth != expected {
		t.Errorf("FlexibleColumnWidth: expected %d, got %d", expected, flexWidth)
	}
}

func TestCompositorWithANSI(t *testing.T) {
	c := NewCompositor(10, 2)

	c.AddColumn(Column{
		Width:    5,
		Enabled:  true,
		Renderer: &mockColorRenderer{char: "X", color: "\033[31m"}, // Red
	})
	c.AddColumn(Column{
		Width:    5,
		Enabled:  true,
		Renderer: &mockRenderer{char: "Y"},
	})

	result := c.Render(nil)
	lines := strings.Split(result, "\n")

	// Visual width should be 10, but with ANSI codes the string is longer
	for _, line := range lines {
		vw := visualWidth(line)
		if vw != 10 {
			t.Errorf("Expected visual width 10, got %d for line %q", vw, line)
		}
	}
}

func TestCalculateColumnWidths(t *testing.T) {
	c := NewCompositor(100, 10)

	c.SetColumns([]Column{
		{Width: 5, Enabled: true},
		{Width: 0, Flexible: true, Enabled: true},
		{Width: 8, Enabled: true},
		{Width: 1, Enabled: true},
	})

	widths := c.calculateColumnWidths()

	if widths[0] != 5 {
		t.Errorf("Column 0: expected 5, got %d", widths[0])
	}
	if widths[1] != 86 { // 100 - 5 - 8 - 1
		t.Errorf("Column 1 (flexible): expected 86, got %d", widths[1])
	}
	if widths[2] != 8 {
		t.Errorf("Column 2: expected 8, got %d", widths[2])
	}
	if widths[3] != 1 {
		t.Errorf("Column 3: expected 1, got %d", widths[3])
	}
}

func TestPadToWidth(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected int // expected visual width
	}{
		{"hello", 10, 10},
		{"hello", 5, 5},
		{"hello", 3, 3},
		{"\033[31mhi\033[0m", 5, 5},    // With ANSI codes
		{"\033[31mhello\033[0m", 3, 3}, // Truncate with ANSI
	}

	for _, tc := range tests {
		result := padToWidth(tc.input, tc.width)
		vw := visualWidth(result)
		if vw != tc.expected {
			t.Errorf("padToWidth(%q, %d): expected visual width %d, got %d (result: %q)",
				tc.input, tc.width, tc.expected, vw, result)
		}
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"\033[31mhello\033[0m", "hello"},
		{"\033[1;32mgreen\033[0m text", "green text"},
		{"no codes here", "no codes here"},
	}

	for _, tc := range tests {
		result := stripANSI(tc.input)
		if result != tc.expected {
			t.Errorf("stripANSI(%q): expected %q, got %q", tc.input, tc.expected, result)
		}
	}
}

func TestVisualWidth(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"\033[31mhello\033[0m", 5},
		{"日本語", 6}, // 3 wide chars
		{"\033[31m日本\033[0m", 4},
	}

	for _, tc := range tests {
		result := visualWidth(tc.input)
		if result != tc.expected {
			t.Errorf("visualWidth(%q): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}
