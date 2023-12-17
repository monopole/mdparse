package loader

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type TreeScanner struct {
	IsAllowedFile, IsAllowedFolder filter
}

var DefaultTreeScanner = &TreeScanner{
	IsAllowedFile:   IsAnAllowedFile,
	IsAllowedFolder: IsAnAllowedFolder,
}

func (ts *TreeScanner) addScannedFolder(parent *MyFolder, name string) error {
	slog.Debug("adding   FOLDER", "name", name, "parent", parent.FullName())
	child, err := ts.readFolderContentsFromDisk(parent, name)
	if err != nil {
		return err
	}
	if child != nil {
		parent.dirs = append(parent.dirs, child)
	}
	return nil
}

func (ts *TreeScanner) readFolderContentsFromDisk(parent *MyFolder, name string) (*MyFolder, error) {
	n := filepath.Join(parent.FullName(), name)
	dirEntries, err := os.ReadDir(n)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read folder %q; %w", n, err)
	}
	var (
		child    MyFolder
		ordering []string
	)
	child.parent = parent
	child.name = name
	for i := range dirEntries {
		var info os.FileInfo
		info, err = dirEntries[i].Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			if ts.IsAllowedFolder(info) {
				err = ts.addScannedFolder(&child, info.Name())
				if err != nil {
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
		if ts.IsAllowedFile(info) {
			err = child.addFile(info.Name())
			if err != nil {
				return nil, err
			}
		}
	}
	if child.IsEmpty() {
		slog.Debug("omitting empty directory", "dir", child.FullName())
		return nil, nil
	}
	child.files = ReorderFiles(child.files, ordering)
	child.dirs = ReorderFolders(child.dirs, ordering)
	return &child, nil
}
