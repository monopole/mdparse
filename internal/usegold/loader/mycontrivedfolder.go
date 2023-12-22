package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type filter func(info os.FileInfo) bool

/*
TODO

This is a mess at the moment.

The thing to do is take one argument to the program.
This arg can be

  - an abs path to some file
    the parent of the file is a SINGLE folder with the full rooted path, e.g.
    "/" or "/home" or "/home/bob"
    the contents of the JUST the final file as a child of the folder object

  - an abs path to some folder
    same as above, but ALL contents of "/" or "/home" or whatever are loaded
    into the single folder.

  - a path that starts with ".."
    reject it

  - a path that starts with "./"
    strip this from the front and eval the result

  - a "." or a "./"
    make this the empty string

  - a non-empty path that doesn't start with "/"
    load the rightmost thing (either a folder or a file)
    and the parent is a folder with that path
    as its entire name (as above)

  - an empty string
    load the current working directory and call
    the encapsulating folder ".".

To test all this, we need to be able to create MyFile and MyFolder
instances from strings that contain these paths, and compare them
to structures built programmatical
(e.g f = newFile, d = newFolder, d.add(f), etc.)
so we need some simple Equals methods.
*/
func MakeTreeItem(fsl *FsLoader, path string) (result MyTreeItem, err error) {
	if path == "" {
		err = fmt.Errorf("need a path")
		return
	}
	if fsl == nil {
		fsl = DefaultFsLoader
	}

	var folder MyFolder
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		folder.name = string(filepath.Separator)
		if err = folder.AbsorbFolderFromDisk(fsl, cleanPath); err != nil {
			return
		}
	} else {
		// For now force the user to exclude anything that starts with "..";
		// they only need to cd to fix it.
		// Goal is a "clean" path for user-visible titles.
		if strings.HasPrefix(cleanPath, "..") {
			err = fmt.Errorf("specify absolute path or something at or below your working directory")
			return
		}
		if strings.HasPrefix(cleanPath, "./") {
			cleanPath = cleanPath[2:]
		}
		if cleanPath == "." {
			// We need an empty name parent
			if err = folder.AbsorbFolderFromDisk(fsl, cleanPath); err != nil {
				return
			}

		}
		{
			var cwd string
			cwd, err = os.Getwd()
			if err != nil {
				return
			}
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
	}
	return &folder, nil

	//if smellsLikeGithubCloneArg(arg) {
	//	return LoadRepo(arg)
	//}
	//var info os.FileInfo
	//info, err = os.Stat(absPath)
	//if err != nil {
	//	return
	//}
	//if info.IsDir() {
	//	if fsl.IsAllowedFolder(info) {
	//		if filepath.IsAbs(arg) {
	//			return m.folderAbs.AbsorbFolderFromDisk(m.fsl, arg)
	//		}
	//		return m.folderRel.AbsorbFolderFromDisk(m.fsl, arg)
	//	}
	//	return fmt.Errorf("illegal folder %q", info.Name())
	//}
	//if m.fsl.IsAllowedFile(info) {
	//	if filepath.IsAbs(arg) {
	//		return m.folderAbs.AbsorbFileFromDisk(arg)
	//	}
	//	return m.folderRel.AbsorbFileFromDisk(arg)
	//}
	//return fmt.Errorf("not a markdown file %q", info.Name())

	//	return nil, nil
}

// MyContrivedFolder is a named grouping of some subset of local disk,
// gathered from both absolute paths and paths specified relative to the
// CWD of the current process, and a list of GitHub repos.
type MyContrivedFolder struct {
	name      string
	fsl       *FsLoader
	repos     []*MyRepo
	folderAbs *MyFolder
	cwd       string
	// folderRel is relative to cwd
	folderRel *MyFolder
}

var _ MyTreeItem = &MyContrivedFolder{}

func (m *MyContrivedFolder) Accept(v TreeVisitor) {
	v.VisitContrivedFolder(m)
}

func (m *MyContrivedFolder) Parent() MyTreeItem {
	return nil
}

func (m *MyContrivedFolder) Name() string {
	return m.name
}

func (m *MyContrivedFolder) FullName() string {
	// TODO: use originalSpecs?
	return m.name
}

func (m *MyContrivedFolder) DirName() string {
	return ""
}

// Initialize processes the given arguments.
// If no error is returned, all the associated arguments are
// available on disk and readable when the func returns.
func (m *MyContrivedFolder) Initialize(
	arg string, ts *FsLoader) error {
	if arg == "" {
		return fmt.Errorf("need an arg")
	}
	if ts == nil {
		ts = DefaultFsLoader
	}
	m.fsl = ts
	m.name = "contrived" // TODO: something better?
	{
		tmp, err := os.Getwd()
		if err != nil {
			return err
		}
		m.cwd = stripTrailingSlash(tmp)
	}
	m.folderAbs = &MyFolder{myTreeItem: myTreeItem{name: "/"}}
	m.folderRel = &MyFolder{myTreeItem: myTreeItem{name: m.cwd}}
	if smellsLikeGithubCloneArg(arg) {
		return m.absorbRepo(arg)
	}
	info, err := os.Stat(arg)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if m.fsl.IsAllowedFolder(info) {
			if filepath.IsAbs(arg) {
				return m.folderAbs.AbsorbFolderFromDisk(m.fsl, arg)
			}
			return m.folderRel.AbsorbFolderFromDisk(m.fsl, arg)
		}
		return fmt.Errorf("illegal folder %q", info.Name())
	}
	if m.fsl.IsAllowedFile(info) {
		if filepath.IsAbs(arg) {
			return m.folderAbs.AbsorbFileFromDisk(arg)
		}
		return m.folderRel.AbsorbFileFromDisk(arg)
	}
	return fmt.Errorf("not a markdown file %q", info.Name())
}

func (m *MyContrivedFolder) absorbRepo(arg string) error {
	n, p, err := extractGithubRepoName(arg)
	if err != nil {
		return err
	}
	for _, r := range m.repos {
		if r.Name() == n {
			return fmt.Errorf("already loaded %s", n)
		}
	}
	r := &MyRepo{
		name: n,
		path: p,
	}
	if err = r.Init(m.fsl); err != nil {
		return err
	}
	m.repos = append(m.repos, r)
	return nil
}

// Cleanup cleans up temp space.
func (m *MyContrivedFolder) Cleanup() {
	for _, r := range m.repos {
		r.CleanUp()
	}
}
