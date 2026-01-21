package ui

import (
	"strings"
	"unicode/utf8"
)

// MinimapRenderer renders a braille-based minimap of the document.
// Standard width is 8 (1 viewport indicator + 6 braille chars + 1 space).
//
// === MINIMAP SPECIFICATION (TODO: implement) ===
//
// Vertical mapping:
//   - 1 braille dot row = 1 visual line (respects word wrap)
//   - Each braille character = 4 visual lines (braille has 4 dot rows)
//   - Minimap height = ceil(total visual lines / 4)
//   - Minimap may be shorter or taller than viewport - not scaled to fit
//
// Horizontal mapping:
//   - 1 braille dot column = 5 source characters
//   - Each braille character = 10 source characters (2 dot columns × 5 chars)
//   - 6 braille characters = 60 source characters max
//   - Lines longer than 60 chars are truncated (not scaled)
//
// Fill logic:
//   - A dot is ON if there are >= 3 non-whitespace characters in that
//     5-character span (i.e., less than 2 char widths of whitespace)
//
// Viewport indicator:
//   - Option A: Current vertical bar │ on left side
//   - Option B: Reverse video on braille chars within viewport range
//
// Mouse interaction:
//   - Clicking on minimap navigates viewport to that location
type MinimapRenderer struct {
	styles  Styles
	enabled bool
}

// NewMinimapRenderer creates a new minimap renderer.
func NewMinimapRenderer(styles Styles) *MinimapRenderer {
	return &MinimapRenderer{
		styles:  styles,
		enabled: false, // Disabled by default
	}
}

// SetStyles updates the styles for runtime theme changes.
func (r *MinimapRenderer) SetStyles(styles Styles) {
	r.styles = styles
}

// SetEnabled enables or disables the minimap.
func (r *MinimapRenderer) SetEnabled(enabled bool) {
	r.enabled = enabled
}

// IsEnabled returns whether the minimap is enabled.
func (r *MinimapRenderer) IsEnabled() bool {
	return r.enabled
}

// Toggle toggles the minimap on/off.
func (r *MinimapRenderer) Toggle() bool {
	r.enabled = !r.enabled
	return r.enabled
}

// Render implements ColumnRenderer.
// Returns braille representation of the document with viewport indicator.
func (r *MinimapRenderer) Render(width, height int, state *RenderState) []string {
	if !r.enabled || width <= 0 || height <= 0 || state == nil {
		rows := make([]string, height)
		for i := range rows {
			rows[i] = strings.Repeat(" ", width)
		}
		return rows
	}

	// Layout: [indicator][braille chars][space]
	// indicator: 1 char showing if this row is in visible viewport
	// braille: width-2 chars of document content
	// space: 1 char padding on right
	brailleWidth := width - 2
	if brailleWidth < 1 {
		brailleWidth = 1
	}

	rows := make([]string, height)

	// Calculate how many document lines each minimap row represents
	totalLines := state.TotalLines
	if totalLines == 0 {
		totalLines = 1
	}

	// Scale factor: how many document lines per minimap row
	linesPerRow := float64(totalLines) / float64(height)
	if linesPerRow < 1 {
		linesPerRow = 1
	}

	// Calculate viewport indicator range
	visibleStart := state.ScrollY
	visibleEnd := state.ScrollY + height
	if state.WordWrap && state.TotalVisualLines > 0 {
		// With word wrap, use visual lines
		totalLines = state.TotalVisualLines
		linesPerRow = float64(totalLines) / float64(height)
		if linesPerRow < 1 {
			linesPerRow = 1
		}
	}

	// Calculate max line length for consistent scaling across all rows
	maxLineLen := 80 // Default minimum width
	for _, line := range state.Lines {
		lineLen := utf8.RuneCountInString(line)
		if lineLen > maxLineLen {
			maxLineLen = lineLen
		}
	}

	// Get theme colors
	ui := r.styles.Theme.UI
	indicatorColor := ColorToANSIFg(ui.MinimapIndicator)
	textColor := ColorToANSIFg(ui.MinimapText)
	resetCode := "\033[0m"

	for row := 0; row < height; row++ {
		var sb strings.Builder

		// Calculate which document lines this minimap row represents
		docLineStart := int(float64(row) * linesPerRow)
		docLineEnd := int(float64(row+1) * linesPerRow)
		if docLineEnd > totalLines {
			docLineEnd = totalLines
		}

		// Viewport indicator
		inViewport := docLineStart < visibleEnd && docLineEnd > visibleStart
		if inViewport {
			sb.WriteString(indicatorColor)
			sb.WriteString("│")
			sb.WriteString(resetCode)
		} else {
			sb.WriteString(" ")
		}

		// Braille representation of document content
		sb.WriteString(textColor)
		braille := r.renderBrailleRow(state.Lines, docLineStart, docLineEnd, brailleWidth, maxLineLen)
		sb.WriteString(braille)
		sb.WriteString(resetCode)

		// Right padding
		sb.WriteString(" ")

		rows[row] = sb.String()
	}

	return rows
}

// renderBrailleRow converts document lines to a braille string.
// Each braille character represents a 2-column x 4-row grid.
// The mapping is proportional: the entire line width maps to the minimap width.
// maxLineLen is the maximum line length in the document for consistent scaling.
func (r *MinimapRenderer) renderBrailleRow(lines []string, startLine, endLine, width, maxLineLen int) string {
	if len(lines) == 0 || startLine >= len(lines) {
		return strings.Repeat(" ", width)
	}

	// Limit to actual document
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine < 0 {
		startLine = 0
	}

	// Sample up to 4 lines for this braille row (braille char height)
	sampleRunes := make([][]rune, 0, 4)
	step := (endLine - startLine) / 4
	if step < 1 {
		step = 1
	}
	for i := startLine; i < endLine && len(sampleRunes) < 4; i += step {
		if i < len(lines) {
			sampleRunes = append(sampleRunes, []rune(lines[i]))
		}
	}

	// Pad to 4 lines
	for len(sampleRunes) < 4 {
		sampleRunes = append(sampleRunes, []rune{})
	}

	// Calculate how many source columns each braille character represents
	// Each braille char has 2 columns, so charsPerBraille is for the whole char
	charsPerBraille := float64(maxLineLen) / float64(width)
	if charsPerBraille < 1 {
		charsPerBraille = 1
	}

	// Generate braille characters
	var result strings.Builder

	for col := 0; col < width; col++ {
		srcColStart := int(float64(col) * charsPerBraille)
		srcColEnd := int(float64(col+1) * charsPerBraille)
		srcColMid := (srcColStart + srcColEnd) / 2

		// Build braille pattern from the 4x2 grid
		// Braille dots are numbered:
		// 1 4
		// 2 5
		// 3 6
		// 7 8
		var pattern rune = 0x2800 // Empty braille

		for rowOffset := 0; rowOffset < 4; rowOffset++ {
			lineRunes := sampleRunes[rowOffset]

			// Left column (dots 1,2,3,7)
			if hasContentAt(lineRunes, srcColStart, srcColMid) {
				switch rowOffset {
				case 0:
					pattern |= 0x01 // dot 1
				case 1:
					pattern |= 0x02 // dot 2
				case 2:
					pattern |= 0x04 // dot 3
				case 3:
					pattern |= 0x40 // dot 7
				}
			}

			// Right column (dots 4,5,6,8)
			if hasContentAt(lineRunes, srcColMid, srcColEnd) {
				switch rowOffset {
				case 0:
					pattern |= 0x08 // dot 4
				case 1:
					pattern |= 0x10 // dot 5
				case 2:
					pattern |= 0x20 // dot 6
				case 3:
					pattern |= 0x80 // dot 8
				}
			}
		}

		result.WriteRune(pattern)
	}

	return result.String()
}

// hasContentAt checks if a line has non-whitespace content in the given column range.
func hasContentAt(lineRunes []rune, start, end int) bool {
	if start < 0 {
		start = 0
	}
	if end > len(lineRunes) {
		end = len(lineRunes)
	}
	for i := start; i < end; i++ {
		if i < len(lineRunes) {
			r := lineRunes[i]
			if r != ' ' && r != '\t' {
				return true
			}
		}
	}
	return false
}

// MinimapWidth returns the standard width for the minimap column.
func MinimapWidth() int {
	return 8 // 1 indicator + 6 braille + 1 space
}

// CalculateMinimapMetrics returns metrics for mouse interaction with minimap.
type MinimapMetrics struct {
	LinesPerRow float64 // Document lines per minimap row
	TotalLines  int
	Height      int
}

// GetMetrics calculates minimap metrics for a given state.
func (r *MinimapRenderer) GetMetrics(height int, state *RenderState) MinimapMetrics {
	totalLines := state.TotalLines
	if state.WordWrap && state.TotalVisualLines > 0 {
		totalLines = state.TotalVisualLines
	}
	if totalLines == 0 {
		totalLines = 1
	}

	linesPerRow := float64(totalLines) / float64(height)
	if linesPerRow < 1 {
		linesPerRow = 1
	}

	return MinimapMetrics{
		LinesPerRow: linesPerRow,
		TotalLines:  totalLines,
		Height:      height,
	}
}

// RowToLine converts a minimap row click to a document line.
func (r *MinimapRenderer) RowToLine(row int, metrics MinimapMetrics) int {
	line := int(float64(row) * metrics.LinesPerRow)
	if line < 0 {
		return 0
	}
	if line >= metrics.TotalLines {
		return metrics.TotalLines - 1
	}
	return line
}

// Helper to get line length in runes
func lineRuneCount(line string) int {
	return utf8.RuneCountInString(line)
}
