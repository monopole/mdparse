package accum

import (
	"github.com/monopole/mdrip/base"
)

// Tutorial represents a book hierarchical form.
type Tutorial interface {
	Accept(v TutVisitor)
	Title() string
	Name() string
	Path() base.FilePath
	Children() []Tutorial
}

// TutVisitor has the ability to visit the items specified in its methods.
type TutVisitor interface {
	VisitCourse(*Course)
	VisitLessonDoc(*LessonDoc)
	VisitCodeBlock(*CodeBlock)
}
