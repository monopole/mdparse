package usegold

import (
	"bytes"
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

func (gm *gomark) Parse(rawData []byte) error {
	gm.rawData = rawData

	gm.doc = gm.p.Parser().Parse(text.NewReader(rawData))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)

	// The AST doesn't hold the original text - it just has byte array offsets.
	// Every Node is a BaseBlock, and each BaseBlock has a ptr to the lines in the
	// source text that make it, and each line is a Segment, and each Segment has
	// a Start and Stop integer index meant for use with a byte array.
	// I confirmed this by sending some different document in.

	w := Walker{}
	// This can walk and accumulate more than one doc!
	return w.WalkDoc(gm.doc, gm.rawData)
}

func (gm *gomark) Dump() {
	gm.doc.Dump(gm.rawData, 0)
}

func (gm *gomark) Render() (string, error) {
	var b bytes.Buffer
	err := gm.p.Renderer().Render(&b, gm.rawData, gm.doc)
	return b.String(), err
}

// DocHolder associates an AST tree with the raw content that
// was parsed to form the tree.
type DocHolder struct {
	doc     ast.Node
	content []byte
}

// Walker walks documents, accumulating them as it goes.
type Walker struct {
	// Used to assign a unique ID to each code block across all the documents.
	codeBlockCounter int
	// depth is for printing
	depth int
	// All the documents that this Walker has walked.
	docs []DocHolder
}

// WalkDoc walks and accumulates a parsed document.
func (w *Walker) WalkDoc(doc ast.Node, content []byte) error {
	w.docs = append(w.docs, DocHolder{doc: doc, content: content})
	w.depth = 0
	return ast.Walk(doc, w.myWalk)
}

const blanks = "                      "

func (w *Walker) myWalk(n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.depth--
		return ast.WalkContinue, nil
	}
	w.depth++
	if n.Kind() == ast.KindFencedCodeBlock {
		w.codeBlockCounter++
		fmt.Printf("fencedCodeBlock %d\n", w.codeBlockCounter)
		labels := []string{fmt.Sprintf("fcb_%03d", w.codeBlockCounter)}
		if prev := n.PreviousSibling(); prev != nil && prev.Kind() == ast.KindHTMLBlock {
			if b, ok := prev.(*ast.HTMLBlock); ok {
				// We have an HTML block.  If it's a comment, try to extract labels.
				labels = append(labels, parseLabels(commentBody(w.nodeText(b)))...)
			}
		}
		for i := range labels {
			fmt.Printf("  %q\n", labels[i])
			// TODO: instead use a fixed key, and store the array as the value?
			n.SetAttributeString(labels[i], labels[i])
		}
		fmt.Printf("     %q\n", w.nodeText(n))
	}
	//s := string(n.Text(w.currentContent()))
	//if len(s) > 30 {
	//	s = s[:30] + "..."
	//}
	// fmt.Printf("%s k=%30s t=%d %s \n", blanks[:gm.depth], n.Kind(), n.Type(), s)
	return ast.WalkContinue, nil
}

func (w *Walker) nodeText(n ast.Node) string {
	return gatherText(n, w.currentContent())
}

func (w *Walker) currentContent() []byte {
	return w.docs[len(w.docs)-1].content
}

func gatherText(n ast.Node, raw []byte) string {
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
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, begin) {
		return ""
	}
	if !strings.HasSuffix(s, end) {
		return ""
	}
	return s[len(begin) : len(s)-len(end)]
}

const labelPrefixChar = uint8('@')

func parseLabels(s string) (result []string) {
	items := strings.Split(s, " ")
	for _, word := range items {
		i := 0
		for i < len(word) && word[i] == labelPrefixChar {
			i++
		}
		if i > 0 && i < len(word) && word[i-1] == labelPrefixChar {
			result = append(result, word[i:])
		}
	}
	return
}
