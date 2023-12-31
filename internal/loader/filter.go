package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// filter returns an error if conditions not met
type filter func(info os.FileInfo) error

var NotMarkDownErr = fmt.Errorf("not a simple markdown file")

// IsMarkDownFile passes markdown files.
func IsMarkDownFile(info os.FileInfo) error {
	if !info.Mode().IsRegular() {
		return NotMarkDownErr
	}
	if filepath.Ext(info.Name()) != ".md" {
		return NotMarkDownErr
	}
	const badLeadingChar = "~.#"
	if strings.Index(badLeadingChar, string(info.Name()[0])) >= 0 {
		return NotMarkDownErr
	}
	return nil
}

var IsADotDirErr = fmt.Errorf("not allowed to load from dot folder")

// IsNotADotDir passes non dot directories (not .git, not .config, etc.)
func IsNotADotDir(info os.FileInfo) error {
	n := info.Name()
	// Allow special dir names.
	if n == currentDir || n == selfPath || n == upDir {
		return nil
	}
	// Ignore .git, etc.
	if len(n) > 1 && string(n[0]) == currentDir {
		return IsADotDirErr
	}
	return nil
}
