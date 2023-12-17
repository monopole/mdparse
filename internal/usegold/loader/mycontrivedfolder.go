package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

type filter func(info os.FileInfo) bool

// MyContrivedFolder is a named grouping of some subset of local disk,
// gather from both absolute paths and paths specified relative to the
// CWD of the current process, and a list of repo specifications.
type MyContrivedFolder struct {
	name            string
	originalSpecs   []string
	repos           []*MyRepo
	folderAbs       *MyFolder
	cwd             string
	folderRel       *MyFolder
	isAllowedFile   filter
	isAllowedFolder filter
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
// available on disk and readable when the func return.
func (m *MyContrivedFolder) Initialize(
	args []string, isAllowedFile, isAllowedFolder filter) error {
	if len(args) == 0 {
		return fmt.Errorf("needs some args")
	}
	m.isAllowedFile = isAllowedFile
	m.isAllowedFolder = isAllowedFolder
	m.name = "contrived" // TODO: something better?
	{
		tmp, err := os.Getwd()
		if err != nil {
			return err
		}
		m.cwd = stripTrailingSlash(tmp)
	}
	m.folderAbs = &MyFolder{
		myTreeItem: myTreeItem{
			name: "/",
		},
	}
	m.folderRel = &MyFolder{
		myTreeItem: myTreeItem{
			name: m.cwd,
		},
	}
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
		if m.isAllowedFolder(info) {
			if filepath.IsAbs(arg) {
				return m.folderAbs.absorbFolder(arg)
			}
			return m.folderRel.absorbFolder(arg)
		}
		return fmt.Errorf("illegal folder %q", info.Name())
	}
	if m.isAllowedFile(info) {
		if filepath.IsAbs(arg) {
			return m.folderAbs.absorbFile(arg)
		}
		return m.folderRel.absorbFile(arg)
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
	if err = r.Init(); err != nil {
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
