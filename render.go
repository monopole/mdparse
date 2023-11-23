package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"io"
)

func renderAsHtml(doc ast.Node) []byte {

	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}

	if doMyStuff {
		opts.RenderNodeHook = myRenderHook
	}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

// a very dummy render hook that will output "code_replacements" instead of
// <code>${content}</code> emitted by html.Renderer
func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node.(type) {
	case *ast.CodeBlock:
		io.WriteString(w, "code_replacement")
		return ast.GoToNext, true
	case *Gallery:
		if entering {
			// note: just for illustration purposes
			// actual implemenation of gallery in HTML / JavaScript is long
			io.WriteString(w, "\n<gallery></gallery>\n\n")
		}
		return ast.GoToNext, true
	default:
		return ast.GoToNext, false
	}
}

// a very dummy render hook that will output "code_replacements" instead of
// <code>${content}</code> emitted by html.Renderer
func renderHookCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	_, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.GoToNext, false
	}
	io.WriteString(w, "code_replacement")
	return ast.GoToNext, true
}

func galleryRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if _, ok := node.(*Gallery); ok {
		if entering {
			// note: just for illustration purposes
			// actual implemenation of gallery in HTML / JavaScript is long
			io.WriteString(w, "\n<gallery></gallery>\n\n")
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
