package loader

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const dotGit = ".git"

// MyRepo is a named group of files and folders.
type MyRepo struct {
	// name is the URL of the repo, e.g.
	// https://githu.com/monopole/mdrip
	name string
	// path is the path of interest inside the repo, ignore
	// everything else.  If this is empty, take the whole
	// repo module the content that doesn't pass filters.
	path string

	// Holds the repo.
	folder *MyFolder
}

var _ MyTreeItem = &MyRepo{}

func (r *MyRepo) Accept(v TreeVisitor) {
	v.VisitRepo(r)
}

func (r *MyRepo) Parent() MyTreeItem {
	return nil
}

func (r *MyRepo) Root() MyTreeItem {
	return r
}

func (r *MyRepo) FullName() string {
	return r.name
}

func (r *MyRepo) Name() string {
	return r.name
}

// smellsLikeGithubCloneArg returns true if the argument seems
// like it could be GitHub url or `git clone` argument.
func smellsLikeGithubCloneArg(arg string) bool {
	arg = strings.ToLower(arg)
	return strings.HasPrefix(arg, "gh:") ||
		strings.HasPrefix(arg, "git@github.com:") ||
		strings.HasPrefix(arg, "https://github.com/")
}

// CloneAndLoadRepo clones a repo locally and loads it.
// The FsLoader should be injected with a real file system,
// since the git command line used here clones to real disk.
func CloneAndLoadRepo(fsl *FsLoader, arg string) (*MyRepo, error) {
	n, p, err := extractGithubRepoName(arg)
	if err != nil {
		return nil, err
	}
	r := &MyRepo{
		name: n,
		path: p,
	}
	var tmpDir string
	tmpDir, err = cloneRepo(r.name)
	if err != nil {
		return nil, err
	}
	r.folder, err = fsl.LoadFolder(filepath.Join(tmpDir, r.path))
	r.folder.name = r.path
	_ = os.RemoveAll(tmpDir)
	return r, err
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
	if strings.HasSuffix(n, dotGit) {
		n = n[0 : len(n)-len(dotGit)]
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

func cloneRepo(repoName string) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("maybe no git program? (%w)", err)
	}
	tmpDir, err := os.MkdirTemp("", "mdrip-git-")
	if err != nil {
		return "", fmt.Errorf("unable to create tmp dir (%w)", err)
	}
	slog.Info("Cloning", "tmpDir", tmpDir, "repoName", repoName)
	cmd := exec.Command(
		gitPath, "clone", "https://github.com/"+repoName+dotGit, tmpDir)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failure (%w)", err)
	}
	slog.Info("Clone complete.")
	return tmpDir, nil
}
