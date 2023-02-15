package bfplus

import (
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	bf "github.com/russross/blackfriday/v2"
)

// WithCodeHighlighting enables syntax highlighting using Chroma
func WithCodeHighlighting(options ...CodeOption) Option {
	return func(r *Renderer) {
		r.Highlighter = &CodeRenderer{
			Style:      styles.Get("monokai"),
			Autodetect: true,
		}
		for _, option := range options {
			option(r.Highlighter)
		}
		r.Highlighter.Formatter = html.New(r.Highlighter.ChromaOptions...)
	}
}

// CodeRenderer contains the required metadata for proper code rendering
type CodeRenderer struct {
	Autodetect    bool
	ChromaOptions []html.Option
	Style         *chroma.Style
	Formatter     *html.Formatter
}

// CodeOption defines the functional option type for the code rendering
type CodeOption func(r *CodeRenderer)

// Style is a function option allowing to set the style used by chroma
// Default : "monokai"
func Style(s string) CodeOption {
	return func(r *CodeRenderer) {
		r.Style = styles.Get(s)
	}
}

// ChromaStyle is an option to directly set the style of the renderer using a
// chroma style instead of a string
func ChromaStyle(s *chroma.Style) CodeOption {
	return func(r *CodeRenderer) {
		r.Style = s
	}
}

// WithoutAutodetect disables chroma's language detection when no codeblock
// extra information is given. It will fallback to a sane default instead of
// trying to detect the language.
func WithoutAutodetect() CodeOption {
	return func(r *CodeRenderer) {
		r.Autodetect = false
	}
}

// ChromaOptions allows to pass Chroma html.Option such as Standalone()
// WithClasses(), ClassPrefix(prefix)...
func ChromaOptions(options ...html.Option) CodeOption {
	return func(r *CodeRenderer) {
		r.ChromaOptions = options
	}
}

// RenderWithChroma will render the given text to the w io.Writer
func (r *CodeRenderer) RenderWithChroma(w io.Writer, text []byte, data bf.CodeBlockData) error {
	var lexer chroma.Lexer

	// Determining the lexer to use
	if len(data.Info) > 0 {
		lexer = lexers.Get(string(data.Info))
	} else if r.Autodetect {
		lexer = lexers.Analyse(string(text))
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, string(text))
	if err != nil {
		return err
	}
	return r.Formatter.Format(w, r.Style, iterator)
}
