package loader

import (
	"path/filepath"
)

func myScanFolder(path string) (*MyFolder, error) {
	return nil, nil
}

type MyTreeItem interface {
	Parent() MyTreeItem
	Name() string
	FullName() string
	DirName() string
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

type MyGitFolder struct {
	repo  string
	files []*MyFile
	dirs  []*MyFolder
}

func (mgf *MyGitFolder) DirName() string {
	return mgf.repo + ":::"
}

func (mgf *MyGitFolder) AbsPath() string {
	return mgf.repo + "://"
}
