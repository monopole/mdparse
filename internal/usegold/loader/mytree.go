package loader

import (
	"path/filepath"
)

type MyTreeItem interface {
	Parent() MyTreeItem
	Name() string
	FullName() string
	DirName() string
	Accept(TreeVisitor)
}

// TreeVisitor has the ability to visit the items specified in its methods.
type TreeVisitor interface {
	VisitRepo(*MyRepo)
	VisitContrivedFolder(*MyContrivedFolder)
	VisitFile(*MyFile)
	VisitFolder(*MyFolder)
}

type myTreeItem struct {
	parent MyTreeItem
	name   string
}

func (ti *myTreeItem) Name() string {
	if ti == nil {
		return ""
	}
	return ti.name
}

func (ti *myTreeItem) Parent() MyTreeItem {
	return ti.parent
}

func (ti *myTreeItem) FullName() string {
	if ti == nil {
		return string(filepath.Separator)
	}
	if ti.parent == nil {
		return string(filepath.Separator) + ti.name
	}
	return filepath.Join(ti.parent.FullName(), ti.name)
}

func (ti *myTreeItem) DirName() string {
	if ti == nil {
		return "{ERROR}"
	}
	if ti.parent == nil {
		return string(filepath.Separator)
	}
	return ti.parent.FullName()
}

func (ti *myTreeItem) Accept(_ TreeVisitor) {
	// do nothing for now
}
