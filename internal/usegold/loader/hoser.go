package loader

import (
	"path/filepath"
)

type absPath string

type MyTreeItem interface {
	Parent() MyTreeItem
	Name() string
	AbsPath() string
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

func (ti *myTreeItem) AbsPath() string {
	if ti == nil {
		return string(filepath.Separator)
	}
	if ti.parent == nil {
		return string(filepath.Separator) + ti.name
	}
	return filepath.Join(ti.parent.AbsPath(), ti.name)
}

func (ti *myTreeItem) DirName() string {
	if ti == nil {
		return string("{ERROR}")
	}
	if ti.parent == nil {
		return string(filepath.Separator)
	}
	return ti.parent.AbsPath()
}

// MyFile is named byte array.
type MyFile struct {
	myTreeItem
	content []byte
}

// MyFolder is a named group of files and folders.
type MyFolder struct {
	myTreeItem
	files []*MyFile
	dirs  []*MyFolder
}

func (fl *MyFolder) AddFile(file *MyFile) {
	file.parent = fl
	fl.files = append(fl.files, file)
}

func (fl *MyFolder) AddFolder(folder *MyFolder) {
	folder.parent = fl
	fl.dirs = append(fl.dirs, folder)
}

// MyContrivedFolder is a named grouping of files and folders
// that doesn't correspond to a "real" folder.
// Its children don't know that it exists.
type MyContrivedFolder struct {
	name  string
	items []string
	files []*MyFile
	dirs  []*MyFolder
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
