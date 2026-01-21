package ui

import (
	"github.com/cornish/textivus-editor/syntax"
)

// ColumnRenderer is the interface that column renderers must implement.
// Each renderer produces exactly `width` visual characters per row.
type ColumnRenderer interface {
	// Render returns exactly `height` rows, each with exactly `width` visual characters.
	// ANSI codes don't count toward width - only visible characters do.
	Render(width, height int, state *RenderState) []string
}

// Column represents a single column in the compositor layout.
type Column struct {
	Width    int            // Fixed width in cells (0 if flexible)
	Flexible bool           // If true, this column takes remaining space
	Enabled  bool           // Whether this column is currently shown
	Renderer ColumnRenderer // The renderer for this column
}

// RenderState holds shared state passed to all column renderers.
// This allows columns to render consistently without direct coupling.
type RenderState struct {
	// Document content
	Lines []string // All lines in the document

	// Cursor position
	CursorLine int
	CursorCol  int

	// Scroll position
	ScrollY int // First visible line (visual line for word wrap)
	ScrollX int // Horizontal scroll offset

	// Selection state (map of line index to selection range)
	Selection map[int]SelectionRange

	// Syntax highlighting (map of line index to color spans)
	LineColors map[int][]syntax.ColorSpan

	// Display options
	WordWrap bool
	TabWidth int // Display width of tabs

	// Total document metrics (used by scrollbar, minimap)
	TotalLines       int // Total buffer lines
	TotalVisualLines int // Total visual lines (with word wrap)

	// Styles for rendering
	Styles Styles
}

// Note: SelectionRange is defined in viewport.go
