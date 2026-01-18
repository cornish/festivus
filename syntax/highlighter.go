package syntax

import (
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// ColorSpan represents a colored region of text
type ColorSpan struct {
	Start int    // Start column (rune index)
	End   int    // End column (rune index, exclusive)
	Color string // ANSI color code
}

// Highlighter provides syntax highlighting for source code
type Highlighter struct {
	lexer   chroma.Lexer
	enabled bool
}

// New creates a new Highlighter for the given filename
func New(filename string) *Highlighter {
	h := &Highlighter{
		enabled: true,
	}
	h.SetFile(filename)
	return h
}

// SetFile updates the lexer based on the filename
func (h *Highlighter) SetFile(filename string) {
	if filename == "" {
		h.lexer = nil
		return
	}
	h.lexer = lexers.Match(filename)
	if h.lexer != nil {
		h.lexer = chroma.Coalesce(h.lexer)
	}
}

// SetEnabled enables or disables syntax highlighting
func (h *Highlighter) SetEnabled(enabled bool) {
	h.enabled = enabled
}

// Enabled returns whether highlighting is enabled
func (h *Highlighter) Enabled() bool {
	return h.enabled
}

// HasLexer returns true if a lexer is available for the current file
func (h *Highlighter) HasLexer() bool {
	return h.lexer != nil
}

// GetLineColors returns color spans for a line
// Returns nil if highlighting is disabled or no lexer is available
func (h *Highlighter) GetLineColors(line string) []ColorSpan {
	if !h.enabled || h.lexer == nil {
		return nil
	}

	iterator, err := h.lexer.Tokenise(nil, line)
	if err != nil {
		return nil
	}

	var spans []ColorSpan
	pos := 0
	for _, token := range iterator.Tokens() {
		color := tokenColor(token.Type)
		tokenLen := utf8.RuneCountInString(token.Value)
		if color != "" && tokenLen > 0 {
			spans = append(spans, ColorSpan{
				Start: pos,
				End:   pos + tokenLen,
				Color: color,
			})
		}
		pos += tokenLen
	}

	return spans
}

// ColorAt returns the color for a specific column position
// Returns empty string if no color applies
func ColorAt(spans []ColorSpan, col int) string {
	for _, span := range spans {
		if col >= span.Start && col < span.End {
			return span.Color
		}
	}
	return ""
}

// tokenColor returns the ANSI color code for a token type
func tokenColor(t chroma.TokenType) string {
	switch {
	// Keywords
	case t == chroma.Keyword,
		t == chroma.KeywordConstant,
		t == chroma.KeywordDeclaration,
		t == chroma.KeywordNamespace,
		t == chroma.KeywordPseudo,
		t == chroma.KeywordReserved,
		t == chroma.KeywordType:
		return "\033[96m" // Bright cyan

	// Strings
	case t == chroma.String,
		t == chroma.StringAffix,
		t == chroma.StringBacktick,
		t == chroma.StringChar,
		t == chroma.StringDelimiter,
		t == chroma.StringDoc,
		t == chroma.StringDouble,
		t == chroma.StringEscape,
		t == chroma.StringHeredoc,
		t == chroma.StringInterpol,
		t == chroma.StringOther,
		t == chroma.StringRegex,
		t == chroma.StringSingle,
		t == chroma.StringSymbol:
		return "\033[92m" // Bright green

	// Comments
	case t == chroma.Comment,
		t == chroma.CommentHashbang,
		t == chroma.CommentMultiline,
		t == chroma.CommentPreproc,
		t == chroma.CommentPreprocFile,
		t == chroma.CommentSingle,
		t == chroma.CommentSpecial:
		return "\033[90m" // Bright black (gray)

	// Numbers
	case t == chroma.Number,
		t == chroma.NumberBin,
		t == chroma.NumberFloat,
		t == chroma.NumberHex,
		t == chroma.NumberInteger,
		t == chroma.NumberIntegerLong,
		t == chroma.NumberOct:
		return "\033[93m" // Bright yellow

	// Operators
	case t == chroma.Operator,
		t == chroma.OperatorWord:
		return "\033[97m" // Bright white

	// Functions
	case t == chroma.NameFunction,
		t == chroma.NameFunctionMagic:
		return "\033[94m" // Bright blue

	// Types/Classes
	case t == chroma.NameClass,
		t == chroma.NameBuiltin,
		t == chroma.NameBuiltinPseudo:
		return "\033[95m" // Bright magenta

	// Constants
	case t == chroma.NameConstant:
		return "\033[93m" // Bright yellow

	// Preprocessor
	case t == chroma.CommentPreproc,
		t == chroma.GenericHeading,
		t == chroma.GenericSubheading:
		return "\033[95m" // Bright magenta

	// Errors
	case t == chroma.Error,
		t == chroma.GenericError:
		return "\033[91m" // Bright red

	default:
		return "" // Default terminal color
	}
}
