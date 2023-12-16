package loader

import (
	"fmt"
	"github.com/monopole/mdrip/base"
	"os"
	"strings"
)

// dataNode holds either markdown data or pointers to markdown data.
type dataNode interface {
	// Origin is some description of the origin of the data node.
	Origin() string
	// Name is the name of the data node.
	Name() string
	// Children is a set of child dataNode.
	// If non-empty, then Content must be nil.
	Children() []dataNode
	// Content is the raw markdown content.
	// If non-empty, then Children must be nil.
	Content() []byte
}

var _ dataNode = &File{}

func (fi *File) Children() []dataNode {
	return nil
}

func (fi *File) Content() []byte {
	return fi.content
}

func (fi *File) Name() string {
	return fi.path.Base()
}

func (fi *File) Origin() string {
	return fi.Name()
}

// File is named byte array.
type File struct {
	path    base.FilePath
	content []byte
}

// Folder is a named grouping of files and folders.
type Folder struct {
	path base.FilePath
	baseFolder
}

// GitFolder holds date cloned from a git server.
type GitFolder struct {
	repoName string
	Folder
}

// baseFolder is a nameless grouping of files and folders,
// likely collected together from command line arguments.
type baseFolder struct {
	files []*File
	dirs  []*Folder
}

// ContrivedFolder is a named grouping of files and folders
// that doesn't correspond to a "real" Folder.  Likely it's
// collected together from command line arguments.
type ContrivedFolder struct {
	name  string
	items []string
	baseFolder
}

var _ dataNode = &ContrivedFolder{}

func (bf *baseFolder) Content() []byte {
	return nil
}

func (bf *baseFolder) Children() (result []dataNode) {
	result = make([]dataNode, len(bf.dirs)+len(bf.files))
	for i := range bf.files {
		result[i] = bf.files[i]
	}
	for i := range bf.dirs {
		result[i+len(bf.files)] = bf.dirs[i]
	}
	return
}

func (bf *baseFolder) IsEmpty() bool {
	return len(bf.dirs) == 0 && len(bf.files) == 0
}

func (cf *ContrivedFolder) Name() string {
	return cf.name
}

func (cf *ContrivedFolder) Origin() string {
	return strings.Join(cf.items, ",")
}

var _ dataNode = &Folder{}

func (fl *Folder) Name() string {
	return fl.path.Base()
}

func (fl *Folder) Origin() string {
	return string(fl.path)
}

func Load(args []string) (dataNode, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("needs some args")
	}
	if len(args) > 0 {
		return loadMany(args)
	}
	return loadOne(args[0])
}

func loadMany(args []string) (*ContrivedFolder, error) {
	var result ContrivedFolder
	for _, arg := range args {
		item, err := loadOne(arg)
		if err != nil {
			return nil, err
		}
		if item == nil {
			continue
		}
		if f, isFile := item.(*File); isFile {
			result.files = append(result.files, f)
		} else {
			d, isDir := item.(*Folder)
			if !isDir {
				panic(fmt.Sprintf("%s isn't a folder or a file", arg))
			}
			result.dirs = append(result.dirs, d)
		}
	}
	if result.IsEmpty() {
		return nil, nil
	}
	return &result, nil
}

func loadOne(arg string) (dataNode, error) {
	if smellsLikeGithubCloneArg(arg) {
		repoName, path, err := extractGithubRepoName(arg)
		if err != nil {
			return nil, err
		}
		return loadFromGithub(repoName, path)
	}
	return loadFromPath(base.FilePath(arg))
}

// errRaw returns "fake" markdown file showing an error message.
func errRaw(arg string, err error) (*File, error) {
	return errFile(base.FilePath(arg), err)
}

// errFile returns "fake" markdown file showing an error message.
func errFile(fp base.FilePath, err error) (*File, error) {
	return &File{
		path:    fp,
		content: []byte(fmt.Sprintf("## Unable to load from %s; %s", fp, err.Error())),
	}, err
}

func loadFromPath(path base.FilePath) (dataNode, error) {
	if path.IsDesirableFile() {
		return scanFile(path)
	}
	if path.IsDesirableDir() {
		return scanDir(path)
	}
	return errFile(path, fmt.Errorf("%s isn't desirable", path))
}

func scanFile(path base.FilePath) (*File, error) {
	contents, err := os.ReadFile(string(path))
	if err != nil {
		return errFile(path, fmt.Errorf("file read error (%w)", err))
	}
	return &File{
		path:    path,
		content: contents,
	}, nil
}

// scanDir assumes that the given directory is desirable.
func scanDir(path base.FilePath) (*Folder, error) {
	dirEntries, err := os.ReadDir(string(path))
	if err != nil {
		return nil, fmt.Errorf("unable to read folder %q; %w", path, err)
	}
	var (
		ordering []string
		files    []*File
		dirs     []*Folder
	)
	for _, f := range dirEntries {
		p := path.Join(f)
		if p.IsDesirableFile() {
			item, _ := scanFile(p)
			files = append(files, item)
			continue
		}
		if p.IsDesirableDir() {
			if item, er := scanDir(p); er == nil && item != nil {
				dirs = append(dirs, item)
			}
			continue
		}
		if p.IsOrderFile() {
			if contents, er := p.Read(); er == nil {
				ordering = strings.Split(contents, "\n")
			}
		}
	}
	if len(dirs) == 0 && len(files) == 0 {
		return nil, nil
	}
	return &Folder{
		path: path,
		baseFolder: baseFolder{
			files: reorderFiles(files, ordering),
			dirs:  reorderFolders(dirs, ordering),
		},
	}, nil
}

func reorderFolders(x []*Folder, ordering []string) []*Folder {
	for i := len(ordering) - 1; i >= 0; i-- {
		x = shiftFolderToTop(x, ordering[i])
	}
	return x
}

func shiftFolderToTop(x []*Folder, top string) []*Folder {
	var first []*Folder
	var remainder []*Folder
	for _, f := range x {
		if f.Name() == top {
			first = append(first, f)
		} else {
			remainder = append(remainder, f)
		}
	}
	return append(first, remainder...)
}

func reorderFiles(x []*File, ordering []string) []*File {
	for i := len(ordering) - 1; i >= 0; i-- {
		x = shiftFileToTop(x, ordering[i])
	}
	return shiftFileToTop(x, "README")
}

func shiftFileToTop(x []*File, top string) []*File {
	var first []*File
	var remainder []*File
	for _, f := range x {
		if f.Name() == top {
			first = append(first, f)
		} else {
			remainder = append(remainder, f)
		}
	}
	return append(first, remainder...)
}
