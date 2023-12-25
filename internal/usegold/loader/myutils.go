package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// filter returns an error if conditions not met
type filter func(info os.FileInfo) error

var selfPath = "." + rootSlash

// FSplit splits a path into a parent path and a single name.
// It differs from filepath.Split in how it handles "." and
// trailing slashes. The resulting parent path never has a trailing
// slash unless it is only a slash.  Also, the resulting parent
// path never has a ".".  It's assumed that if the parent path
// has no leading "/", then it's a *relative* path.  An empty
// parent path means the current directory.  See tests.
// The goal is to get a parent path that "looks" nice.
func FSplit(path string) (string, string) {
	dir, name := filepath.Split(path)
	if strings.HasPrefix(dir, selfPath) {
		dir = dir[2:]
	}
	if dir == rootSlash {
		return dir, name
	}
	dir = stripTrailingSlash(dir)
	if dir == "" && name == "." {
		return "", ""
	}
	return dir, name
}

func stripTrailingSlash(path string) string {
	if strings.HasSuffix(path, rootSlash) {
		return path[:len(path)-1]
	}
	return path
}

var NotMarkDownErr = fmt.Errorf("not a simple markdown file")

// IsMarkDownFile passes markdown files.
func IsMarkDownFile(info os.FileInfo) error {
	if !info.Mode().IsRegular() {
		return NotMarkDownErr
	}
	if filepath.Ext(info.Name()) != ".md" {
		return NotMarkDownErr
	}
	base := filepath.Base(info.Name())
	const badLeadingChar = "~.#"
	if strings.Index(badLeadingChar, string(base[0])) >= 0 {
		return NotMarkDownErr
	}
	return nil
}

var IsADotDirErr = fmt.Errorf("not allowed to load from dot folder")

// IsNotADotDir passes non dot directories (not .git, not .config, etc.)
func IsNotADotDir(info os.FileInfo) error {
	n := info.Name()
	// Allow special dir names.
	if n == "." || n == selfPath || n == ".." {
		return nil
	}
	// Ignore .git, etc.
	base := filepath.Base(n)
	if len(base) > 1 && string(base[0]) == "." {
		return IsADotDirErr
	}
	return nil
}
