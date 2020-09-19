package bfplus

import (
	"fmt"
	"io"
	"strings"

	bf "gopkg.in/russross/blackfriday.v2"
)

// WithHeadingAnchors enables anchor insertion for headings
func WithHeadingAnchors(options ...AnchorOption) Option {
	return func(r *Renderer) {
		r.Anchor = &AnchorRenderer{
			IDPrefix: "a-",
			Classes:  []string{"anchor"},
			Content:  "#",
		}
		for _, opt := range options {
			opt(r.Anchor)
		}
	}
}

// AnchorRenderer contains the required metadata for proper code rendering
type AnchorRenderer struct {
	IDPrefix string
	Classes  []string
	Content  string
}

// AnchorOption defines the functional option type for the code rendering
type AnchorOption func(r *AnchorRenderer)

// IDPrefix is a prefix that will be put before the ID of the matching heading
// content. Although you can set it to "", that will duplicate IDs in your DOM,
// so it not a good practice. Always use a prefix.
// Default: "a-"
func IDPrefix(s string) AnchorOption {
	return func(r *AnchorRenderer) {
		r.IDPrefix = s
	}
}

// Classes is a list of classes that will be applied to the anchor.
// Default: []string{"anchor"}
func Classes(in ...string) AnchorOption {
	return func(r *AnchorRenderer) {
		r.Classes = in
	}
}

// Content is what will be displayed. This can be any HTML code, or a simple
// char.
// Default: "#"
func Content(s string) AnchorOption {
	return func(r *AnchorRenderer) {
		r.Content = s
	}
}

// RenderNode is a simple function that needs to be called at the right moment
// when rendering markdown. Note that the AnchorRenderer itself doesn't
// implement all the methods required to be considered a true renderer. This is
// by design.
func (a *AnchorRenderer) RenderNode(w io.Writer, node *bf.Node, entering bool) {
	if !entering {
		w.Write([]byte(
			fmt.Sprintf(
				` <a id="%s%s" class="%s" href="#%s">%s</a>`,
				a.IDPrefix,
				node.HeadingID,
				strings.Join(a.Classes, " "),
				node.HeadingID,
				a.Content,
			),
		))
	}
}
