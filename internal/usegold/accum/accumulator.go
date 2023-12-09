package accum

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// DocAccumulator accumulates document ASTs, looking for fenced code blocks.
type DocAccumulator struct {
	// All the documents that this DocAccumulator has accumulated.
	docs []*docHolder
}

// Accumulate accepts a document and examines it for fenced code blocks.
func (dac *DocAccumulator) Accumulate(doc ast.Node, content []byte) error {
	dh := &docHolder{doc: doc, content: content}
	dac.docs = append(dac.docs, dh)
	return ast.Walk(doc, dh.myWalk)
}

func (dac *DocAccumulator) Dump() {
	for _, dh := range dac.docs {
		dh.Dump()
	}
}

func (dac *DocAccumulator) Render(r renderer.Renderer) (string, error) {
	// Assumes only one doc and renders it
	return dac.docs[0].Render(r)
}
