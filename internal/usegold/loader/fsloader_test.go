package loader_test

import (
	"fmt"
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
			folder.Accept(NewVisitorDump(l))
		})
	}
}

// Permission bits
// The file or folder's owner:
// 400      Read
// 200      Write
// 100      Execute/search
//
//	        Other users in the file or folder's group:
//	40      Read
//	20      Write
//	10      Execute/search
//
//	        Other users not in the group:
//	 4      Read
//	 2      Write
//	 1      Execute/search
const (
	RW  os.FileMode = 0644
	RWX os.FileMode = 0755
)

// Define a bunch of markdown files and their contents.
// Any file whose name ends in ".md" is considered a markdown file.
var (
	m0, m1, m2, m3 = NewFile("m0.md", []byte("# m0")),
		NewFile("m1.md", []byte("# m1")),
		NewFile("m2.md", []byte("# m2")),
		NewFile("m3.md", []byte("# m3"))
)

func makeSmallAbsFs(t *testing.T, fs afero.Fs) {
	assert.NoError(t, afero.WriteFile(fs, "/m0.md", m0.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/aaa/m1.md", m1.C(), RW))
}

func makeSmallRelFs(t *testing.T, fs afero.Fs) {
	assert.NoError(t, afero.WriteFile(fs, "m0.md", m0.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "aaa/m1.md", m1.C(), RW))
}

func makeMediumAbsFs(t *testing.T, fs afero.Fs) {
	assert.NoError(t, afero.WriteFile(fs, "/m0.md", m0.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/aaa/bbb/m1.md", m1.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/aaa/m2.md", m2.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/aaa/ccc/m3.md", m3.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/aaa/ccc/ignore", []byte("not markdown"), RW))
	assert.NoError(t, fs.MkdirAll("/aaa/empty", RWX))
}

func makeMediumRelFs(t *testing.T, fs afero.Fs) {
	assert.NoError(t, afero.WriteFile(fs, "m0.md", m0.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "aaa/bbb/m1.md", m1.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "aaa/m2.md", m2.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "aaa/ccc/m3.md", m3.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "aaa/ccc/ignore", []byte("not markdown"), RW))
	assert.NoError(t, fs.MkdirAll("aaa/empty", RWX))
}

func TestAferoNonRootPath(t *testing.T) {
	// There's no notion of a "current" working directory
	// that you can change when writing to afero.
	// You can read and write absolute paths, which always start with rootSlash,
	// or you can read/write paths that don't start with rootSlash.
	// There's no way to "cd" in the rootSlash file system,
	// and have that change the behavior read/write.
	fs := afero.NewMemMapFs() // afero.NewOsFs()
	assert.NoError(t, afero.WriteFile(fs, "m0.md", m0.C(), RW))
	assert.NoError(t, afero.WriteFile(fs, "/m1.md", m1.C(), RW))
	data, err := afero.ReadFile(fs, "m0.md")
	assert.NoError(t, err)
	assert.Equal(t, m0.C(), data)
	data, err = afero.ReadFile(fs, "/m1.md")
	assert.NoError(t, err)
	assert.Equal(t, m1.C(), data)
	_, err = afero.ReadFile(fs, "m1.md")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}

func TestLoadFolderFromMemoryHappy(t *testing.T) {
	type testC struct {
		fillFs      func(*testing.T, afero.Fs)
		pathToLoad  string
		expectedFld func() *MyFolder
		errMsg      string
	}
	for n, tc := range map[string]testC{
		"nothingOk": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				// don't make any files.
			},
			pathToLoad: "/",
			expectedFld: func() *MyFolder {
				return nil
			},
		},
		"nothingWithError": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				// don't make any files.
			},
			pathToLoad: "/a.md",
			errMsg:     "file does not exist",
		},
		"nonExistentFile": {
			fillFs:     makeSmallAbsFs,
			pathToLoad: "/monkey",
			errMsg:     "does not exist",
		},
		"noGoingUp": {
			fillFs:     makeSmallAbsFs,
			pathToLoad: "../zzz",
			errMsg:     "specify absolute path or something at or below your working directory",
		},
		"oneFile": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(tt, afero.WriteFile(fs, "/m1.md", m1.C(), RW))
			},
			pathToLoad: "/m1.md",
			expectedFld: func() *MyFolder {
				return NewFolder("/").AddFileObject(m1)
			},
		},
		"oneFileButAskForWrongFile": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(tt, afero.WriteFile(fs, "/m1.md", m1.C(), RW))
			},
			pathToLoad: "/m2.md",
			errMsg:     "file does not exist",
		},
		"oneEmptyFolder": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(t, fs.MkdirAll("/aaa", RWX))
			},
			pathToLoad: "/",
			expectedFld: func() *MyFolder {
				return nil
			},
		},
		"oneEmptyFolderAgain": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(t, fs.MkdirAll("/aaa", RWX))
			},
			pathToLoad: "/aaa",
			expectedFld: func() *MyFolder {
				return nil
			},
		},
		"justOneDir": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(tt, afero.WriteFile(fs, "/aaa/m1.md", m1.C(), RW))
			},
			pathToLoad: "/",
			expectedFld: func() *MyFolder {
				return NewFolder("/").AddFolderObject(NewFolder("aaa").AddFileObject(m1))
			},
		},
		"justOneSubDir": {
			fillFs: func(tt *testing.T, fs afero.Fs) {
				assert.NoError(tt, afero.WriteFile(fs, "/aaa/m1.md", m1.C(), RW))
			},
			pathToLoad: "/aaa",
			expectedFld: func() *MyFolder {
				return NewFolder("/").AddFolderObject(NewFolder("aaa").AddFileObject(m1))
			},
		},
		"allOfSmallAbsFs": {
			fillFs:     makeSmallAbsFs,
			pathToLoad: "/",
			expectedFld: func() *MyFolder {
				aaa := NewFolder("aaa").AddFileObject(m1)
				return NewFolder("/").AddFileObject(m0).AddFolderObject(aaa)
			},
		},
		"allOfSmallRelFs": {
			fillFs:     makeSmallRelFs,
			pathToLoad: ".",
			expectedFld: func() *MyFolder {
				aaa := NewFolder("aaa").AddFileObject(m1)
				return NewFolder(".").AddFileObject(m0).AddFolderObject(aaa)
			},
		},
		"fromMediumAbsEverything": {
			fillFs:     makeMediumAbsFs,
			pathToLoad: "/",
			expectedFld: func() *MyFolder {
				ccc := NewFolder("ccc").AddFileObject(m3)
				bbb := NewFolder("bbb").AddFileObject(m1)
				aaa := NewFolder("aaa").AddFileObject(m2).
					AddFolderObject(bbb).AddFolderObject(ccc)
				return NewFolder("/").AddFileObject(m0).AddFolderObject(aaa)
			},
		},
		"fromMediumRelEverything": {
			fillFs:     makeMediumAbsFs,
			pathToLoad: ".",
			expectedFld: func() *MyFolder {
				ccc := NewFolder("ccc").AddFileObject(m3)
				bbb := NewFolder("bbb").AddFileObject(m1)
				aaa := NewFolder("aaa").AddFileObject(m2).
					AddFolderObject(bbb).AddFolderObject(ccc)
				return NewFolder(".").AddFileObject(m0).AddFolderObject(aaa)
			},
		},
		"fromMediumAbsJustAAA": {
			fillFs:     makeMediumAbsFs,
			pathToLoad: "/aaa",
			expectedFld: func() *MyFolder {
				ccc := NewFolder("ccc").AddFileObject(m3)
				bbb := NewFolder("bbb").AddFileObject(m1)
				aaa := NewFolder("aaa").AddFileObject(m2).
					AddFolderObject(bbb).AddFolderObject(ccc)
				return NewFolder("/").AddFolderObject(aaa)
			},
		},
		"fromMediumAbsJustM0": {
			fillFs:     makeMediumAbsFs,
			pathToLoad: "/m0.md",
			expectedFld: func() *MyFolder {
				return NewFolder("/").AddFileObject(m0)
			},
		},
		"fromMediumAbsJustM3": {
			fillFs:     makeMediumAbsFs,
			pathToLoad: "/aaa/ccc/m3.md",
			expectedFld: func() *MyFolder {
				return NewFolder("/aaa/ccc").AddFileObject(m3)
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			fs := afero.NewMemMapFs() // afero.NewOsFs()
			tc.fillFs(t, fs)
			ldr := NewFsLoader(fs)
			fld, err := ldr.LoadFolder(tc.pathToLoad)
			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				return
			}
			assert.NoError(t, err)
			if !assert.True(t, tc.expectedFld().Equals(fld)) {
				t.Errorf("Didn't get expected folder.")
				t.Log("Loaded:")
				fld.Accept(NewVisitorDump(ldr))
				t.Log("Expected:")
				tc.expectedFld().Accept(NewVisitorDump(ldr))
			}
		})
	}
}
