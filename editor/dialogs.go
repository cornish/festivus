package editor

import (
	"strings"
)

// overlayLineAt overlays the dropdown line on top of the viewport line at the given offset,
// preserving viewport content on both sides of the dropdown
func overlayLineAt(dropLine, viewportLine string, offset int) string {
	// Calculate the visual width of the dropdown line (strip ANSI codes)
	dropWidth := visualWidth(dropLine)

	// Get the viewport content as runes (stripped of ANSI for positioning)
	vpRunes := []rune(stripAnsi(viewportLine))

	// Build the result: prefix + dropdown + suffix
	var result strings.Builder

	// Prefix: viewport content before the dropdown (or spaces if line is short)
	if offset > 0 {
		if len(vpRunes) >= offset {
			// Use viewport content as prefix
			result.WriteString(string(vpRunes[:offset]))
		} else {
			// Viewport line is shorter than offset - use what we have plus padding
			result.WriteString(string(vpRunes))
			result.WriteString(strings.Repeat(" ", offset-len(vpRunes)))
		}
	}

	// The dropdown itself
	result.WriteString(dropLine)

	// Suffix: viewport content after the dropdown
	suffixStart := offset + dropWidth
	if suffixStart < len(vpRunes) {
		result.WriteString(string(vpRunes[suffixStart:]))
	}

	return result.String()
}

// stripAnsi removes ANSI escape sequences from a string
func stripAnsi(s string) string {
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

// visualWidth calculates the visible width of a string (ignoring ANSI codes)
func visualWidth(s string) int {
	return len([]rune(stripAnsi(s)))
}

// overlayAboutDialog overlays the about dialog centered on the viewport
func (e *Editor) overlayAboutDialog(viewportContent string) string {
	// Use the stored quote (selected when dialog opened)
	quote := e.aboutQuote
	if quote == "" {
		quote = "A Festivus for the rest of us!"
	}

	// ASCII art from festivus.txt - art is 62 chars, box is 64 for padding
	boxWidth := 64
	centerText := func(s string) string {
		sLen := len(s)
		if sLen >= boxWidth {
			// Truncate if too long
			return s[:boxWidth]
		}
		padLeft := (boxWidth - sLen) / 2
		padRight := boxWidth - sLen - padLeft
		return strings.Repeat(" ", padLeft) + s + strings.Repeat(" ", padRight)
	}

	// Split quote into lines if too long (max 60 chars per line)
	maxLineWidth := 60
	var quoteLines []string
	quotedText := "\"" + quote + "\""
	if len(quotedText) <= maxLineWidth {
		// Fits on one line
		quoteLines = []string{centerText(quotedText)}
	} else {
		// Split at word boundary
		words := strings.Fields(quote)
		line1 := "\""
		line2 := ""
		for _, word := range words {
			testLine := line1
			if len(line1) > 1 {
				testLine += " "
			}
			testLine += word
			if line2 == "" && len(testLine) <= maxLineWidth {
				line1 = testLine
			} else {
				if line2 != "" {
					line2 += " "
				}
				line2 += word
			}
		}
		line2 += "\""
		quoteLines = []string{centerText(line1), centerText(line2)}
	}

	// Choose logo based on ASCII mode
	var logoLines []string
	if e.box.Lock == "*" {
		// ASCII mode - use asterisk art (64 chars wide to match boxWidth)
		logoLines = []string{
			"      *****  *****   ****  *****  ***  *   *  *   *   ****      ",
			"      *      *      *        *     *   *   *  *   *  *          ",
			"      ****   ****    ***     *     *   *   *  *   *   ***       ",
			"      *      *          *    *     *    * *   *   *      *      ",
			"      *      *****  ****     *    ***    *     ***   ****       ",
			"                                                                ",
		}
	} else {
		// Unicode mode - use block art
		logoLines = []string{
			" ███████╗███████╗███████╗████████╗██╗██╗   ██╗██╗   ██╗███████╗ ",
			" ██╔════╝██╔════╝██╔════╝╚══██╔══╝██║██║   ██║██║   ██║██╔════╝ ",
			" █████╗  █████╗  ███████╗   ██║   ██║██║   ██║██║   ██║███████╗ ",
			" ██╔══╝  ██╔══╝  ╚════██║   ██║   ██║╚██╗ ██╔╝██║   ██║╚════██║ ",
			" ██║     ███████╗███████║   ██║   ██║ ╚████╔╝ ╚██████╔╝███████║ ",
			" ╚═╝     ╚══════╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝   ╚═════╝ ╚══════╝ ",
		}
	}

	aboutLines := []string{strings.Repeat(" ", boxWidth)}
	aboutLines = append(aboutLines, logoLines...)
	aboutLines = append(aboutLines,
		strings.Repeat(" ", boxWidth),
		centerText("A Text Editor for the Rest of Us"),
		strings.Repeat(" ", boxWidth),
		centerText("Version 0.1.0"),
		centerText("github.com/cornish/festivus"),
		centerText("Copyright (c) 2025"),
		strings.Repeat(" ", boxWidth),
	)
	aboutLines = append(aboutLines, quoteLines...)
	aboutLines = append(aboutLines,
		strings.Repeat(" ", boxWidth),
		centerText("Press any key to continue..."),
		strings.Repeat(" ", boxWidth),
	)
	boxHeight := len(aboutLines)

	// Calculate centering
	startX := (e.width - boxWidth) / 2
	if startX < 0 {
		startX = 0
	}
	startY := (e.viewport.Height() - boxHeight) / 2
	if startY < 0 {
		startY = 0
	}

	viewportLines := strings.Split(viewportContent, "\n")

	for i, aboutLine := range aboutLines {
		viewportY := startY + i
		if viewportY >= 0 && viewportY < len(viewportLines) {
			// Build the styled about line with cyan background
			var styledLine strings.Builder
			styledLine.WriteString("\033[46;30m") // Cyan bg, black text
			styledLine.WriteString(aboutLine)
			styledLine.WriteString("\033[0m")

			// Overlay on viewport line
			viewportLines[viewportY] = overlayLineAt(styledLine.String(), viewportLines[viewportY], startX)
		}
	}

	return strings.Join(viewportLines, "\n")
}

// overlayHelpDialog overlays the help dialog centered on the viewport
func (e *Editor) overlayHelpDialog(viewportContent string) string {
	// Two-column layout for keyboard shortcuts
	boxWidth := 72
	innerWidth := boxWidth - 2 // 70
	colWidth := 33             // Each column width
	// Layout: colWidth (33) + separator "  │ " (4) + colWidth (33) = 70

	padText := func(s string, width int) string {
		if len(s) > width {
			return s[:width]
		}
		return s + strings.Repeat(" ", width-len(s))
	}

	centerText := func(s string, width int) string {
		if len(s) >= width {
			return s[:width]
		}
		padLeft := (width - len(s)) / 2
		padRight := width - len(s) - padLeft
		return strings.Repeat(" ", padLeft) + s + strings.Repeat(" ", padRight)
	}

	// Define shortcuts in two columns
	leftCol := []string{
		"  FILE",
		"  Ctrl+N       New file",
		"  Ctrl+O       Open file",
		"  Ctrl+W       Close file",
		"  Ctrl+S       Save file",
		"  Ctrl+Q       Quit",
		"",
		"  EDIT",
		"  Ctrl+Z       Undo",
		"  Ctrl+Y       Redo",
		"  Ctrl+X       Cut",
		"  Ctrl+C       Copy",
		"  Ctrl+V       Paste",
		"  Ctrl+A       Select all",
		"  Ctrl+F       Find",
	}

	rightCol := []string{
		"  NAVIGATION",
		"  Arrows       Move cursor",
		"  Ctrl+Left/Right  Move by word",
		"  Home/End     Start/end of line",
		"  Ctrl+Home    Start of file",
		"  Ctrl+End     End of file",
		"  PgUp/PgDn    Page up/down",
		"",
		"  SELECTION",
		"  Shift+Arrows    Select text",
		"  Ctrl+Shift+L/R  Select word",
		"  Shift+Home/End  Select to line",
		"",
		"  MOUSE: Click, Drag, Scroll",
	}

	// Build help lines
	var helpLines []string

	// Top border with title
	title := " Keyboard Shortcuts "
	titlePadLeft := (innerWidth - len(title)) / 2
	titlePadRight := innerWidth - len(title) - titlePadLeft
	helpLines = append(helpLines, e.box.TopLeft+strings.Repeat(e.box.Horizontal, titlePadLeft)+title+strings.Repeat(e.box.Horizontal, titlePadRight)+e.box.TopRight)

	// Empty line
	helpLines = append(helpLines, e.box.Vertical+strings.Repeat(" ", innerWidth)+e.box.Vertical)

	// Build two-column content
	maxRows := len(leftCol)
	if len(rightCol) > maxRows {
		maxRows = len(rightCol)
	}

	colSep := "  " + e.box.Vertical + " "
	for i := 0; i < maxRows; i++ {
		left := ""
		right := ""
		if i < len(leftCol) {
			left = leftCol[i]
		}
		if i < len(rightCol) {
			right = rightCol[i]
		}
		line := padText(left, colWidth) + colSep + padText(right, colWidth)
		helpLines = append(helpLines, e.box.Vertical+line+e.box.Vertical)
	}

	// Empty line
	helpLines = append(helpLines, e.box.Vertical+strings.Repeat(" ", innerWidth)+e.box.Vertical)

	// Options section
	helpLines = append(helpLines, e.box.Vertical+centerText("OPTIONS: Ctrl+L Line Numbers", innerWidth)+e.box.Vertical)
	helpLines = append(helpLines, e.box.Vertical+centerText("MENUS: F10 or Alt+F/E/O/H", innerWidth)+e.box.Vertical)

	// Empty line
	helpLines = append(helpLines, e.box.Vertical+strings.Repeat(" ", innerWidth)+e.box.Vertical)

	// Footer
	helpLines = append(helpLines, e.box.Vertical+centerText("Press any key to continue...", innerWidth)+e.box.Vertical)

	// Bottom border
	helpLines = append(helpLines, e.box.BottomLeft+strings.Repeat(e.box.Horizontal, innerWidth)+e.box.BottomRight)

	boxHeight := len(helpLines)

	// Calculate centering
	startX := (e.width - boxWidth) / 2
	if startX < 0 {
		startX = 0
	}
	startY := (e.viewport.Height() - boxHeight) / 2
	if startY < 0 {
		startY = 0
	}

	viewportLines := strings.Split(viewportContent, "\n")

	for i, helpLine := range helpLines {
		viewportY := startY + i
		if viewportY >= 0 && viewportY < len(viewportLines) {
			// Build the styled help line with cyan background
			var styledLine strings.Builder
			styledLine.WriteString("\033[46;30m") // Cyan bg, black text
			styledLine.WriteString(helpLine)
			styledLine.WriteString("\033[0m")

			// Overlay on viewport line
			viewportLines[viewportY] = overlayLineAt(styledLine.String(), viewportLines[viewportY], startX)
		}
	}

	return strings.Join(viewportLines, "\n")
}
