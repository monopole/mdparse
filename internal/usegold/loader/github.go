package loader

import (
	"bytes"
	"fmt"
	"github.com/monopole/mdrip/base"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func loadFromGithub(repoName, relPath string) (dataNode, error) {
	tmpDir, err := cloneRepo(repoName)
	defer deleteTheClone(tmpDir)
	if err != nil {
		return errRaw(repoName, err)
	}
	fullPath := tmpDir
	if len(relPath) > 0 {
		fullPath = filepath.Join(tmpDir, relPath)
	}
	return loadFromPath(base.FilePath(fullPath))
}

func cloneRepo(repoName string) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("maybe no git program? (%w)", err)
	}
	tmpDir, err := os.MkdirTemp("", "mdrip-git-")
	if err != nil {
		return "", fmt.Errorf("unable to create tmp dir (%w)", err)
	}
	slog.Info("Cloning to " + tmpDir)
	cmd := exec.Command(gitPath, "clone", "https://github.com/"+repoName+".git", tmpDir)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failure (%w)", err)
	}
	slog.Info("Clone complete.")
	return tmpDir, nil
}

func deleteTheClone(tmpDir string) {
	_ = os.RemoveAll(tmpDir)
	slog.Info("Deleted " + tmpDir)
}

// smellsLikeGithubCloneArg returns true if the argument seems
// like it could be GitHub url or `git clone` argument.
func smellsLikeGithubCloneArg(arg string) bool {
	arg = strings.ToLower(arg)
	return strings.HasPrefix(arg, "gh:") ||
		strings.HasPrefix(arg, "git@github.com:") ||
		strings.HasPrefix(arg, "https://github.com/")
}

// extractGithubRepoName parses strings like git@github.com:monopole/mdrip.git or
// https://github.com/monopole/mdrip, extracting the repository name
// and the path inside the repository.
func extractGithubRepoName(n string) (string, string, error) {
	for _, p := range []string{
		// Order matters here.
		"gh:", "https://", "http://", "git@", "github.com:", "github.com/"} {
		if strings.ToLower(n[:len(p)]) == p {
			n = n[len(p):]
		}
	}
	if strings.HasSuffix(n, ".git") {
		n = n[0 : len(n)-len(".git")]
	}
	i := strings.Index(n, "/")
	if i < 1 {
		return "", "", fmt.Errorf("no separator in github spec")
	}
	j := strings.Index(n[i+1:], "/")
	if j < 0 {
		// No path, so show entire repo.
		return n, "", nil
	}
	j += i + 1
	return n[:j], n[j+1:], nil
}
