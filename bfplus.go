package bfplus

import (
	"io"

	bf "github.com/russross/blackfriday/v2"
)

// Renderer extends a base renderer, allowing to add a custom anchor
// element inside every heading block with a matching `a-id` ID and an `anchor`
// class for further CSS customization
type Renderer struct {
	Base bf.Renderer // Base renderer, the one we will extend

	Highlighter *CodeRenderer
	Anchor      *AnchorRenderer
	Admonition  *AdmonitionRenderer
}

// Option defines the functional option type
type Option func(r *Renderer)

// NewRenderer returns a new renderer with all the options applied
func NewRenderer(options ...Option) *Renderer {
	r := &Renderer{}
	for _, opt := range options {
		opt(r)
	}
	return r
}

// Extend allows to specify the blackfriday renderer which is extended
func Extend(br bf.Renderer) Option {
	return func(r *Renderer) {
		r.Base = br
	}
}

// RenderNode satisfies the Renderer interface
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Heading:
		if r.Anchor != nil {
			r.Anchor.RenderNode(w, node, entering)
		}
		return r.Base.RenderNode(w, node, entering)
	case bf.CodeBlock:
		if r.Highlighter != nil {
			if err := r.Highlighter.RenderWithChroma(w, node.Literal, node.CodeBlockData); err != nil {
				return r.Base.RenderNode(w, node, entering)
			}
			return bf.SkipChildren
		}
		return r.Base.RenderNode(w, node, entering)
	default:
		if r.Admonition != nil {
			return r.Admonition.RenderNode(w, node, entering, r.Base)
		}
		return r.Base.RenderNode(w, node, entering)
	}
}

// RenderHeader satisfies the Renderer interface
func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
	r.Base.RenderHeader(w, ast)
}

// RenderFooter satisfies the Renderer interface
func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
	r.Base.RenderFooter(w, ast)
}
