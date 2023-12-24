package loader_test

import (
	"fmt"
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TODO: make lots more tests
// var AppFs = afero.NewMemMapFs()
//
// or
//
// var AppFs = afero.NewOsFs()
// then
//    os.Open("/tmp/foo")
// becomes
//     AppFs.Open("/tmp/foo")
//

func TestExist(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	assert.NoError(t, appFS.MkdirAll("src/a", 0755))
	assert.NoError(t, afero.WriteFile(appFS, "src/a/b", []byte("file b"), 0644))
	assert.NoError(t, afero.WriteFile(appFS, "src/c", []byte("file c"), 0644))
	name := "src/c"
	_, err := appFS.Stat(name)
	if os.IsNotExist(err) {
		t.Errorf("file \"%s\" does not exist.\n", name)
	}
}

func TestLoadFolderFromFsSad(t *testing.T) {
	type testC struct {
		root   string
		path   string
		errMsg string
	}
	l := NewFsLoader(afero.NewOsFs())
	for n, tc := range map[string]testC{
		"e1": {
			root:   "/home/jregan/myrepos/github.com/monopole",
			path:   "yugga",
			errMsg: "unable to read folder",
		},
	} {
		t.Run(n, func(t *testing.T) {
			parent := NewFolder(tc.root)
			_, err := l.LoadSubFolder(parent, tc.path)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestLoadFolderFromFsHappy(t *testing.T) {
	{
		var cwd string
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return
		}
		fmt.Println("cwd of test =", cwd)
	}
	type testC struct {
		root string
		path string
	}
	l := NewFsLoader(afero.NewOsFs())
	for n, tc := range map[string]testC{
		"t1": {
			root: "/home/jregan/myrepos/github.com/monopole",
			path: "mdparse",
		},
		"t2": {
			root: "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader",
			path: ".",
		},
	} {
		t.Run(n, func(t *testing.T) {
			folder := NewFolder(tc.root)
			f, err := l.LoadSubFolder(folder, tc.path)
			assert.NoError(t, err)
			assert.NotNil(t, f)
			f.Accept(NewVisitorDump(l))
		})
	}
}
