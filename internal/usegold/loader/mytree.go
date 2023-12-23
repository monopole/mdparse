package loader

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type MyTreeItem interface {
	Parent() MyTreeItem
	Name() string
	FullName() string
	DirName() string
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
		fsl = NewFsLoader(afero.NewOsFs())
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

	if filepath.IsAbs(cleanPath) {

		// return either a file or a folder
		// if it's a file, the parent of the file is a SINGLE folder
		// with the full rooted path, e.g.
		// "/" or "/home" or "/home/bob"
		// the contents of the JUST the final file as a child of the folder object
		// if folder, then same as above but
		// ALL contents of "/" or "/home" or whatever are loaded
		// into the single folder.

		var folder MyFolder
		folder.name = string(filepath.Separator)
		err = folder.AbsorbFolderFromDisk(fsl, cleanPath)
		return &folder, err
	}

	if strings.HasPrefix(cleanPath, "./") {
		// The user's intent is clear, but we don't need this
		// in the final result.
		cleanPath = cleanPath[2:]
	}
	if cleanPath == "." {
		// load the current working directory and call the encapsulating
		// folder ".".
	} else {
		// Load the rightmost thing - either a folder or a file
		// and the parent folder should have cleanPath as its name.
	}
	{
		var cwd string
		cwd, err = os.Getwd()
		if err != nil {
			return
		}
		var folder MyFolder
		folder.name = stripTrailingSlash(cwd)
	}
	//dir, name := FSplit(path)
	//folder := fl
	//if dir != "" {
	//	folder = fl.buildParentTree(dir)
	//}
	//if err = folder.AbsorbFolderFromDisk(fsl, cleanPath); err != nil {
	//	return
	//}

	// TODO UNCOMMENT THIS!!! DO OT DELETE TILL YOU HAVE REPOS WORKING
	//var info os.FileInfo
	//info, err = os.Stat(absPath)
	//if err != nil {
	//	return
	//}
	//if info.IsDir() {
	//	if fsl.IsAllowedFolder(info) {
	//		if filepath.IsAbs(path) {
	//			return m.folderAbs.AbsorbFolderFromDisk(m.fsl, path)
	//		}
	//		return m.folderRel.AbsorbFolderFromDisk(m.fsl, path)
	//	}
	//	return fmt.Errorf("illegal folder %q", info.Name())
	//}
	//if m.fsl.IsAllowedFile(info) {
	//	if filepath.IsAbs(path) {
	//		return m.folderAbs.AbsorbFileFromDisk(path)
	//	}
	//	return m.folderRel.AbsorbFileFromDisk(path)
	//}
	//return fmt.Errorf("not a markdown file %q", info.Name())

	return nil, nil
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
	if err = r.Init(fsl); err != nil {
		return nil, err
	}
	return r, nil
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
