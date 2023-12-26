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

const (
	orderingFile = "README_ORDER.txt"
	rootSlash    = string(filepath.Separator)
	currentDir   = "."
	selfPath     = currentDir + rootSlash
	upDir        = ".."
)

// LoadFolder loads the files at or below a path into memory, returning them
// inside an MyFolder instance.
//
// Only "allowed" files and folders are included, per the filter functions
// provided. If filtering leaves a folder empty, it is discarded.  If nothing
// makes it through, the function returns a nil folder and no error.
//
// If an "orderingFile" is found in a directory, it's used to sort the files
// and sub-folders in that folder's in-memory representation. An ordering file
// is just text, with one name per line. Ordered files appear first, with the
// remainder in the order imposed by fs.ReadDir.
//
// If the path is a file, only that file is loaded.  Since LoadFolder must
// return a folder, the folder's name is the path to that file minus the file's
// name.  The path might be absolute (starting with the rootSlash) or relative,
// starting with any legal character other than rootSlash. If the path is only
// a file name (no folder names), then the name of the returned folder is ".",
// and it contains only one file.
//
// Examples:
//
//	             path | returned folder name | contents
//	------------------+----------------------+--------------
//	           foo.md |                    . | foo.md
//	         ./foo.md |                    . | foo.md
//	        ../foo.md |            {illegal} | {illegal}
//	/usr/local/foo.md |           /usr/local | foo.md
//	       bar/foo.md |                  bar | foo.md
//
// If the path is a folder, only that folder is loaded; folders rooted higher
// in the tree are ignored.
//
// Examples:
//
//	             path | returned folder name | contents
//	------------------+----------------------+--------------
//	              foo |                  foo | {whatever}
//	            ./foo |                  foo | {whatever}
//	           ../foo |            {illegal} | {illegal}
//	   /usr/local/foo |       /usr/local/foo | {whatever}
//	          bar/foo |              bar/foo | {whatever}
//
// Any error returned will be from the file system.
func (fsl *FsLoader) LoadFolder(rawPath string) (fld *MyFolder, err error) {
	cleanPath := filepath.Clean(rawPath)
	// For now, disallow paths that start with upDir, because in the task at
	// hand we want a clear root directory for display.
	if strings.HasPrefix(cleanPath, upDir) {
		return nil, fmt.Errorf(
			"specify absolute path or something at or below your working directory")
	}
	var info os.FileInfo
	info, err = fsl.fs.Stat(cleanPath)
	if err != nil {
		return
	}
	dir, base := filepath.Dir(cleanPath), filepath.Base(cleanPath)
	// Behavior:
	//	             path  |       dir  | base
	//	-------------------+------------+-----------
	//	   {empty string}  |         .  |  .
	//	                .  |         .  |  .
	//	               ./  |         .  |  .
	//                  /  |         /  |  /
	//	            ./foo  |         .  | foo
	//	           ../foo  |        ..  | foo
	//	             /foo  |         /  | foo
	//	   /usr/local/foo  | /usr/local | foo

	fld = &MyFolder{myTreeItem: myTreeItem{name: dir}}

	if !info.IsDir() {
		if err = fsl.IsAllowedFile(info); err != nil {
			// If user explicitly asked for a disallowed file, complain.
			// Deeper in, when absorbing folders, they are just ignored.
			err = fmt.Errorf("illegal file %q; %w", info.Name(), err)
			return
		}
		err = fsl.loadSubFile(fld, base)
		return
	}
	if err = fsl.IsAllowedFolder(info); err != nil {
		// If user explicitly asked for a disallowed folder, complain.
		// Deeper in, when absorbing folders, they are just ignored.
		err = fmt.Errorf("illegal folder %q; %w", info.Name(), err)
		return
	}
	if base == "" {
		// Special case - user asking for the entire root file system.
		if dir != rootSlash {
			err = fmt.Errorf("something wrong with dir/base split")
			return
		}
		fld, err = fsl.loadSubFolder(fld, base)
		if fld != nil {
			fld.name = rootSlash
		}
		return
	}
	_, err = fsl.loadSubFolder(fld, base)
	if fld.IsEmpty() {
		fld = nil
	}
	return
}

// loadSubFolder returns a fully loaded MyFolder instance.
// The arguments are a parent folder, and the simple name of a folder inside
// the parent - no path separators in the folderName.
// The parent's name must be either a full absolute path or a relative
// path that makes sense with respect to the process' working directory.
// This function returns a new folder object that knows about its parent, and
// the parent is informed of the child.
func (fsl *FsLoader) loadSubFolder(
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
				if _, err = fsl.loadSubFolder(&fld, info.Name()); err != nil {
					return nil, err
				}
			}
			continue
		}
		if IsOrderingFile(info) {
			ordering, err = LoadOrderFile(info)
			if err != nil {
				return nil, err
			}
			continue
		}
		if err = fsl.IsAllowedFile(info); err == nil {
			if err = fsl.loadSubFile(&fld, info.Name()); err != nil {
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
	parent.dirs = append(parent.dirs, &fld)
	return &fld, nil
}

// loadSubFile loads a file into a folder.
func (fsl *FsLoader) loadSubFile(fld *MyFolder, n string) (err error) {
	fi := NewEmptyFile(n)
	fld.AddFileObject(fi)
	fi.content, err = fsl.fs.ReadFile(fi.FullName())
	return
}
