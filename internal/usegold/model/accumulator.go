package accum

import (
	"github.com/yuin/goldmark/renderer"
)

// DocAccumulator accumulates document ASTs, looking for fenced code blocks.
type DocAccumulator struct {
	// All the documents that this DocAccumulator has accumulated.
	docs []*LessonDoc
}

// Accumulate accepts a document and examines it for fenced code blocks.
func (dac *DocAccumulator) Accumulate(ld *LessonDoc) {
	dac.docs = append(dac.docs, ld)
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
