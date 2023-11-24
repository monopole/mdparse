package main

import (
	_ "embed"
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/monopole/mdparse/internal/file"
	"github.com/monopole/mdparse/internal/parse"
	"github.com/monopole/mdparse/internal/render"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"os"
)

//go:embed internal/testdata/small.md
var mds string

const (
	version   = "v0.2.2"
	shortHelp = "Clone or rebase the repositories specified in the input file."
	doMyStuff = true
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

		// RunE:         doItWithGoMarkDown,
		RunE: doItWithGoldMark,

		SilenceUsage: true,
	}
}

func doItWithGoMarkDown(cmd *cobra.Command, args []string) (err error) {
	md := []byte(mds)
	p := parse.NewMarkdownParser(doMyStuff)
	doc := p.Parse(md)

	ast.PrintWithPrefix(os.Stdout, doc, "  ")
	//myWalk(doc)
	//	_, err = fmt.Printf("--- Markdown:\n%s\n\n", md)
	_, err = fmt.Printf("--- HTML:\n%s\n", render.RenderAsHtml(doc, doMyStuff))
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

	//var b bytes.Buffer
	md := []byte(mds)

	doc := markdown.Parser().Parse(text.NewReader(md))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	fmt.Printf("%T %+v\n", doc, doc)
	//err = markdown.Renderer().Render(&b, md, doc)
	//fmt.Println(b.String())
	return
}

// General plan:
// The goal here is to leverage someone else's work in both lexing into an AST,
// and rendering an AST, so that we don't have to implement, say, mermaid handling.
//
//

// https://github.com/yuin/goldmark has lots of single person activity and lots of extensions,
// and has some PRs being ignored by the maintainer.  Might be gone.
//
// https://github.com/gomarkdown/markdown/graphs/contributors is a fork of blackfriday
// has no open pull requests, but has better documentation than goldmark, and has clear
// access to the AST.
// and seems like it could support mermaid : https://github.com/gomarkdown/markdown/issues/284
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
