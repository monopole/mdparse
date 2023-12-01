package usegold

import (
	bytes2 "bytes"
	"fmt"
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"strings"
)

type gomark struct {
	doMyStuff bool
	p         goldmark.Markdown
	rawData   []byte
	doc       ast.Node
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
	return &gomark{doMyStuff: doMyStuff, p: markdown}
}

func (gm *gomark) Parse(bytes []byte) error {
	gm.rawData = bytes

	gm.doc = gm.p.Parser().Parse(text.NewReader(bytes))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	fmt.Printf("%T %+v\n", gm.doc, gm.doc)
	// Dump and Render need the original source text because the AST doesn't
	// hold the original text - it just has byte array offsets.
	// Every Node is a BaseBlock, and each BaseBlock has a ptr to the lines in the
	// source text that make it, and each line is a Segment, and each Segment has
	// a Start and Stop integer index meant for use with a byte array.
	// I confirmed this by sending some different document in.

	return gm.WalkIt()
}

func (gm *gomark) Dump() {
	gm.doc.Dump(gm.rawData, 0)
}

func (gm *gomark) Render() (string, error) {
	var b bytes2.Buffer
	err := gm.p.Renderer().Render(&b, gm.rawData, gm.doc)
	return b.String(), err
}

func (gm *gomark) WalkIt() error {
	gm.depth = 0
	return ast.Walk(gm.doc, gm.myWalk)
}

const blanks = "                      "

func (gm *gomark) myWalk(n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		gm.depth--
		return ast.WalkContinue, nil
	}
	gm.depth++
	s := string(n.Text(gm.rawData))
	if len(s) > 30 {
		s = s[:30] + "..."
	}
	if n.Kind() == ast.KindFencedCodeBlock {
		prev := n.PreviousSibling()
		if prev != nil && prev.Kind() == ast.KindHTMLBlock {
			htmlBlock, ok := prev.(*ast.HTMLBlock)
			var labels []string
			if ok {
				labels = recoverLabels(rawText(htmlBlock, gm.rawData))
				for i := range labels {
					fmt.Printf("  %q\n", labels[i])
					n.SetAttributeString(labels[i], "")
				}
			}
		}
		fmt.Println("fencedCodeBlock")
		fmt.Printf("  %q\n", rawText(n, gm.rawData))
	}
	// fmt.Printf("%s k=%30s t=%d %s \n", blanks[:gm.depth], n.Kind(), n.Type(), s)
	return ast.WalkContinue, nil
}

func rawText(n ast.Node, raw []byte) string {
	var buff strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		s := n.Lines().At(i)
		buff.Write(raw[s.Start:s.Stop])
	}
	return buff.String()
}

func commentBody(s string) string {
	const (
		begin = "<!--"
		end   = "-->"
	)
	if !strings.HasPrefix(s, begin) {
		return ""
	}
	if !strings.HasSuffix(s, end) {
		return ""
	}
	return s[len(begin) : len(s)-len(end)]
}

func recoverLabels(s string) (result []string) {
	if !strings.HasPrefix(s, "<!--") {
		// Ignore the HTML if it isn't a comment.
		//
		return
	}
	fmt.Println(s)
	items := strings.Split(s, " ")
	for _, word := range items {
		if len(word) > 0 && word[0] == uint8('@') {
			result = append(result, word[1:])
		}
	}
	return
}
