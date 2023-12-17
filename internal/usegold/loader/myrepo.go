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

// MyRepo is a named group of files and folders.
type MyRepo struct {
	// name is the URL of the repo, e.g.
	// https://githu.com/monopole/mdrip
	name string
	// path is the path of interest inside the repo, ignore
	// everything else.  If this is empty, take the whole
	// repo module the content that doesn't pass filters.
	path string

	// Directory that holds the clone
	tmpDir string

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

func (r *MyRepo) FullName() string {
	return r.name
}

func (r *MyRepo) DirName() string {
	return ""
}

func (r *MyRepo) Name() string {
	return r.name
}

func (r *MyRepo) Init(tb *TreeScanner) (err error) {
	r.CleanUp()
	r.tmpDir, err = cloneRepo(r.name)
	if err != nil {
		return
	}
	base := filepath.Base(r.tmpDir)
	dir := filepath.Dir(r.tmpDir)
	if len(r.path) > 0 {
		base = filepath.Join(base, r.path)
	}
	//fmt.Printf("r.tmpDir = %s\n", r.tmpDir)
	//fmt.Printf("  r.path = %s\n", r.path)
	//fmt.Printf("    base = %s\n", base)
	//fmt.Printf("     dir = %s\n", dir)
	r.folder = &MyFolder{
		myTreeItem: myTreeItem{
			name: dir,
		},
	}
	return r.folder.AbsorbFolderFromDisk(tb, base)
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

func (r *MyRepo) CleanUp() {
	if r.tmpDir == "" {
		return
	}
	_ = os.RemoveAll(r.tmpDir)
	slog.Info("Deleted " + r.tmpDir)
	r.tmpDir = ""
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
