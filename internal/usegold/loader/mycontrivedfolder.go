package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

type filter func(info os.FileInfo) bool

// MyContrivedFolder is a named grouping of some subset of local disk,
// gathered from both absolute paths and paths specified relative to the
// CWD of the current process, and a list of GitHub repos.
type MyContrivedFolder struct {
	name          string
	originalSpecs []string
	fsl           *FsLoader
	repos         []*MyRepo
	folderAbs     *MyFolder
	cwd           string
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
	args []string, ts *FsLoader) error {
	if len(args) == 0 {
		return fmt.Errorf("needs some args")
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
	m.originalSpecs = make([]string, len(args))
	for i := range args {
		if err := m.absorb(args[i]); err != nil {
			return err
		}
		m.originalSpecs[i] = args[i]
	}
	return nil
}

func (m *MyContrivedFolder) absorb(arg string) error {
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
