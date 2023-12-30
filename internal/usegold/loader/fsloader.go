package loader

import (
	"fmt"
	"github.com/spf13/afero"
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
// Files or folders that don't pass the provided filters are excluded.
// If filtering leaves a folder empty, it is discarded.  If nothing
// makes it through, the function returns a nil folder and no error.
//
// If an "orderingFile" is found in a directory, it's used to sort the files
// and sub-folders in that folder's in-memory representation. An ordering file
// is just lines of text, one name per line. Ordered files appear first, with
// the remainder in the order imposed by fs.ReadDir.
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
//	   {empty string} |                    . | {whatever}
//	                . |                    . | {whatever}
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
	// hand we want a clear root directory for display. Might allow this later.
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

	if !info.IsDir() {
		if err = fsl.IsAllowedFile(info); err != nil {
			// If user explicitly asked for a disallowed file, complain.
			// Deeper in, when absorbing folders, they are simply ignored.
			err = fmt.Errorf("illegal file %q; %w", info.Name(), err)
			return
		}
		if base != info.Name() {
			panic("assumption 1 about filepath.Base vs filepath.Dir broken")
		}
		fi := NewEmptyFile(base)
		fi.content, err = fsl.fs.ReadFile(cleanPath)
		if err != nil {
			return nil, err
		}
		fld = &MyFolder{myTreeItem: myTreeItem{name: dir}}
		fld.AddFileObject(fi)
		return
	}
	if err = fsl.IsAllowedFolder(info); err != nil {
		// If user explicitly asked for a disallowed folder, complain.
		// Deeper in, when absorbing folders, they are simply ignored.
		err = fmt.Errorf("illegal folder %q; %w", info.Name(), err)
		return
	}
	if base == rootSlash || base == currentDir {
		if dir != base {
			panic("assumption 2 about filepath.Base vs filepath.Dir broken")
		}
		fld, err = fsl.loadPath(cleanPath)
		if err != nil {
			return
		}
		if fld != nil {
			fld.name = base
		}
		return
	}
	fld, err = fsl.loadPath(cleanPath)
	if err != nil {
		return
	}
	if fld != nil {
		fld.name = cleanPath
	}
	if fld.IsEmpty() {
		fld = nil
	}
	return
}

// loadPath loads from the path, returning an instance of MyFolder.
func (fsl *FsLoader) loadPath(path string) (*MyFolder, error) {
	var (
		result   MyFolder
		subFld   *MyFolder
		ordering []string
	)
	dirEntries, err := fsl.fs.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read folder %q; %w", path, err)
	}
	for i := range dirEntries {
		info := dirEntries[i]
		if info.IsDir() {
			if err = fsl.IsAllowedFolder(info); err == nil {
				subPath := filepath.Join(path, info.Name())
				if subFld, err = fsl.loadPath(subPath); err != nil {
					return nil, err
				}
				if !subFld.IsEmpty() {
					subFld.name = info.Name()
					result.AddFolderObject(subFld)
				}
			}
			continue
		}
		if IsOrderingFile(info) {
			// load it and keep it for use at end of function.
			// TODO: test this.  i don't think the path is correct
			if ordering, err = LoadOrderFile(info); err != nil {
				return nil, err
			}
			continue
		}
		if err = fsl.IsAllowedFile(info); err == nil {
			fi := NewEmptyFile(info.Name())
			subPath := filepath.Join(path, info.Name())
			fi.content, err = fsl.fs.ReadFile(subPath)
			if err != nil {
				return nil, err
			}
			result.AddFileObject(fi)
		}
	}
	if result.IsEmpty() {
		return nil, nil
	}
	result.files = ReorderFiles(result.files, ordering)
	result.dirs = ReorderFolders(result.dirs, ordering)
	return &result, nil
}

//// loadSubFile loads a file into a folder.
//func (fsl *FsLoader) loadSubFile(fld *MyFolder, n string) (err error) {
//	fi := NewEmptyFile(n)
//	fld.AddFileObject(fi)
//	fi.content, err = fsl.fs.ReadFile(fi.FullName())
//	return
//}
