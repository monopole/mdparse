package loader

import (
	"os"
	"path/filepath"
	"strings"
)

const blanks = "                                                                "

func fSplit(path string) (dir, name string) {
	dir, name = filepath.Split(path)
	if strings.HasSuffix(dir, string(filepath.Separator)) {
		dir = dir[:len(dir)-1]
	}
	return
}

func isAnAllowedFile(info os.FileInfo) bool {
	if !info.Mode().IsRegular() {
		return false
	}
	if filepath.Ext(info.Name()) != ".md" {
		return false
	}
	base := filepath.Base(info.Name())
	const badLeadingChar = "~.#"
	return strings.Index(badLeadingChar, string(base[0])) < 0
}

func isAnAllowedFolder(info os.FileInfo) bool {
	n := info.Name()
	// Allow special dir names.
	if n == "." || n == "./" || n == ".." {
		return true
	}
	// Ignore .git, etc.
	return !strings.HasPrefix(filepath.Base(n), ".")
}

// MyIsOrderFile returns true if the file appears to be a "reorder"
// file specifying how to re-order the files in the directory
// in some fashion other than directory order.
func MyIsOrderFile(n string) bool {
	s, err := os.Stat(n)
	if err != nil {
		return false
	}
	if s.IsDir() {
		return false
	}
	if !s.Mode().IsRegular() {
		return false
	}
	return filepath.Base(s.Name()) == "README_ORDER.txt"
}
