package main

import (
	_ "embed"
	"github.com/gomarkdown/markdown/ast"
	"github.com/monopole/mdparse/internal/file"
	"github.com/spf13/cobra"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed hoser.md
var mds string

const (
	version   = "v0.2.2"
	shortHelp = "Clone or rebase the repositories specified in the input file."
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

			extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
			p := parser.NewWithExtensions(extensions)
			doc := p.Parse(md)

			ast.PrintWithPrefix(os.Stdout, doc, "  ")
			//	_, err = fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, renderAsHtml(doc))
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

func renderAsHtml(doc ast.Node) []byte {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}
