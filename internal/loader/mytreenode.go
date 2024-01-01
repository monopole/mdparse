package loader

import (
	"path/filepath"
)

type MyTreeNode interface {
	Parent() MyTreeNode
	Name() string
	FullName() string
	Root() MyTreeNode
	Accept(TreeVisitor)
}

// TreeVisitor has the ability to visit the items specified in its methods.
type TreeVisitor interface {
	VisitFile(*MyFile)
	VisitFolder(*MyFolder)
}

// myTreeNode is the commonality between a file and a folder
type myTreeNode struct {
	parent MyTreeNode
	name   string
}

var _ MyTreeNode = &myTreeNode{}

// Root returns the "highest" non-nil tree item.
func (ti *myTreeNode) Root() MyTreeNode {
	if ti == nil {
		return nil
	}
	if ti.parent == nil {
		// This is how it stops.
		return ti
	}
	return ti.parent.Root()
}

// Name is the base name of the item.
func (ti *myTreeNode) Name() string {
	if ti == nil {
		return ""
	}
	return ti.name
}

// FullName is the fully qualified name of the item, including parents.
func (ti *myTreeNode) FullName() string {
	if ti == nil {
		return rootSlash
	}
	if ti.parent == nil {
		return ti.name
	}
	return filepath.Join(ti.parent.FullName(), ti.name)
}

// Parent is the parent of the item.
func (ti *myTreeNode) Parent() MyTreeNode {
	return ti.parent
}

func (ti *myTreeNode) Accept(_ TreeVisitor) {
	// do nothing for now
}
