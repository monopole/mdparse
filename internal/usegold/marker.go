package usegold

import (
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/usegold/accum"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type marker struct {
	doMyStuff bool
	p         goldmark.Markdown
	acc       *accum.DocAccumulator
	depth     int
}

func NewMarker(doMyStuff bool) ifc.Marker {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	return &marker{doMyStuff: doMyStuff, p: markdown, acc: &accum.DocAccumulator{}}
}

func (gm *marker) Parse(rawData []byte) error {

	doc := gm.p.Parser().Parse(text.NewReader(rawData))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)

	return gm.acc.Accumulate(doc, rawData)
}

func (gm *marker) Dump() {
	gm.acc.Dump()
}

func (gm *marker) Render() (string, error) {
	return gm.acc.Render(gm.p.Renderer())
}
