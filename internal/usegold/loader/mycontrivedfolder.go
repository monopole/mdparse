package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

// MyContrivedFolder is a named grouping of some subset of local disk,
// gather from both absolute paths and paths specified relative to the
// CWD of the current process, and a list of repo specifications.
type MyContrivedFolder struct {
	name          string
	originalSpecs []string
	repos         []*MyRepo
	folderAbs     *MyFolder
	cwd           string
	folderRel     *MyFolder
}

func (m *MyContrivedFolder) Dump() {
	fmt.Printf("%s, nRepos=%d, %v\n", m.name, len(m.repos), m.originalSpecs)
	for i := range m.repos {
		m.repos[i].Dump()
	}
	m.folderAbs.Dump(0)
	m.folderRel.Dump(0)
}

func (m *MyContrivedFolder) Parent() MyTreeItem {
	return nil
}

func (m *MyContrivedFolder) Name() string {
	return m.name
}

func (m *MyContrivedFolder) FullName() string {
	return m.name
}

func (m *MyContrivedFolder) DirName() string {
	return ""
}

var _ MyTreeItem = &MyContrivedFolder{}

// NewMyContrivedFolder returns an instance with validated arguments.
// If this returns without error, all the associated arguments are
// available on disk and readable (at the moment).
func NewMyContrivedFolder(args []string) (*MyContrivedFolder, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("needs some args")
	}
	var f MyContrivedFolder
	f.name = "contrived" // TODO: something better?
	f.folderAbs = &MyFolder{
		myTreeItem: myTreeItem{
			name: "{ABS}",
		},
	}
	f.folderRel = &MyFolder{
		myTreeItem: myTreeItem{
			name: "{REL}",
		},
	}
	f.cwd, _ = os.Getwd()
	f.originalSpecs = make([]string, len(args))
	for i := range args {
		if err := f.absorb(args[i]); err != nil {
			return nil, err
		}
		f.originalSpecs[i] = args[i]
	}
	return &f, nil
}

func (m *MyContrivedFolder) absorb(arg string) error {
	if smellsLikeGithubCloneArg(arg) {
		repoName, path, err := extractGithubRepoName(arg)
		if err != nil {
			return err
		}
		return m.absorbRepo(repoName, path)
	}
	info, err := os.Stat(arg)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if isAnAllowedFolder(info) {
			if filepath.IsAbs(arg) {
				return m.folderAbs.absorbFolder(arg)
			}
			return m.folderRel.absorbFolder(arg)
		}
		return fmt.Errorf("illegal folder %q", info.Name())
	}
	if isAnAllowedFile(info) {
		if filepath.IsAbs(arg) {
			return m.folderAbs.absorbFile(arg)
		}
		return m.folderRel.absorbFile(arg)
	}
	return fmt.Errorf("not a markdown file %q", info.Name())
}

func (m *MyContrivedFolder) absorbRepo(repoName, path string) error {
	fmt.Println("*** pretending to absorb repo")
	return nil
}

// Reload returns all the content.
func (m *MyContrivedFolder) Reload() (interface{}, error) {
	panic("no impl")
	return nil, nil
}
