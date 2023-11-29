package main

import (
	_ "embed"
	"fmt"
	"github.com/monopole/mdparse/internal/file"
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/parseblue"
	"github.com/monopole/mdparse/internal/parsegold"
	"github.com/spf13/cobra"
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

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var m ifc.Marker
			if useGoldmark {
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
				m = parsegold.NewMarkdownParser(doMyStuff)
			} else {
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
				m = parseblue.NewMarkdownParser(doMyStuff)
			}
			if err = m.Parse([]byte(mds)); err != nil {
				return
			}
			var s string
			s, err = m.Render()
			if err != nil {
				return
			}
			fmt.Println(s)
			return
		},

		SilenceUsage: true,
	}
}

// General plan:
//
// Leverage someone else's work in both lexing into an AST,
// and rendering an AST, so that we don't have to implement,
// say, mermaid handling.
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
