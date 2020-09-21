module github.com/Depado/bfplus

go 1.14

require (
	github.com/alecthomas/chroma v0.8.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
