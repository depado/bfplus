package bfplus

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	bf "gopkg.in/russross/blackfriday.v2"
)

// WithAdmonition enables admonition support
func WithAdmonition() Option {
	return func(r *Renderer) {
		r.Admonition = &AdmonitionRenderer{
			re: regexp.MustCompile(`^!!!\s?(\w+(?: +\w+)*)(?: +"(.*?)")? *\n`),
		}
	}
}

// AdmonitionRenderer is a custom Blackfriday renderer that attempts to find admonition
// style markdown and render it
type AdmonitionRenderer struct {
	in   bool
	buff bytes.Buffer
	w    *bufio.Writer
	re   *regexp.Regexp
}

func (r *AdmonitionRenderer) check(n *bf.Node) bool {
	matches := r.re.FindSubmatch(n.Literal)
	remain := bytes.SplitN(n.Literal, []byte{'\n'}, 2)
	return matches == nil || len(remain) != 2
}

// RenderNode will render the node and try to find admonitions
//nolint:errcheck,gosec
func (r *AdmonitionRenderer) RenderNode(w io.Writer, node *bf.Node, entering bool, base bf.Renderer) bf.WalkStatus {
	// First we check if we enter a paragraph. If so, we check if the first child
	// matches with our regex so we don't generate an extra useless <p> tag
	if node.Type == bf.Paragraph && entering && node.FirstChild != nil {
		if r.check(node.FirstChild) { // This doesn't match, keep going
			return base.RenderNode(w, node, entering)
		}
		return bf.GoToNext
	}
	if !r.in {
		matches := r.re.FindSubmatch(node.Literal)
		remain := bytes.SplitN(node.Literal, []byte{'\n'}, 2)
		if matches == nil || len(remain) != 2 { // This doesn't match, keep going
			return base.RenderNode(w, node, entering)
		}
		r.in = true
		node.Literal = remain[1]
		r.buff = bytes.Buffer{}
		r.w = bufio.NewWriter(&r.buff)
		t, title := matches[1], matches[2]
		r.w.WriteString(`<div class="admonition `)
		r.w.Write(t)
		r.w.WriteString(`">`)
		r.w.Write([]byte{'\n', '\t'})
		if len(title) > 0 {
			r.w.WriteString(`<p class="admonition-title">`)
			r.w.Write(title)
			r.w.WriteString(`</p>`)
			r.w.Write([]byte{'\n', '\t'})
		}
		r.w.WriteString("<p>")
		r.w.WriteRune('\n')
		return base.RenderNode(r.w, node, entering)
	}
	if r.in && node.Type == bf.Paragraph && !entering {
		r.in = false
		r.w.Write([]byte{'\n', '\t'})
		r.w.WriteString("</p>")
		r.w.WriteRune('\n')
		r.w.WriteString("</div>")
		r.w.WriteRune('\n')
		r.w.Flush()
		r.buff.WriteTo(w)
	}

	return base.RenderNode(r.w, node, entering)
}
