package usegold

import (
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/usegold/accum"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type marker struct {
	doMyStuff bool
	p         goldmark.Markdown
	depth     int
	doc       *accum.LessonDoc
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
	return &marker{doMyStuff: doMyStuff, p: markdown}
}

func (gm *marker) Load(rawData []byte) error {
	//gm.doc = gm.p.Parser().Parse(text.NewReader(rawData))
	var err error
	gm.doc, err = accum.NewLessonDocFromBytes(gm.p.Parser(), rawData)
	if err != nil {
		return err
	}
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	//gm.acc.Accumulate(ld)
	return nil
}

func (gm *marker) Dump() {
	gm.doc.Dump()
}

func (gm *marker) Render() (string, error) {
	return gm.doc.Render(gm.p.Renderer())
}
