# Claude Code Instructions

## Project Overview

Festivus is a terminal-based text editor inspired by DOS EDIT, built with Go and the Bubbletea TUI framework. "A text editor for the rest of us."

## Tech Stack

- Go 1.21+
- Bubbletea - TUI framework (MVU architecture)
- Lipgloss - Styling
- go-runewidth - Unicode width calculation
- Chroma - Syntax highlighting

## Build & Run

```bash
go build
./festivus [filename]
```

## Project Structure

- `editor/` - Core editor logic (buffer, cursor, selection, undo, dialogs, file browser)
- `ui/` - UI components (menubar, statusbar, viewport, styles)
- `clipboard/` - Clipboard handling (native xclip/xsel/wl-clipboard, OSC52 for SSH)
- `syntax/` - Syntax highlighting (Chroma-based)
- `config/` - Configuration file handling

## Code Patterns

- Use direct ANSI escape codes for menu bar, status bar, find/replace bar backgrounds (lipgloss nesting causes color issues)
- Gap buffer for text storage
- Visual line counting for word wrap positioning
- Clipboard uses native tools (xclip/xsel/wl-clipboard) with OSC52 fallback for SSH

## Git Workflow

- Do not commit and push without explicit user request
