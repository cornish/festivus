package ui

import (
	"strings"
	"unicode/utf8"
)

// LineNumberRenderer renders line numbers in a column.
// Standard width is 5 (4 digits + 1 space separator).
type LineNumberRenderer struct {
	styles Styles
}

// NewLineNumberRenderer creates a new line number renderer.
func NewLineNumberRenderer(styles Styles) *LineNumberRenderer {
	return &LineNumberRenderer{styles: styles}
}

// SetStyles updates the styles for runtime theme changes.
func (r *LineNumberRenderer) SetStyles(styles Styles) {
	r.styles = styles
}

// Render implements ColumnRenderer.
// Returns line numbers for visible lines, with the cursor line highlighted.
func (r *LineNumberRenderer) Render(width, height int, state *RenderState) []string {
	if width <= 0 || height <= 0 {
		return make([]string, height)
	}

	rows := make([]string, height)
	numWidth := width - 1 // Reserve 1 char for separator space

	if state.WordWrap {
		r.renderWrapped(rows, width, numWidth, height, state)
	} else {
		r.renderNoWrap(rows, width, numWidth, height, state)
	}

	return rows
}

// renderNoWrap renders line numbers without word wrap.
func (r *LineNumberRenderer) renderNoWrap(rows []string, width, numWidth, height int, state *RenderState) {
	// Get colors from theme
	ui := r.styles.Theme.UI
	normalColor := ColorToANSIFg(ui.LineNumber)
	activeColor := ColorToANSIFg(ui.LineNumberActive)
	resetCode := "\033[0m"

	for row := 0; row < height; row++ {
		lineIdx := state.ScrollY + row

		var sb strings.Builder
		if lineIdx < len(state.Lines) {
			// Real line - show number
			lineNum := lineIdx + 1 // 1-indexed
			numStr := padLeftStr(itoaLocal(lineNum), numWidth)

			if lineIdx == state.CursorLine {
				sb.WriteString(activeColor)
			} else {
				sb.WriteString(normalColor)
			}
			sb.WriteString(numStr)
			sb.WriteString(resetCode)
			sb.WriteString(" ")
		} else {
			// Past end of file - empty gutter
			sb.WriteString(strings.Repeat(" ", width))
		}
		rows[row] = sb.String()
	}
}

// renderWrapped renders line numbers with word wrap.
// Only the first visual line of each buffer line shows the number.
func (r *LineNumberRenderer) renderWrapped(rows []string, width, numWidth, height int, state *RenderState) {
	// Get colors from theme
	ui := r.styles.Theme.UI
	normalColor := ColorToANSIFg(ui.LineNumber)
	activeColor := ColorToANSIFg(ui.LineNumberActive)
	resetCode := "\033[0m"

	// Calculate text width (we need this to determine wrap points)
	// This is a bit of a hack - we don't know the text column width here.
	// For now, estimate based on a typical layout.
	// TODO: Pass text width through RenderState
	textWidth := 80 // Default estimate

	// Find which buffer line corresponds to ScrollY visual line
	visualLine := 0
	bufferLine := 0
	wrapOffset := 0

	for bufferLine < len(state.Lines) && visualLine < state.ScrollY {
		lineLen := utf8.RuneCountInString(state.Lines[bufferLine])
		wrappedCount := countWrappedLinesForWidth(lineLen, textWidth)

		if visualLine+wrappedCount > state.ScrollY {
			// Start partway through this line
			wrapOffset = state.ScrollY - visualLine
			break
		}
		visualLine += wrappedCount
		bufferLine++
	}

	// Render visible rows
	for row := 0; row < height; row++ {
		var sb strings.Builder

		if bufferLine >= len(state.Lines) {
			// Past end of file
			sb.WriteString(strings.Repeat(" ", width))
			rows[row] = sb.String()
			continue
		}

		lineLen := utf8.RuneCountInString(state.Lines[bufferLine])
		wrappedCount := countWrappedLinesForWidth(lineLen, textWidth)

		if wrapOffset == 0 {
			// First visual line of buffer line - show number
			lineNum := bufferLine + 1
			numStr := padLeftStr(itoaLocal(lineNum), numWidth)

			if bufferLine == state.CursorLine {
				sb.WriteString(activeColor)
			} else {
				sb.WriteString(normalColor)
			}
			sb.WriteString(numStr)
			sb.WriteString(resetCode)
			sb.WriteString(" ")
		} else {
			// Continuation line - empty gutter
			sb.WriteString(strings.Repeat(" ", width))
		}

		rows[row] = sb.String()

		// Move to next visual line
		wrapOffset++
		if wrapOffset >= wrappedCount {
			wrapOffset = 0
			bufferLine++
		}
	}
}

// countWrappedLinesForWidth returns how many visual lines a buffer line takes.
func countWrappedLinesForWidth(lineLen, textWidth int) int {
	if textWidth <= 0 {
		return 1
	}
	if lineLen == 0 {
		return 1
	}
	return (lineLen + textWidth - 1) / textWidth
}

// padLeftStr pads a string with spaces on the left to reach the target width.
func padLeftStr(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

// itoaLocal converts an int to string (local copy to avoid import).
func itoaLocal(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
