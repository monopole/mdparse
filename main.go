package main

import (
	_ "embed"
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/monopole/mdparse/internal/file"
	"github.com/monopole/mdparse/internal/parse"
	"github.com/monopole/mdparse/internal/render"
	"github.com/spf13/cobra"
	"os"
)

//go:embed small.md
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			md := []byte(mds)

			p := parse.NewMarkdownParser(doMyStuff)
			doc := p.Parse(md)

			ast.PrintWithPrefix(os.Stdout, doc, "  ")
			//myWalk(doc)
			//	_, err = fmt.Printf("--- Markdown:\n%s\n\n", md)
			_, err = fmt.Printf("--- HTML:\n%s\n", render.RenderAsHtml(doc, doMyStuff))
			return
		},
		SilenceUsage: true,
	}
}

func main() {
	if err := newCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
