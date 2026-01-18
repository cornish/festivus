# Claude Code Instructions

## Project Overview

Festivus is a terminal-based text editor inspired by DOS EDIT, built with Go and the Bubbletea TUI framework. "A text editor for the rest of us."

## Tech Stack

- Go 1.21+
- Bubbletea - TUI framework (MVU architecture)
- Lipgloss - Styling
- go-runewidth - Unicode width calculation

## Build & Run

```bash
go build
./festivus [filename]
```

## Project Structure

- `editor/` - Core editor logic (buffer, cursor, selection, undo)
- `ui/` - UI components (menubar, statusbar, viewport, styles)
- `clipboard/` - Clipboard handling (OSC52 for SSH, local fallback)

## Code Patterns

- Use direct ANSI escape codes for menu bar, status bar, find bar backgrounds (lipgloss nesting causes color issues)
- Gap buffer for text storage
- Visual line counting for word wrap positioning

## Git Workflow

- Do not commit and push without explicit user request
