package loader

import (
	"os"
	"path/filepath"
	"strings"
)

// MyIsOrderFile returns true if the file appears to be a "reorder"
// file specifying how to re-order the files in the directory
// in some fashion other than directory order.
func MyIsOrderFile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}
	if !info.Mode().IsRegular() {
		return false
	}
	return filepath.Base(info.Name()) == "README_ORDER.txt"
}

// LoadOrderFile returns a list of names specify file name order priority.
func LoadOrderFile(info os.FileInfo) ([]string, error) {
	contents, err := os.ReadFile(info.Name())
	if err != nil {
		return nil, err
	}
	return strings.Split(string(contents), "\n"), nil
}

func ReorderFolders(x []*MyFolder, ordering []string) []*MyFolder {
	for i := len(ordering) - 1; i >= 0; i-- {
		x = ShiftFolderToTop(x, ordering[i])
	}
	return x
}

func ShiftFolderToTop(x []*MyFolder, top string) []*MyFolder {
	var first []*MyFolder
	var remainder []*MyFolder
	for _, f := range x {
		if f.Name() == top {
			first = append(first, f)
		} else {
			remainder = append(remainder, f)
		}
	}
	return append(first, remainder...)
}

func ReorderFiles(x []*MyFile, ordering []string) []*MyFile {
	for i := len(ordering) - 1; i >= 0; i-- {
		x = ShiftFileToTop(x, ordering[i])
	}
	return ShiftFileToTop(x, "README")
}

func ShiftFileToTop(x []*MyFile, top string) []*MyFile {
	var first []*MyFile
	var remainder []*MyFile
	for _, f := range x {
		if f.Name() == top {
			first = append(first, f)
		} else {
			remainder = append(remainder, f)
		}
	}
	return append(first, remainder...)
}
