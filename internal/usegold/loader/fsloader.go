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
	dir, name := FSplit(cleanPath)
	fld = &MyFolder{myTreeItem: myTreeItem{name: dir}}
	if info.IsDir() {
		err = fsl.IsAllowedFolder(info)
		if err != nil {
			err = fmt.Errorf("illegal folder %q; %w", info.Name(), err)
			return
		}
		var sub *MyFolder
		sub, err = fsl.LoadSubFolder(fld, name)
		if err != nil {
			return
		}
		fld.dirs = append(fld.dirs, sub)
		return
	}
	err = fsl.IsAllowedFile(info)
	if err != nil {
		err = fmt.Errorf("illegal file %q; %w", info.Name(), err)
		return
	}
	fld.files = append(fld.files, NewFile(name))
	return
}

// LoadSubFolder returns a MyFolder instance representing an FS folder.
// The arguments are a parent folder, and the simple name of a folder inside
// the parent - no path separators in the folderName.
// The parent's name must be either a full absolute path or a relative
// path that makes sense with respect to the process's working directory.
// This function returns a new folder object, loaded with all approved
// sub-folders and their files. The new folder knows about its parent, but
// the parent doesn't know about it. The function returns nil if the folder
// is empty or has no approved sub-folders or files.
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
				var child *MyFolder
				child, err = fsl.LoadSubFolder(&fld, info.Name())
				if err != nil {
					return nil, err
				}
				if child != nil {
					fld.dirs = append(fld.dirs, child)
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
			if err = fld.loadFileFromFs(info.Name()); err != nil {
				return nil, err
			}
		}
	}
	if fld.IsEmpty() {
		slog.Debug("omitting empty directory", "dir", fld.FullName())
		return nil, nil
	}
	fld.files = ReorderFiles(fld.files, ordering)
	fld.dirs = ReorderFolders(fld.dirs, ordering)
	return &fld, nil
}
