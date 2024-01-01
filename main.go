package main

import (
	_ "embed"
	"fmt"
	"github.com/monopole/mdparse/internal/loader"
	"github.com/spf13/afero"
	"os"

	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/useblue"
	"github.com/monopole/mdparse/internal/usegold"
	"github.com/spf13/cobra"
)

//go:embed internal/usegold/accum/testdata/small.md
var mds string

const (
	version   = "v0.2.2"
	shortHelp = "Clone or rebase the repositories specified in the input file."
	doMyStuff = false
)

// General plan:
//
// Leverage someone else's work in both lexing into an AST,
// and rendering an AST, so that we don't have to implement,
// say, mermaid handling.
//
// Goal is to wrap every code
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

func newCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "mdparse {fileName}",
		Short:   shortHelp,
		Long:    shortHelp + " " + version,
		Example: "  mdparse some/directory",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var fld *loader.MyFolder
			fld, err = loadData(args)
			if err != nil {
				return err
			}
			usegold.NewVisitorDump2().VisitFolder(fld)
			var m ifc.Marker
			if useGoldmark := true; useGoldmark {
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
				m = usegold.NewMarker(doMyStuff)
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
				m = useblue.NewMarker(doMyStuff)
			}
			// try:   go run . internal/loader/README.md
			if err = m.Load([]byte(mds)); err != nil {
				return
			}
			m.Dump()
			var s string
			s, err = m.Render()
			if err != nil {
				return
			}
			if printIt := false; printIt {
				fmt.Println(s)
			}
			return
		},

		SilenceUsage: true,
	}
}

func loadData(args []string) (*loader.MyFolder, error) {
	ldr := loader.NewFsLoader(afero.NewOsFs())
	if len(args) < 2 {
		arg := "." // By default, read the current directory.
		if len(args) == 1 {
			arg = args[0]
		}
		return ldr.LoadTree(arg)
	}
	// Make one folder to hold all the argument folders.
	wrapper := loader.NewFolder("multiArgWrapper")
	for i := range args {
		fld, err := ldr.LoadTree(args[i])
		if err != nil {
			return nil, err
		}
		wrapper.AddFolderObject(fld)
	}
	return wrapper, nil
}
