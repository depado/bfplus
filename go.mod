module github.com/Depado/bfplus

go 1.12

require (
	github.com/alecthomas/chroma v0.6.6
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
