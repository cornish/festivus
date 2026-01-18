# Festivus

**A Text Editor for the Rest of Us**

Festivus is a terminal-based text editor inspired by the classic DOS EDIT, built with Go and the [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework.

![Festivus Screenshot](screenshot.png)

## Features

- **Instant startup** - No bloat, just editing
- **Classic DOS EDIT styling** - Dark blue menu bar and status bar with cyan highlights
- **Modern keyboard shortcuts** - Ctrl+S, Ctrl+C, Ctrl+V, Ctrl+Z, etc.
- **Mouse support** - Click to position cursor, drag to select, scroll wheel
- **Shift+Arrow selection** - Select text the modern way
- **Word wrap** - Toggle via Options menu
- **Line numbers** - Toggle via Options menu or Ctrl+L
- **Syntax highlighting** - Auto-detected by file extension
- **Find & Replace** - Ctrl+F to find, Ctrl+H to find and replace
- **Go to Line** - Ctrl+G to jump to a specific line
- **Cut Line** - Ctrl+K to cut the entire current line (like nano)
- **Word & Character counts** - Displayed in the status bar
- **Clipboard support** - Native X11/Wayland support, OSC52 for SSH
- **Undo/Redo** - Ctrl+Z / Ctrl+Y with full history

## Installation

### From Source

Requires Go 1.21 or later.

```bash
git clone https://github.com/cornish/festivus.git
cd festivus
go build
./festivus [filename]
```

### Clipboard Support (Linux)

For clipboard integration with other applications, install one of:

```bash
# X11
sudo apt install xclip
# or
sudo apt install xsel

# Wayland
sudo apt install wl-clipboard
```

Without these tools, copy/paste will only work within Festivus.

## Keyboard Shortcuts

### File Operations
| Action | Shortcut |
|--------|----------|
| New | Ctrl+N |
| Open | Ctrl+O |
| Save | Ctrl+S |
| Close | Ctrl+W |
| Quit | Ctrl+Q |

### Editing
| Action | Shortcut |
|--------|----------|
| Undo | Ctrl+Z |
| Redo | Ctrl+Y |
| Cut | Ctrl+X |
| Copy | Ctrl+C |
| Paste | Ctrl+V |
| Cut Line | Ctrl+K |
| Select All | Ctrl+A |

### Search
| Action | Shortcut |
|--------|----------|
| Find | Ctrl+F |
| Find Next | F3 |
| Replace | Ctrl+H |
| Go to Line | Ctrl+G |

### Navigation
| Action | Shortcut |
|--------|----------|
| Start of file | Ctrl+Home |
| End of file | Ctrl+End |
| Start of line | Home |
| End of line | End |
| Word left/right | Ctrl+Left/Right |
| Page up/down | PgUp/PgDn |

### Selection
| Action | Shortcut |
|--------|----------|
| Select with cursor | Shift+Arrow |
| Select word | Ctrl+Shift+Left/Right |
| Select to line start/end | Shift+Home/End |
| Select to file start/end | Ctrl+Shift+Home/End |

### Options
| Action | Shortcut |
|--------|----------|
| Toggle Line Numbers | Ctrl+L |

## Menu Navigation

- **F10** or click to open File menu
- **Alt+F** File, **Alt+E** Edit, **Alt+S** Search, **Alt+O** Options, **Alt+H** Help
- Arrow keys to navigate within menus
- Press underlined letter to select item
- Enter to select, Escape to close

## Status Bar

The status bar shows:
- Filename (with * if modified)
- Word count (W:xxx)
- Character count (C:xxx)
- Current line and column
- File encoding (UTF-8)

## Why "Festivus"?

> "A Festivus for the rest of us!"

Named after the holiday from Seinfeld, because every text editor tries to be Vim or Emacs. Festivus is for the rest of us who just want to edit text.

## Built With

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [go-runewidth](https://github.com/mattn/go-runewidth) - Unicode width calculation
- [Chroma](https://github.com/alecthomas/chroma) - Syntax highlighting

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Feel free to submit issues and pull requests.

---

*"I got a lot of problems with you people!"* - Frank Costanza
