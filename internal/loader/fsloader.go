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
	ReadmeFileName   = "README.md"
	OrderingFileName = "README_ORDER.txt"
	rootSlash        = string(filepath.Separator)
	currentDir       = "."
	selfPath         = currentDir + rootSlash
	upDir            = ".."
)

// LoadTree loads a file tree from disk, possibly after first cloning a repo from GitHub.
func (fsl *FsLoader) LoadTree(rawPath string) (*MyFolder, error) {
	if smellsLikeGithubCloneArg(rawPath) {
		return CloneAndLoadRepo(fsl, rawPath)
	}
	return fsl.LoadFolder(rawPath)
}

// LoadFolder loads the files at or below a path into memory, returning them
// inside an MyFolder instance.
//
// Files or folders that don't pass the provided filters are excluded.
// If filtering leaves a folder empty, it is discarded.  If nothing
// makes it through, the function returns a nil folder and no error.
//
// If an "OrderingFileName" is found in a directory, it's used to sort the files
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
	// If rawPath is empty, cleanPath ends up with "."
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

	if info.IsDir() {
		if err = fsl.IsAllowedFolder(info); err != nil {
			// If user explicitly asked for a disallowed folder, complain.
			// Deeper in, when absorbing folders, they are simply ignored.
			err = fmt.Errorf("illegal folder %q; %w", info.Name(), err)
			return
		}
		fld, err = fsl.loadFolder(cleanPath)
		if err != nil {
			return
		}
		if !fld.IsEmpty() {
			fld.name = cleanPath
			return
		}
		return nil, nil
	}
	// Load just one file.
	if err = fsl.IsAllowedFile(info); err != nil {
		// If user explicitly asked for a disallowed file, complain.
		// Deeper in, when absorbing folders, they are simply ignored.
		err = fmt.Errorf("illegal file %q; %w", info.Name(), err)
		return
	}
	dir, base := DirBase(cleanPath)
	var c []byte
	c, err = fsl.fs.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}
	fld = NewFolder(dir).AddFileObject(NewFile(base, c))
	return
}

// loadFolder loads the folder specified by the path.
// This is the recursive part of the LoadFolder entrypoint.
// The path must point to a folder.
// For example, given a file system like
//
//	/home/bob/
//	  f1.md
//	  games/
//	    doom.md
//
// The argument /home/bob should yield an unnamed, unparented folder containing
// 'f1.md' and the folder 'game' (with 'doom.md' inside 'game').
//
// The same thing is returned if the file system is
//
//	./
//	  f1.md
//	  games/
//	    doom.md
//
// and the argument passed in is simply "." or an empty string.
func (fsl *FsLoader) loadFolder(path string) (*MyFolder, error) {
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
		subPath := filepath.Join(path, info.Name())
		if info.IsDir() {
			if err = fsl.IsAllowedFolder(info); err == nil {
				if subFld, err = fsl.loadFolder(subPath); err != nil {
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
			if ordering, err = LoadOrderFile(fsl.fs, subPath); err != nil {
				return nil, err
			}
			continue
		}
		if err = fsl.IsAllowedFile(info); err == nil {
			fi := NewEmptyFile(info.Name())
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
