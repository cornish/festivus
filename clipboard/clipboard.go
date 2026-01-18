package clipboard

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aymanbagabas/go-osc52/v2"
)

// ClipboardTool represents an available clipboard tool
type ClipboardTool int

const (
	ToolNone ClipboardTool = iota
	ToolXclip
	ToolXsel
	ToolWlClipboard
)

// Clipboard provides unified clipboard access with OSC52 support for SSH.
type Clipboard struct {
	// Internal clipboard for when no system clipboard is available
	internal string
	// Whether we're likely in an SSH session
	isSSH bool
	// Output writer for OSC52 sequences (typically os.Stdout)
	output io.Writer
	// Detected clipboard tool
	tool ClipboardTool
	// Whether we've warned about missing clipboard tools
	warned bool
}

// New creates a new Clipboard instance.
func New(output io.Writer) *Clipboard {
	if output == nil {
		output = os.Stdout
	}
	return &Clipboard{
		isSSH:  isSSHSession(),
		output: output,
		tool:   detectClipboardTool(),
	}
}

// isSSHSession detects if we're running in an SSH session.
func isSSHSession() bool {
	// Check common SSH environment variables
	if os.Getenv("SSH_TTY") != "" {
		return true
	}
	if os.Getenv("SSH_CLIENT") != "" {
		return true
	}
	if os.Getenv("SSH_CONNECTION") != "" {
		return true
	}
	return false
}

// detectClipboardTool finds an available clipboard tool
func detectClipboardTool() ClipboardTool {
	// Check for Wayland first if WAYLAND_DISPLAY is set
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if _, err := exec.LookPath("wl-copy"); err == nil {
			if _, err := exec.LookPath("wl-paste"); err == nil {
				return ToolWlClipboard
			}
		}
	}

	// Check for X11 tools
	if os.Getenv("DISPLAY") != "" {
		if _, err := exec.LookPath("xclip"); err == nil {
			return ToolXclip
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return ToolXsel
		}
	}

	return ToolNone
}

// Copy copies the given text to the clipboard.
// In SSH sessions, it uses OSC52 escape sequences.
// Locally, it tries native clipboard tools first.
func (c *Clipboard) Copy(text string) error {
	// Always store internally as a last resort
	c.internal = text

	if c.isSSH {
		// In SSH, always use OSC52
		return c.copyOSC52(text)
	}

	// Try native clipboard tool
	err := c.copyNative(text)
	if err == nil {
		return nil
	}

	// Fall back to OSC52
	return c.copyOSC52(text)
}

// copyNative copies text using native clipboard tools
func (c *Clipboard) copyNative(text string) error {
	var cmd *exec.Cmd

	switch c.tool {
	case ToolXclip:
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case ToolXsel:
		cmd = exec.Command("xsel", "--clipboard", "--input")
	case ToolWlClipboard:
		cmd = exec.Command("wl-copy")
	default:
		return &ClipboardError{Message: "no clipboard tool available"}
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// copyOSC52 copies text using OSC52 escape sequence.
func (c *Clipboard) copyOSC52(text string) error {
	seq := osc52.New(text)
	_, err := io.WriteString(c.output, seq.String())
	return err
}

// Paste returns text from the clipboard.
// Note: OSC52 paste (OSC52 query) is not widely supported.
// We rely on native clipboard tools or the internal buffer.
func (c *Clipboard) Paste() (string, error) {
	// Try native clipboard tool first
	text, err := c.pasteNative()
	if err == nil && text != "" {
		return text, nil
	}

	// Fall back to internal clipboard
	return c.internal, nil
}

// pasteNative reads from clipboard using native tools
func (c *Clipboard) pasteNative() (string, error) {
	var cmd *exec.Cmd

	switch c.tool {
	case ToolXclip:
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case ToolXsel:
		cmd = exec.Command("xsel", "--clipboard", "--output")
	case ToolWlClipboard:
		cmd = exec.Command("wl-paste", "-n")
	default:
		return "", &ClipboardError{Message: "no clipboard tool available"}
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// HasContent returns true if there's content available to paste.
func (c *Clipboard) HasContent() bool {
	// Check native clipboard
	text, err := c.pasteNative()
	if err == nil && text != "" {
		return true
	}

	// Check internal clipboard
	return c.internal != ""
}

// Clear clears the internal clipboard.
func (c *Clipboard) Clear() {
	c.internal = ""
}

// IsSSH returns true if we're in an SSH session.
func (c *Clipboard) IsSSH() bool {
	return c.isSSH
}

// HasNativeClipboard returns true if a native clipboard tool is available.
func (c *Clipboard) HasNativeClipboard() bool {
	return c.tool != ToolNone
}

// ToolName returns the name of the detected clipboard tool.
func (c *Clipboard) ToolName() string {
	switch c.tool {
	case ToolXclip:
		return "xclip"
	case ToolXsel:
		return "xsel"
	case ToolWlClipboard:
		return "wl-clipboard"
	default:
		return "none"
	}
}

// ClipboardError represents a clipboard operation error
type ClipboardError struct {
	Message string
}

func (e *ClipboardError) Error() string {
	return e.Message
}
