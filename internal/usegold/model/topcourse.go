package accum

import (
	"github.com/monopole/mdrip/base"
)

// A TopCourse is a Course that isn't created from a real physical directory.
// It's still a list of files and directories, but they don't live in
// one place, as they were gathered from a command line containing an
// unrelated list of files, directories, URLs, etc. Its "name" field must be
// computed in some special way. There's only one of these in a tree,
// it's at the root, and it's only necessary if the tree was built from
// more than one element in a command line.
type TopCourse struct {
	Course
}

var _ Tutorial = &TopCourse{}

// NewTopCourse makes a new TopCourse.
func NewTopCourse(n string, p base.FilePath, c []Tutorial) *TopCourse {
	return &TopCourse{Course{n, p, c}}
}

// Accept accepts a visitor.
func (t *TopCourse) Accept(v TutVisitor) { v.VisitTopCourse(t) }
