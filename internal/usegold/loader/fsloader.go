package loader

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// FsLoader navigates and reads a file system.
// It would be a good place to hold a mockable file system pointer,
// to remove the dependence on the "os" package.
type FsLoader struct {
	IsAllowedFile, IsAllowedFolder filter
}

// DefaultFsLoader is an FsLoader instance configured with
// default file and folder filters.
var DefaultFsLoader = &FsLoader{
	IsAllowedFile:   IsAnAllowedFile,
	IsAllowedFolder: IsAnAllowedFolder,
}

// loadFolderFromFs accepts a folderName and a fully formed parent folder.
// The parent must be fully formed as this allows the folderName to be located.
// It assumes that the folderName is approved/desirable.  It returns a
// new folder object, loaded with all approved sub-folders and their files.
// It returns nil if the folder is empty or has no approved sub-folders or files.
func (fsl *FsLoader) loadFolderFromFs(parent *MyFolder, folderName string) (*MyFolder, error) {
	n := filepath.Join(parent.FullName(), folderName)
	dirEntries, err := os.ReadDir(n)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read folder %q; %w", n, err)
	}
	var (
		fld      MyFolder
		ordering []string
	)
	fld.parent = parent
	fld.name = folderName
	for i := range dirEntries {
		var info os.FileInfo
		info, err = dirEntries[i].Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			if fsl.IsAllowedFolder(info) {
				var child *MyFolder
				child, err = fsl.loadFolderFromFs(&fld, info.Name())
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
		if fsl.IsAllowedFile(info) {
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
