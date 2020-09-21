package main

import (
	"fmt"

	bfp "github.com/Depado/bfplus"
	"github.com/alecthomas/chroma/formatters/html"
	bf "gopkg.in/russross/blackfriday.v2"
)

var exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.DefinitionLists | bf.Footnotes

var flags = bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.SmartypantsLatexDashes | bf.TOC

var md = `# Title
## Subtitle

!!! note "My Note Title"
	First Line
	Second Line
	*Italic*
	**Bold**

!!! warning
	**This is very dangerous, think again!**

!!! danger "Dangerous Stuff Ahead"
	This is a simple test.
	This could even be another test to be honest.


` + "```go" + `
package main

import "fmt"

fmt.Println("Hello World")
` + "```" + `

Let's go back to non-admonition markdown now. 
**And see if that works properly.**
`

func main() {
	r := bfp.NewRenderer(
		bfp.WithAdmonition(),
		bfp.WithCodeHighlighting(
			bfp.Style("monokai"),
			bfp.WithoutAutodetect(),
			bfp.ChromaOptions(html.WithClasses(true)),
		),
		bfp.WithHeadingAnchors(),
		bfp.Extend(bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags})),
	)
	fmt.Println(string(bf.Run(
		[]byte(md),
		bf.WithRenderer(r),
		bf.WithExtensions(exts),
	)))
}
