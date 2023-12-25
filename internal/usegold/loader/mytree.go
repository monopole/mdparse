package loader

import (
	"fmt"
	"path/filepath"
	"strings"
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
	VisitRepo(*MyRepo)
	VisitFile(*MyFile)
	VisitFolder(*MyFolder)
}

type myTreeItem struct {
	parent MyTreeItem
	name   string
}

/*
To test all this, we need to be able to create MyFile and MyFolder
instances from strings that contain these paths, and compare them
to structures built programmatical
(e.g f = newFile, d = newFolder, d.add(f), etc.)
so we need some simple Equals methods.
*/

// MakeTreeItem returns a MyTreeItem representing the argument.
func MakeTreeItem(fsl *FsLoader, path string) (result MyTreeItem, err error) {
	if path == "" {
		err = fmt.Errorf("need some kind of URL or path to a file or folder")
		return
	}
	if fsl == nil {
		err = fmt.Errorf("need a File System loader")
		return
	}
	if smellsLikeGithubCloneArg(path) {
		return absorbRepo(fsl, path)
	}
	cleanPath := filepath.Clean(path)
	// Disallow paths that start with "..", because the notion
	// of the root name is ambiguous, and we use the root name as a
	// user visible title.  The user's *intent* isn't clear.
	if strings.HasPrefix(cleanPath, "..") {
		err = fmt.Errorf("specify absolute path or something at or below your working directory")
		return
	}
	return fsl.LoadFolder(path)
}

func absorbRepo(fsl *FsLoader, arg string) (*MyRepo, error) {
	n, p, err := extractGithubRepoName(arg)
	if err != nil {
		return nil, err
	}
	r := &MyRepo{
		name: n,
		path: p,
	}
	err = r.Init(fsl)
	return r, err
}

func (ti *myTreeItem) IsRoot() bool {
	return ti.Root() == ti
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
		return rootSlash
	}
	if ti.parent == nil {
		return ti.name
	}
	return filepath.Join(ti.parent.FullName(), ti.name)
}

func (ti *myTreeItem) Accept(_ TreeVisitor) {
	// do nothing for now
}
