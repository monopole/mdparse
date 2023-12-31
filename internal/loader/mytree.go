package loader

import (
	"path/filepath"
)

type MyTreeItem interface {
	Parent() MyTreeItem
	Name() string
	FullName() string
	Root() MyTreeItem
	Accept(TreeVisitor)
}

// TreeVisitor has the ability to visit the items specified in its methods.
type TreeVisitor interface {
	VisitFile(*MyFile)
	VisitFolder(*MyFolder)
}

// myTreeItem is the commonality between a file and a folder
type myTreeItem struct {
	parent MyTreeItem
	name   string
}

// Root returns the "highest" non-nil tree item.
func (ti *myTreeItem) Root() MyTreeItem {
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
func (ti *myTreeItem) Name() string {
	if ti == nil {
		return ""
	}
	return ti.name
}

// FullName is the fully qualified name of the item, including parents.
func (ti *myTreeItem) FullName() string {
	if ti == nil {
		return rootSlash
	}
	if ti.parent == nil {
		return ti.name
	}
	return filepath.Join(ti.parent.FullName(), ti.name)
}

// Parent is the parent of the item.
func (ti *myTreeItem) Parent() MyTreeItem {
	return ti.parent
}

func (ti *myTreeItem) Accept(_ TreeVisitor) {
	// do nothing for now
}
