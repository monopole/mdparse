package usegold

import (
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/usegold/model"
	"github.com/monopole/mdrip/base"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

func (gm *marker) Load(set *base.DataSet) error {
	panic(`
	We need a "loader" to convert the dataset to a Tutorial.
	 doc := gm.p.Parser().Parse(text.NewReader(rawData))

the argument to the command should be an array of datasoruces
andy GH urls in those shuld be convert to file trees.
so ultimately we always pass an array holding absolute paths,
some of which are directories on the local machine, and some
of those being cloned from a repo.
So we need clean data for that.
A file or a directory can come from a repo.
A sub dir comes from the place the parent dir came from.

`)
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	// gm.acc.Accumulate(doc, rawData)
	return nil
}

func (gm *marker) Dump() {
	gm.acc.Dump()
}

func (gm *marker) Render() (string, error) {
	return gm.acc.Render(gm.p.Renderer())
}
