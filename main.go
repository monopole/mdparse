package main

import (
	"bytes"
	_ "embed"
	"fmt"
	gomAst "github.com/gomarkdown/markdown/ast"
	"github.com/monopole/mdparse/internal/file"
	"github.com/monopole/mdparse/internal/parse"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"os"
)

//go:embed testdata/small.md
var mds string

const (
	version     = "v0.2.2"
	shortHelp   = "Clone or rebase the repositories specified in the input file."
	doMyStuff   = false
	useGoldmark = true
)

func newCommand() *cobra.Command {
	//var body []byte
	return &cobra.Command{
		Use:   "mdparse {fileName}",
		Short: shortHelp,
		Long:  shortHelp + " " + version,
		Example: "  mdparse " + file.DefaultConfigFileName() + `

`,
		//Args: func(_ *cobra.Command, args []string) (err error) {
		//	filePath, err := file.GetFilePath(args)
		//	body, err = os.ReadFile(string(filePath))
		//	if err != nil {
		//		return fmt.Errorf("unable to read config file %q", filePath)
		//	}
		//	return err
		//},

		//RunE: doItWithGoMarkDown,

		RunE: whichImpl(),

		SilenceUsage: true,
	}
}

func whichImpl() func(cmd *cobra.Command, args []string) (err error) {
	// https://github.com/yuin/goldmark
	// GOOD:
	//   - One active, dedicated maintainer.
	//   - lots of extensions, proven framework.
	//   - goldmark is now the markdown renderer for Hugo, replacing blackfriday
	//   - It already supports mermaid via an extension.
	//   - It has 81 releases!  https://github.com/yuin/goldmark/releases
	//     The latest on Oct 28 2023.
	//   - 83% coverage
	//
	// BAD:
	//   - There are some PRs being ignored by the maintainer.
	//   - It doesn't yet support block level attributes, but is thinking about it
	//
	if useGoldmark {
		return doItWithGoldMark
	}

	// https://github.com/gomarkdown/markdown/graphs/contributors
	// GOOD:
	//   - It has no open pull requests (responsive owners)
	//   - Much better documentation than goldmark.
	//   - Clear access to the AST, as the API requires you to hold it
	//     in between
	//   - The AST has all the document contents.
	//   - It supports block level attributes: {#id3 .myclass fontsize="tiny"}' on (at least)
	//     header blocks and code blocks, which is all i need.
	//
	// BAD
	//   - It could support mermaid : https://github.com/gomarkdown/markdown/issues/284, but I
	//     don't seen an extension.
	//   - The number of contributors is unclear, since it is a fork of blackfriday.
	//   - It has zero official releases.
	return doItWithGoMarkDown
}

func doItWithGoMarkDown(cmd *cobra.Command, args []string) (err error) {
	md := []byte(mds)
	p := parse.NewMarkdownParser(doMyStuff)
	doc := p.Parse(md)

	gomAst.PrintWithPrefix(os.Stdout, doc, "  ")
	myWalk(doc)
	//	_, err = fmt.Printf("--- Markdown:\n%s\n\n", md)
	//_, err = fmt.Printf("--- HTML:\n%s\n", render.RenderAsHtml(doc, doMyStuff))
	return
}
func doItWithGoldMark(cmd *cobra.Command, args []string) (err error) {
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

	md := []byte(mds)
	var doc ast.Node
	doc = markdown.Parser().Parse(text.NewReader(md))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	fmt.Printf("%T %+v\n", doc, doc)
	doc.Dump(md, 0)
	// Dump and Render need the original source text because the AST doesn't
	// hold the original text - it just has byte array offsets.
	// Every Node is a BaseBlock, and each BaseBlock has a ptr to the lines in the
	// source text that make it, and each line is a Segment, and each Segment has
	// a Start and Stop integer index meant for use with a byte array.
	// I confirmed this by sending some different document in.

	doc.Dump(md, 0)

	var b bytes.Buffer
	err = markdown.Renderer().Render(&b, md, doc)
	fmt.Println(b.String())
	return
}

// General plan:
// The goal here is to leverage someone else's work in both lexing into an AST,
// and rendering an AST, so that we don't have to implement, say, mermaid handling.
//
//

// The trick with either implementation is to
//   - write a special 'code block comment' parser, creating a new AST entry
//   - modify the AST tree to create special, runnable, identifiable, codeblocks with
//     additional parameters (name, unique index, etc.)
//     -
//
// as the latter
// seems to be one ddue
func main() {
	if err := newCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
