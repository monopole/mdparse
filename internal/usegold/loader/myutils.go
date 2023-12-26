package loader

import (
	"path/filepath"
	"strings"
)

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
	if dir == "" && name == currentDir {
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
