package render

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/monopole/mdparse/internal/parseblue"
	"io"
)

func RenderAsHtml(doc ast.Node, doMyStuff bool) []byte {
	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}
	if doMyStuff {
		opts.RenderNodeHook = myRenderHook
	}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node.(type) {
	case *ast.CodeBlock:
		if entering {
			io.WriteString(w, "code_replacement\n")
		}
		return ast.GoToNext, true
	case *parseblue.Gallery:
		if entering {
			io.WriteString(w, "\n<gallery></gallery>\n\n")
		}
		return ast.GoToNext, true
	default:
		return ast.GoToNext, false
	}
}
