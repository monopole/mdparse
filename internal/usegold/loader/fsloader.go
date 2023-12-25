package loader

import (
	"fmt"
	"github.com/spf13/afero"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// FsLoader navigates and reads a file system.
type FsLoader struct {
	IsAllowedFile, IsAllowedFolder filter
	fs                             *afero.Afero
}

// NewFsLoader returns a file system (FS) loader with default filters.
// For an in-memory FS, inject afero.NewMemMapFs().
// For a "real" disk-based system, inject afero.NewOsFs().
func NewFsLoader(fs afero.Fs) *FsLoader {
	return &FsLoader{
		IsAllowedFile:   IsMarkDownFile,
		IsAllowedFolder: IsNotADotDir,
		fs:              &afero.Afero{Fs: fs},
	}
}

// Read returns the contents of a file.
func (fsl *FsLoader) Read(fi *MyFile) ([]byte, error) {
	return fsl.fs.ReadFile(fi.FullName())
	//// myErrFile returns "fake" markdown file showing an error message.
	//func myErrFile(path string, err error) (*MyFile, error) {
	//	return &MyFile{
	//		myTreeItem: myTreeItem{
	//			parent: nil,
	//			name:   path,
	//		},
	//		content: []byte(fmt.Sprintf("## Unable to load from %s; %s", path, err.Error())),
	//	}, err
	//}
}

func (fsl *FsLoader) LoadFolder(path string) (fld *MyFolder, err error) {
	cleanPath := filepath.Clean(path)
	// Disallow paths that start with "..", because in the task at hand
	// we want a clear root directory for display.
	if strings.HasPrefix(cleanPath, "..") {
		return nil, fmt.Errorf(
			"specify absolute path or something at or below your working directory")
	}
	var info os.FileInfo
	info, err = fsl.fs.Stat(cleanPath)
	if err != nil {
		return
	}
	dir, base := FSplit(cleanPath)
	if dir == "" {
		dir = "."
	}
	fld = &MyFolder{myTreeItem: myTreeItem{name: dir}}
	if info.IsDir() {
		if err = fsl.IsAllowedFolder(info); err != nil {
			err = fmt.Errorf("illegal folder %q; %w", info.Name(), err)
			return
		}
		_, err = fsl.LoadSubFolder(fld, base)
		return
	}
	if err = fsl.IsAllowedFile(info); err != nil {
		err = fmt.Errorf("illegal file %q; %w", info.Name(), err)
		return
	}
	fld.AddFileObject(NewFile(base))
	return
}

// LoadSubFolder returns a MyFolder instance representing an FS folder.
// The arguments are a parent folder, and the simple name of a folder inside
// the parent - no path separators in the folderName.
// The parent's name must be either a full absolute path or a relative
// path that makes sense with respect to the process' working directory.
// This function returns a new folder object, loaded with all approved
// sub-folders and their files. The new folder knows about its parent, and
// the parent is informed of the child.
// If "order" files are encountered in a given sub-folder, they are obeyed
// to sort the files and sub-folders at a given level.
// Any error returned will be from the file system.
func (fsl *FsLoader) LoadSubFolder(
	parent *MyFolder, folderName string) (*MyFolder, error) {
	fullName := filepath.Join(parent.FullName(), folderName)
	dirEntries, err := fsl.fs.ReadDir(fullName)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read folder %q; %w", fullName, err)
	}
	var (
		fld      MyFolder
		ordering []string
	)
	fld.parent = parent
	fld.name = folderName
	for i := range dirEntries {
		info := dirEntries[i]
		if info.IsDir() {
			if err = fsl.IsAllowedFolder(info); err == nil {
				if _, err = fsl.LoadSubFolder(&fld, info.Name()); err != nil {
					return nil, err
				}
			}
			continue
		}
		if MyIsOrderFile(info) {
			ordering, err = LoadOrderFile(info)
			if err != nil {
				return nil, err
			}
			continue
		}
		if err = fsl.IsAllowedFile(info); err == nil {
			fld.AddFileObject(NewFile(info.Name()))
			// TODO: Read the file contents
		}
	}
	if fld.IsEmpty() {
		slog.Debug("omitting empty directory", "dir", fld.FullName())
		return nil, nil
	}
	fld.files = ReorderFiles(fld.files, ordering)
	fld.dirs = ReorderFolders(fld.dirs, ordering)
	parent.dirs = append(parent.dirs, &fld)
	return &fld, nil
}
